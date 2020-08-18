package weather

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// 一些配置信息
const (
	// gin 的 addr
	ginServeAddr = ":8080"
	gaodeKey     = "e92a6d262f97f19923fd04e3c8413b81"

	// 向高德请求位置的url(此时还缺少ip信息)
	locationUrl = `https://restapi.amap.com/v3/ip?output=json&key=` + gaodeKey + "&ip="

	// 向高德请求天气情况的url(此时还缺少adcode)
	weatherUrl = `https://restapi.amap.com/v3/weather/weatherInfo?&extensions=base&output=json&key=` + gaodeKey + "&city="

	// 高德传回来的reportTime的格式
	gaodeReportTimeFormat = `2006-01-02 15:04:05`

	// redis 的 地址
	redisAddr = "localhost:6379"
	// 连接池中的最大空闲连接数量
	maxIdle = 3
	// 连接池中空闲连接的最大等待时间
	idleTimeout = 240 * time.Second
	// 当连接不够用时，get会被阻塞直到有可用连接
	wait = true

	// 天气的过期时间：8小时
	expireDuration = 8 * time.Hour

	// logrus 的日志保存路径
	logPath = "./weather_server_log.txt"
)

// 配置 log 的打印位置

func init() {
	var file, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could Not Open Log File : " + err.Error())
		fmt.Println("log will output to console")
	} else {
		logrus.SetOutput(file)
		// todo 这个file还需要关闭么？
	}
}
