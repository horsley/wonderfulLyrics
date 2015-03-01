package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
	"wonderfulLyrics/qqmusic"
	"wonderfulLyrics/xiami"
)

var StorageSongData DataFile
var StorageSongDataLock sync.RWMutex

//数据存储文件格式
type DataFile struct {
	Data []SongData
}

type SongData struct {
	Mp3Url     string //mp3
	BgImgUrl   string //背景图
	BestLyric  string //最美歌词
	SongInfo   string //歌曲信息
	XiamiSid   int    //虾米音乐id
	QQMusicMID string //qq音乐mid
	QQMusicID  int    //qq音乐id
}

//读取文件数据
func readSongDataFile() error {
	bin, err := ioutil.ReadFile(*dataFile)
	if err != nil {
		log.Println("readSongDataFile:read data err:", err)
		return err
	}

	var d DataFile
	err = json.Unmarshal(bin, &d)
	if err != nil {
		log.Println("readSongDataFile:read data err:", err)
		return err
	}

	StorageSongDataLock.Lock()
	defer StorageSongDataLock.Unlock()
	StorageSongData = d

	log.Println("readSongDataFile:read data finish, song count:", len(d.Data))
	return nil
}

//监视数据文件更新
func monDataFileUpdate() {
	var lastMtime time.Time
	tk := time.NewTicker(5 * time.Second)

	for _ = range tk.C {
		if stat, err := os.Stat(*dataFile); err != nil {
			log.Println("monDataFileUpdate stat data file error:", err)
		} else {
			if stat.ModTime().Sub(lastMtime) > 0 {
				log.Println("monDataFileUpdate find data file update")

				if readSongDataFile() == nil && fillSongMp3Url() == nil {
					if bLog, err := build(); err == nil { //自动重新构建
						log.Println("autobuild succeed:\n" + bLog)
					} else {
						log.Println("autobuild fail:", err)
					}
				}

				lastMtime = stat.ModTime()
			}
		}

	}
}

//检查并填充歌曲地址
func fillSongMp3Url() error {
	StorageSongDataLock.Lock()
	defer StorageSongDataLock.Unlock()

	for k, v := range StorageSongData.Data {
		if v.Mp3Url != "" { //已经获取过地址的 检查是否还可用
			if !test200(v.Mp3Url) { //已失效,置空
				log.Println("fillSongMp3Url: song:", v.SongInfo, "Mp3Url test200 fail!")
				v.Mp3Url = ""
			} else {
				continue
			}
		}

		if v.Mp3Url == "" { //失效的和新的
			var err error
			switch {
			case v.XiamiSid != 0:
				StorageSongData.Data[k].Mp3Url, err = xiami.GetSongUrl(v.XiamiSid)
			case v.QQMusicID != 0:
				StorageSongData.Data[k].Mp3Url, err = qqmusic.GetSongUrlByID(v.QQMusicID)
			case v.QQMusicMID != "":
				StorageSongData.Data[k].Mp3Url, err = qqmusic.GetSongUrlByMID(v.QQMusicMID)
			}

			if err != nil {
				log.Println("fillSongMp3Url:fetch music mp3 url error:", err, "info:", v.SongInfo)
				return err
			} else {
				log.Println("fillSongMp3Url: fill url:", StorageSongData.Data[k].Mp3Url, "for:", v.SongInfo)
			}
		}
	}

	return nil
}
