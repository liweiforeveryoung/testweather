package weather

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

// 根据 ip 向 高德请求 位置信息
type Location struct {
	Status string `json:"status"`
	Info   string `json:"info"`
	Adcode string `json:"adcode"` // 如果是局域网，adcode将为空
}

// 给定一个ip，返回该 ip 的 adcode（调用者需要检查错误）
func getAdcodeFromIP(ip string) (string, error) {
	url := locationUrl + ip
	locationResponse, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer locationResponse.Body.Close()
	// 检查 状态是否为 200
	if locationResponse.StatusCode != http.StatusOK {
		return "", errors.New("location Http response error")
	}
	// 检查 json 的 status 是否为1（1代表成功）
	var loc Location
	err = json.NewDecoder(locationResponse.Body).Decode(&loc)
	if err != nil {
		// 如果 Unmarshal失败的话，较大可能是ip为局域网ip或者国外ip
		// 高德有点反人类呀，当ip合法的时候，他给我返回的adcode是一个字符串
		// 当ip不合法的是否，它给我一个数组
		// 那便将ip不合法的情况当作错误返回算了
		// {"status":"1","info":"OK","infocode":"10000","province":[],"city":[],"adcode":[],"rectangle":[]}`
		// {"status":"1","info":"OK","infocode":"10000","province":"河北省","city":"衡水市","adcode":"131100","rectangle":"115.5630291,37.66458626;115.8063269,37.80817343"}
		return "", err
	}
	// 如果不为1，记录json中的info
	if loc.Status != "1" {
		return "", errors.New("request location data from gaode failed : " + loc.Info)
	}
	// 返回 adcode
	return loc.Adcode, nil

}
