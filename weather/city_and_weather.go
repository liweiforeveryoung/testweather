package weather

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 根据 adcode 查询城市和天气信息
// 对客户端发出的weather请求的响应
type Live struct {
	City        string  `json:"city"`
	Weather     string  `json:"weather"`
	Temperature float32 `json:"temperature,string"`
	ReportTime  string  `json:"reporttime"` // 该天气信息的产生时间
}

type CityAndWeather struct {
	Status string `json:"status"`
	Info   string `json:"info"`
	Count  int32  `json:"count,string"`
	Lives  []Live `json:"lives"`
}

func getWeatherImformationResponseFromAdcode(adcode string) *WeatherImformationResponse {
	responsePtr := getWeatherImforFromRedis(adcode)
	if responsePtr != nil {
		// 说明redis中有记录
		responsePtr.Status = 0
		return responsePtr
	} else {
		// 说明redis中无记录，需要向高德请求数据
		var response WeatherImformationResponse
		cityAndWeather, err := getCityAndWeatherImforFromGaode(adcode)
		if err != nil {
			// 打印错误
			logrus.WithFields(logrus.Fields{"err": err, "adcode": adcode}).Warn("get city and weather data from adcode (gaode) failed")
			// 发送 status 为1的消息
			response.Status = 1
			return &response
		}
		lives := cityAndWeather.Lives
		if len(lives) <= 0 {
			// 打印错误
			logrus.WithFields(logrus.Fields{"err": err, "adcode": adcode, "lives length": len(lives)}).Warn("get city and weather data from adcode failed (gaode) (invalid lives length)")
			response.Status = 1
			return &response
		} else {
			response.Status = 0
			live := lives[0]
			response.City = live.City
			response.Temperature = live.Temperature
			response.Weather = live.Weather
			response.ReportTime = live.ReportTime
			// 将数据存入redis
			addWeatherImforToRedis(adcode, &response)
			return &response
		}
	}
}

// 从高德获取天气信息
// 记得在外面检查 cityandweather 的count是否大于0，否则可能会出现越界
func getCityAndWeatherImforFromGaode(adcode string) (*CityAndWeather, error) {
	// 拼接 url
	url := weatherUrl + adcode
	// 向高德发出请求
	weatherRspn, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer weatherRspn.Body.Close()
	// 检查http状态
	if weatherRspn.StatusCode != http.StatusOK {
		return nil, errors.New("weather Http response error")
	}
	// 检查 json 的 status状态
	var cw CityAndWeather
	err = json.NewDecoder(weatherRspn.Body).Decode(&cw)
	if err != nil {
		return nil, err
	}
	if cw.Status != "1" {
		return nil, errors.New("request weather data from gaode failed" + cw.Info)
	}
	return &cw, nil
}

// 将 adcode city 天气情况 温度 放入redis
func addWeatherImforToRedis(adcode string, weatherImfor *WeatherImformationResponse) {
	if weatherImfor == nil {
		return
	}
	reportTime, err := time.ParseInLocation(gaodeReportTimeFormat, weatherImfor.ReportTime, time.Local)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reporttime from gaode": weatherImfor.ReportTime}).Warn(
			"parse time failed")
		return
	}
	// 计算时间差
	reportTimeToNowInterval := time.Now().Sub(reportTime)
	// interval 不一定合法，可能会小于零
	if reportTimeToNowInterval < 0 {
		return
	}
	// 如果时间差比过期时间大，那就没必须保存在redis里了，直接返回
	if reportTimeToNowInterval > expireDuration {
		return
	}

	// 获得连接
	conn := pool.Get()
	defer pool.Close()

	reply, err := redis.String(conn.Do("hmset", "adcode:"+adcode, "city", weatherImfor.City, "weather",
		weatherImfor.Weather, "temperature", weatherImfor.Temperature))
	if err != nil {
		logrus.Errorf("add addWeatherImforToRedis failed,%v ", err)
		return
	}
	// 检查reply是否是ok
	if reply != "OK" {
		logrus.WithFields(logrus.Fields{
			"adcode": adcode, "city": weatherImfor.City,
			"weather": weatherImfor.Weather, "temperature": weatherImfor.Temperature,
		}).Warn("addWeatherImforToRedis failed")
		return
	}
	// 得到过期时间
	leftTime := expireDuration - reportTimeToNowInterval
	// 设置过期时间
	conn.Do("expire", "adcode:"+adcode, int32(leftTime.Seconds()))
	// todo 检测设置过期时间是否成功
}

// 从 redis 取出 city 天气情况 温度
func getWeatherImforFromRedis(adcode string) *WeatherImformationResponse {
	// 查询redis中是否具有该key
	// 如果没有就直接返回nil
	conn := pool.Get()
	defer conn.Close()

	values, err := redis.StringMap(conn.Do("hgetall", "adcode:"+adcode))
	if err != nil || len(values) == 0 {
		// 说明map中并不存在该天气信息
		return nil
	} else {
		var response WeatherImformationResponse
		response.City = values["city"]
		response.Weather = values["weather"]
		tempTemperature, _ := strconv.ParseFloat(values["temperature"], 32)
		response.Temperature = float32(tempTemperature)
		return &response
	}
}
