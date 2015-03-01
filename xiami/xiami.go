// xiami project xiami.go
package xiami

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	XIAMI_INFO_JSONP = `http://www.xiami.com/song/playlist/id/%d/object_name/default/object_id/0/cat/json?r=song%2Fdetail&song_id=%d&callback=jsonp1`
)

type jsonpRst struct { //只要关键字段
	Data jsonpData `json:"data"`
}

type jsonpData struct {
	TrackList []SongInfo `json:"trackList"`
}
type SongInfo struct { //只要关键字段
	Location string `json:"location"`
}

//根据虾米id获取歌曲url
func GetSongUrl(songId int) (string, error) {
	resp, err := http.Get(fmt.Sprintf(XIAMI_INFO_JSONP, songId))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if bytes.Compare([]byte(" jsonp1("), bin[:8]) == 0 {
		bin = bin[8 : len(bin)-1]
	}

	var rst jsonpRst
	err = json.Unmarshal(bin, &rst)
	if err != nil {
		return "", err
	}

	return LocationDecode(rst.Data.TrackList[0].Location), nil
}

func LocationDecode(loc string) string {
	line := int(loc[0] - '0')
	stack := make([]string, line)
	lCount := (len(loc) - 1) % line //长的行数
	sLen := (len(loc) - 1) / line   //短的行长度

	lastPos := 1
	for i := 0; i < line; i++ {
		if i < lCount {
			stack[i] = loc[lastPos : lastPos+sLen+1]
			lastPos = lastPos + sLen + 1
		} else {
			stack[i] = loc[lastPos : lastPos+sLen]
			lastPos = lastPos + sLen
		}
	}

	//fmt.Printf("%#v", stack)

	var rst string
	for j := 0; j < sLen+1; j++ {
		for i := 0; i < line; i++ { //行
			if j+1 > len(stack[i]) {
				break
			}
			//fmt.Println(i, j, string(stack[i][j]))
			rst += string(stack[i][j])
		}
	}

	rst, _ = url.QueryUnescape(rst)
	rst = strings.Replace(rst, "^", "0", -1)

	return rst
}
