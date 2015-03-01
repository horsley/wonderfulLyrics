package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DetailPageData struct {
	SongData

	//下面的不在数据文件中记录
	PageTitle string //网页标题
	HomeUrl   string //主页url
}

//首页数据
type IndexPageData struct {
	Songs       []DetailPageData
	NextPageUrl string
	PrevPageUrl string
}

//构建过程
//返回构建日志和错误
//错误信息也会包含在构建日志中
func build() (string, error) {
	var buildLog string

	StorageSongDataLock.RLock()
	defer StorageSongDataLock.RUnlock()

	//复制一份倒序的用于首页列表
	indexReverseData := make([]DetailPageData, len(StorageSongData.Data))

	//播放页面构建
	for k, s := range StorageSongData.Data {
		dpd := DetailPageData{SongData: s}
		dpd.PageTitle = fmt.Sprintf("%s - %s", dpd.BestLyric, SITE_TITLE)
		dpd.HomeUrl = *webRootPrefix + `/index.html`

		pageData := parseHtml(DEFAULT_TMPL, dpd)

		//子目录和绝对路径替换问题
		pageData = bytes.Replace(pageData, []byte(`href="/`),
			[]byte(`href="`+*webRootPrefix+`/`), -1)
		pageData = bytes.Replace(pageData, []byte(`src="/`),
			[]byte(`src="`+*webRootPrefix+`/`), -1)

		//打上build时间标记
		timeTag := fmt.Sprint("<!--", "build time:", time.Now(), "-->")
		pageData = append(pageData, []byte(timeTag)...)

		//写入文件,这里文件名是歌曲信息的sha1
		buildHash := sha1([]byte(fmt.Sprintf("%s", dpd.BestLyric, dpd.SongInfo)))
		buildFile := filepath.Join(*buildTgtDir, buildHash+".html")
		err := ioutil.WriteFile(buildFile, pageData, 0666)
		if err != nil {
			buildLog += fmt.Sprintln("build:ioutil.WriteFile err:", err)
			return buildLog, err
		} else {
			buildLog += fmt.Sprintln("build: build file:", buildFile, "size:", len(pageData), "succeed")
		}

		//复用这个字段用于制作首页链接
		dpd.HomeUrl = *webRootPrefix + `/` + buildHash + ".html"
		indexReverseData[len(StorageSongData.Data)-k-1] = dpd
	}

	//首页构建
	buildFile := filepath.Join(*buildTgtDir, "index.html")
	pageData := parseHtml(INDEX_TMPL, IndexPageData{Songs: indexReverseData})
	err := ioutil.WriteFile(buildFile, pageData, 0666)
	if err != nil {
		buildLog += fmt.Sprintln("build:ioutil.WriteFile err:", err)
		return buildLog, err
	} else {
		buildLog += fmt.Sprintln("build: build file:", buildFile, "size:", len(pageData), "succeed")
	}

	return buildLog, nil
}

//清理老页面过程
func clean() (string, error) {
	var buildLog string

	dir, err := os.Open(*buildTgtDir)
	if err != nil {
		buildLog += fmt.Sprintln("clean:build clean err:", err.Error())
		return buildLog, err
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		buildLog += fmt.Sprintln("clean:build clean err:", err.Error())
		return buildLog, err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasSuffix(f.Name(), ".html") {
			if err := os.Remove(filepath.Join(*buildTgtDir, f.Name())); err != nil {
				buildLog += fmt.Sprintln("clean:remove file fail:", f.Name())
				return buildLog, err
			} else {
				buildLog += fmt.Sprintln("clean:remove file:", f.Name())
			}
		}
	}

	return buildLog, nil
}

//parseHtml 解析html模板
func parseHtml(file string, data interface{}) []byte {
	t, err := template.ParseFiles(file)
	if err != nil {
		return []byte("parse tmpl error:" + err.Error())
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return []byte("exec tmpl error:" + err.Error())
	}

	return buf.Bytes()
}
