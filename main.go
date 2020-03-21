package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"videocvt/video"
)


var protocol string
var host string
var httpPort int
var filePath string
var filePort int
var loglevel string
var ffmpegPath string

func main() {
	flag.StringVar(&protocol,"protocol","http","protocol")
	flag.StringVar(&host,"host","127.0.0.1","listen host")
	flag.IntVar(&httpPort,"port",80,"http file server port")
	flag.StringVar(&filePath,"fpath","/var/www/video/hls/","http file server path")
	flag.IntVar(&filePort,"fport",81,"http port")
	flag.StringVar(&loglevel,"level","info","loglevel")
	flag.StringVar(&ffmpegPath,"ffmpegPath","/usr/local/bin/ffmpeg","ffmpeg path")
	flag.Parse()
	video.FfmpegPath = ffmpegPath
	l,_ := logrus.ParseLevel(loglevel)
	logrus.SetLevel(l)
	go StartFileServer(filePort,filePath)

	StartHttpServer(httpPort)
	//cmd.Run("docker", "pull","nginx:latest")
	//time.Sleep(time.Hour)



}
