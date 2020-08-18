package weather

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 开始服务
func Serve() {
	g := gin.Default()
	g.GET("/weather", requestWeatherHandler)
	err := g.Run(ginServeAddr) // 默认运行在 8080 端口
	if err != nil {
		logrus.Fatal("start server failed,%v", err)
	}
}

// 对客户端发出的weather请求的响应
type WeatherImformationResponse struct {
	Status      int8    `json:"status,string"` // 状态：若为零表示正常，非零为异常
	City        string  `json:"city"`
	Weather     string  `json:"weather"`
	Temperature float32 `json:"temperature"`
	ReportTime  string  `json:"-"`
}

func requestWeatherHandler(context *gin.Context) {
	// 获得 context 的ip
	ip := context.ClientIP()
	// 向 高德 请求 adcode
	adcode, err := getAdcodeFromIP(ip)
	if err != nil {
		// 打印错误
		logrus.WithFields(logrus.Fields{"err": err, "ip": ip}).Warn("get adcode from ip failed")
		// 发送 status 为1的消息
		var response WeatherImformationResponse
		response.Status = 1
		context.JSON(http.StatusOK, response)
	} else {
		response := getWeatherImformationResponseFromAdcode(adcode)
		context.JSON(http.StatusOK, response)
	}
}
