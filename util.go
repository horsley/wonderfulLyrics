package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//stdout 和 http 同步打印错误
func printError(rw http.ResponseWriter, a ...interface{}) {
	str := fmt.Sprint(a...)
	http.Error(rw, str, http.StatusInternalServerError)
	log.Println(str)
}

//sha1
func sha1(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

//改变工作目录到可执行文件所在目录
//否则相对路径读写的文件可能会有问题
//这里的方法在windows下不可用
func switchPwd() {
	if exe, err := os.Readlink("/proc/self/exe"); err != nil {
		log.Println("switchPwd:read exe path err:", err)
		os.Exit(1)
	} else {
		wd := filepath.Dir(exe)
		if err := os.Chdir(wd); err != nil {
			log.Println("switchPwd:chdir to path:", wd, " err:", err)
			os.Exit(1)
		}
	}
}

//初始化文件和stdout双输出日志
func initLog(file string) {
	logFile, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("open log file:", file, "error:", err)
		os.Exit(1)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}
