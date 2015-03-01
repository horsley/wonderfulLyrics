// wonderfulLyrics project main.go
package main

import (
	"flag"
	"log"
	"net/http"
)

const (
	DEFAULT_LISTEN    = "127.0.0.1:15215"
	DEFAULT_BUILD_DIR = "build/"
	DEFAULT_DATA      = "data.json"
	DEFAULT_LOGFILE   = "log.txt"
	DEFAULT_TMPL      = "tmpl/detail.html"
	INDEX_TMPL        = "tmpl/index.html"
	SITE_TITLE        = "最美歌词"
)

var (
	//监听地址
	listen = flag.String("listen", DEFAULT_LISTEN, "listen addr")
	//构建输出目录
	buildTgtDir = flag.String("build", DEFAULT_BUILD_DIR, "build dir")
	//web子目录前缀
	webRootPrefix = flag.String("prefix", "", "asset resource prefix")
	//数据文件路径
	dataFile = flag.String("data", DEFAULT_DATA, "song data")
	//日志文件路径
	logFile = flag.String("log", DEFAULT_LOGFILE, "output log file")
)

func init() {
	//为了让相对路径生效，切换pwd到文件所在目录
	switchPwd()

	flag.Parse()
	initLog(*logFile)
}

func main() {
	log.Println("wonderfulLyrics server start!")

	//数据文件监控,读取,更新
	go monDataFileUpdate()
	//内存歌曲数据mp3检查更新  @todo:图片失效检测
	go trackMp3UrlValid()

	//纯静态文件serve，主要是本地测试用
	http.Handle("/", http.FileServer(http.Dir("build")))

	//这个是构建方法，构建以后全是静态的文件，可以不用本程序serve
	//如果要用本程序serve，前端的nginx应该拦截/build这个内部接口
	http.HandleFunc("/build", buildHandler)

	log.Println(http.ListenAndServe(*listen, nil))
}

//手工触发构建
func buildHandler(rw http.ResponseWriter, r *http.Request) {
	var buildLog string
	if r.URL.Query().Get("clean") != "" {
		//可选清理老文件,一般不需要,因为文件名是稳定的,重新构建会覆盖旧文件
		if cLog, err := clean(); err != nil {
			log.Println(cLog)
			rw.Write([]byte(cLog))
			return
		} else {
			buildLog = cLog
		}
	}

	bLog, _ := build()
	buildLog += bLog

	log.Println(buildLog)
	rw.Write([]byte(buildLog))
	return
}
