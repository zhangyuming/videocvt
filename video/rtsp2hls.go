package video

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
	"videocvt/cmd"
	"videocvt/util"

)

var FfmpegPath = ""
var pullStreamTimeout = 30
var s = []context{}

type context struct {
	ID string
	Source string
	Result string
	PID int
}
const HlsSuffix  = "index.m3u8"


type Rtsp2Hls struct {}

func (r *Rtsp2Hls)Convert(source string,dist string)(result string,err error)  {

	c,err := getContextBySource(source)
	if err == nil {
		return c.Result,nil
	}


	c.ID = util.Md5Sum(source)
	c.Source = source
	c.Result = c.ID

	os.MkdirAll(dist + "/" + c.ID + "/",0644)


	stc := make(chan cmd.CmdStatus,3)

	go cmd.Run(stc,FfmpegPath, "-i", source, "-fflags", "flush_packets",
		"-max_delay", "2", "-flags", "-global_header", "-hls_time", "2",
		"-hls_list_size", "3", "-vcodec", "copy", "-y", dist + "/" + c.ID + "/" + HlsSuffix)
	//
	//go cmd.Run(stc,"C:/program/ffmpeg/bin/ffmpeg.exe", "-i", source, "-fflags", "flush_packets",
	//	"-max_delay", "2", "-flags", "-global_header", "-hls_time", "2",
	//	"-hls_list_size", "3", "-vcodec", "copy", "-y", dist + "/" + c.ID + "/" + HlsSuffix)

	v := <-stc
	if v.Status == "running" {
		c.PID = v.Pid
	}else if v.Error != nil {
		return "",v.Error
	}

	start := time.Now()
	for{
		if b,_ := util.PathExists(dist + "/" + c.ID + "/" + HlsSuffix); b {
			break
		}else{
			time.Sleep(time.Millisecond*100)
		}
		d,_ := time.ParseDuration(strconv.Itoa(pullStreamTimeout)+"s")
		if time.Now().After(start.Add(d)){
			logrus.Info("timeout pull stream ",source)
			cmd.Kill(c.PID)
			os.RemoveAll(dist + "/" + c.ID + "/")
			return "",errors.New("pull stream timeout")
		}
	}
	addContext2S(c)
	return c.Result,nil
}

func (r *Rtsp2Hls)TeardownById(key string,dist string)(error){

	c, err := getContextByID(key)
	if c.ID == "" {
		return nil
	}
	if err !=nil {
		os.RemoveAll(dist + "/" + c.ID + "/")
		deleteContext(c)
		logrus.Warn("no found process: ", key , err)

		return nil
	}
	if err := cmd.Kill(c.PID); err != nil {
		return err
	}
	deleteContext(c)
	os.RemoveAll(dist + "/" + c.ID + "/")
	return nil

}

func (r *Rtsp2Hls)TeardownBySource(source string,dist string)error{
	c, err := getContextBySource(source)
	if c.ID == "" {
		logrus.Debug(source , " context already remove")
		return nil
	}
	if err !=nil {
		os.RemoveAll(dist + "/" + c.ID + "/")
		deleteContext(c)
		logrus.Warn("no found process: ", source , err)
		return nil
	}
	if err := cmd.Kill(c.PID); err != nil {
		return err
	}
	deleteContext(c)
	os.RemoveAll(dist + "/" + c.ID + "/")
	return nil

}


func (r *Rtsp2Hls)Reset(source string,dist string)(result string,err error){

	r.TeardownBySource(source,dist)

	c,err := getContextBySource(source)
	if err == nil {
		return c.Result,nil
	}


	c.ID = util.Md5Sum(source)
	c.Source = source
	c.Result = c.ID

	os.MkdirAll(dist + "/" + c.ID + "/",0644)


	stc := make(chan cmd.CmdStatus,3)

	go cmd.Run(stc,FfmpegPath, "-i", source, "-fflags", "flush_packets",
		"-max_delay", "2", "-flags", "-global_header", "-hls_time", "2",
		"-hls_list_size", "3", "-vcodec", "copy", "-y", dist + "/" + c.ID + "/" + HlsSuffix)
	//go cmd.Run(stc,"C:/program/ffmpeg/bin/ffmpeg.exe", "-i", source, "-fflags", "flush_packets",
	//	"-max_delay", "2", "-flags", "-global_header", "-hls_time", "2",
	//	"-hls_list_size", "3", "-vcodec", "copy", "-y", dist + "/" + c.ID + "/" + HlsSuffix)

	v := <-stc
	if v.Status == "running" {
		c.PID = v.Pid
	}else if v.Error != nil {
		return "",v.Error
	}


	start := time.Now()
	for{
		if b,_ := util.PathExists(dist + "/" + c.ID + "/" + HlsSuffix); b {
			break
		}else{
			time.Sleep(time.Millisecond*100)
		}
		d,_ := time.ParseDuration(strconv.Itoa(pullStreamTimeout)+"s")
		if time.Now().After(start.Add(d)){

			logrus.Info("timeout pull stream ",source)
			cmd.Kill(c.PID)
			os.RemoveAll(dist + "/" + c.ID + "/")
			return "",errors.New("pull stream timeout")
		}
	}
	addContext2S(c)
	return c.Result,nil


}




func addContext2S(c context){
	for _,a := range s{
		if a.ID == c.ID{
			return
		}
	}
	s = append(s, c)
}
func getContextByID(id string) (context,error){
	for _,s1 := range s {
		if s1.ID == id {
			return s1,nil
		}
	}
	return context{},errors.New("no found")

}

func getContextBySource(source string) (context,error)  {
	for _,s1 := range s {
		if s1.Source == source {
			return s1,nil
		}
	}
	return context{},errors.New("no found")
}

func deleteContext(c context)  {
	for index,c1 := range s{
		if c1.ID == c.ID {
			s = append(s[:index], s[index+1:]...)
		}
	}
}

func GetContextIdbySource(source string) string  {
	for _,s1 := range s {
		if s1.Source == source {
			return s1.ID
		}
	}
	return ""
}