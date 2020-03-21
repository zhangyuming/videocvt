package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
	"videocvt/video"
)

// 如果某个请求流超过当前时间则把后台转换任务停掉，默认为秒
const reuestTimeout  = 30


func StartFileServer(port int, path string)  {

	router := gin.Default()
	router.Use(Cors())
	router.Use(Filter())
	router.Static("/", path)
	router.Run(":" + strconv.Itoa(port))

}


func StartHttpServer(port int){
	router := gin.Default()
	router.Use(Cors())
	router.POST("/convertRtsp2Hls",func(c *gin.Context) {
		bt,err := c.GetRawData()
		if err != nil {
			logrus.Error("get post body error",err)
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"get post body error, " + err.Error(),
			})
			return
		}
		body := string(bt)
		body = strings.TrimSpace(body)
		if body == "" {
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"request body is none",
			})
			return
		}
		logrus.Debug("request body is ", body)

		rthls := video.Rtsp2Hls{}

		result,err := rthls.Convert(body,filePath)
		IdMap[result] = 1
		checkBackJob(result)

		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{
				"message":err.Error(),
			})
		}else{
			c.JSON(http.StatusOK,gin.H{
				"message": protocol + "://" +host + ":" + strconv.Itoa(filePort) + "/" + result + "/" + video.HlsSuffix,
			})
		}


	})


	router.POST("/resetRtsp2Hls",func(c *gin.Context) {
		bt,err := c.GetRawData()
		if err != nil {
			logrus.Error("get post body error",err)
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"get post body error, " + err.Error(),
			})
			return
		}
		body := string(bt)
		body = strings.TrimSpace(body)
		if body == "" {
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"request body is none",
			})
			return
		}
		logrus.Debug("request body is ", body)

		rthls := video.Rtsp2Hls{}
		result,err := rthls.Convert(body,filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{
				"message":err.Error(),
			})
		}else{
			c.JSON(http.StatusOK,gin.H{
				"message": protocol + "://" +host + ":" + strconv.Itoa(filePort) + "/" + result + "/" + video.HlsSuffix,
			})
		}

	})

	router.POST("/shutdownRtsp2Hls", func(c *gin.Context) {
		bt,err := c.GetRawData()
		if err != nil {
			logrus.Error("get post body error",err)
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"get post body error, " + err.Error(),
			})
			return
		}
		body := string(bt)
		body = strings.TrimSpace(body)
		if body == "" {
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"request body is none",
			})
			return
		}
		logrus.Debug("request body is ", body)

		rthls := video.Rtsp2Hls{}
		if err := rthls.TeardownBySource(body,filePath); err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{
				"message":err.Error(),
			})
		}else{
			id := video.GetContextIdbySource(body)
			delete(IdMap,id)
			c.JSON(http.StatusOK,gin.H{
				"message": "success",
			})
		}

	})

	router.Run(":"+strconv.Itoa(port))
}


var IdMap = make(map[string]int)

func checkBackJob(key string)  {
	//每秒钟 把固定的请求加1
	go func(m map[string]int, key string) {
		for range time.Tick(time.Second){
			if v,ok :=m[key]; ok {
				logrus.Trace("add key[",key," time +1 current value is ",v)
				m[key] = m[key] + 1
			}else{
				return
			}
		}

	}(IdMap,key)
	// 规定时间内如果固定请求的map值从超过指定时间 则干掉执行的后台任务
	go func(m map[string]int, key string) {
		for range time.Tick(time.Second*5){
			if v,ok :=m[key]; ok {
				logrus.Trace("check key[",key,"] time is " ,v)
				if v > reuestTimeout {
					rt := video.Rtsp2Hls{}
					rt.TeardownById(key,filePath)
					delete(m,key)
				}
			}else{
				return
			}
		}
	}(IdMap,key)
}

func Filter()gin.HandlerFunc{
	return func(c *gin.Context) {
		url := c.Request.RequestURI
		///011ef62186fea63bff04d7c42f2122d9/index5.ts
		key := strings.Split(url,"/")[1]
		if len(key) == 32 {
			IdMap[key] = 1
		}
		c.Next()
	}
}


func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method      //请求方法
		origin := c.Request.Header.Get("Origin")        //请求头部
		var headerKeys []string                             // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")        // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")      //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")      // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")        // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")       //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")       // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next()        //  处理请求
	}
}

