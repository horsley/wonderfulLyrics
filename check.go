package main

import (
	"log"
	"net/http"
	"time"
)

//检查200
func test200(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		log.Println("test200 err:", err, "url:", url)
		return false
	}

	defer resp.Body.Close()
	return resp.StatusCode == 200
}

//定时检查mp3是否失效
func trackMp3UrlValid() {
	tk := time.NewTicker(time.Minute)
	for _ = range tk.C {
		fillSongMp3Url()
	}
}
