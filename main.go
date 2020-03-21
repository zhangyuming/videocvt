package main

import (
	"flag"
	"fmt"
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
	flag.IntVar(&filePort,"fport",81,"http port")
	flag.StringVar(&loglevel,"level","info","loglevel")
	//flag.StringVar(&filePath,"fpath","/var/www/video/hls/","http file server path")
	//flag.StringVar(&ffmpegPath,"ffmpegPath","/usr/local/bin/ffmpeg","ffmpeg path")
	flag.StringVar(&ffmpegPath,"ffmpegPath","C:/program/ffmpeg/bin/ffmpeg.exe","ffmpeg path")
	flag.StringVar(&filePath,"fpath","D:/tmp/","http file server path")
	flag.Parse()
	video.FfmpegPath = ffmpegPath
	//l,_ := logrus.ParseLevel(loglevel)
	//logrus.SetLevel(l)
	//go StartFileServer(filePort,filePath)
	//
	//StartHttpServer(httpPort)
	//cmd.Run("docker", "pull","nginx:latest")
	//time.Sleep(time.Hour)

	var sl = []int{3,5,34,21,43,43}
	for index,s := range sl{
		fmt.Println(index, s)
		if s == 43{
			fmt.Print(index)
			sl = append(sl[:index], sl[index+1:]...)
		}

	}
	fmt.Print(sl)


}
