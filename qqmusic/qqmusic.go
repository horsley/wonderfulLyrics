package qqmusic

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	QQMUSIC_PLAY_PAGE  = `http://data.music.qq.com/playsong.html?songid=%d`
	QQMUSIC_PLAY_PAGE1 = `http://data.music.qq.com/playsong.html?songmid=%s`
	IPHONE_UA          = `Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X; en-us) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11A465 Safari/9537.53`
)

var qqh5re = regexp.MustCompile(`<audio.*?src="(.*?)"`)

//通过数字id获取歌曲
func GetSongUrlByID(songId int) (string, error) {
	return getSongUrl(fmt.Sprintf(QQMUSIC_PLAY_PAGE, songId))
}

//通过字符串mid获取歌曲
func GetSongUrlByMID(songMID string) (string, error) {
	return getSongUrl(fmt.Sprintf(QQMUSIC_PLAY_PAGE1, songMID))
}

//fetch and match
func getSongUrl(url string) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", IPHONE_UA)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//fmt.Println(string(bin))

	if match := qqh5re.FindSubmatch(bin); match == nil {
		return "", errors.New("regexp match error")
	} else {
		return string(match[1]), nil
	}

}
