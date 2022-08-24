package zabbix

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"os"
	"time"
	"zabbixGraph/utils"
)

// JsonData 请求json结构体
type JsonData struct {
	Jsonrpc string                 `json:"jsonrpc,omitempty"`
	Method  string                 `json:"method,omitempty"`
	Params  map[string]interface{} `json:"params,omitempty"`
	Id      string                 `json:"id,omitempty"`
	Auth    string                 `json:"auth,omitempty"`
}

// GetJsonStr 构造json字符串
func GetJsonStr(token string, method string, params map[string]interface{}) string {
	var m = JsonData{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      "1",
		Auth:    token,
	}
	jsonStr, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	return string(jsonStr)
}

/*
// Params 参数结构体
type Params struct {
	User     string                 `json:"user,omitempty"`
	Password string                 `json:"password,omitempty"`
	Output   string                 `json:"output,omitempty"`
	Filter   map[string]interface{} `json:"filter,omitempty"`
}
*/

// RequestJson 通过json请求数据，结果为string类型
func RequestJson(json string, url string) string {
	request := gorequest.New()
	_, body, errs := request.Post(url).
		Set("Content-Type", "application/json").
		Send(json).
		Timeout(3 * time.Second).
		End()
	if errs != nil {
		fmt.Println("请求错误:", errs)
		os.Exit(1)
	}
	return body
}

/*
// JsonBody 响应json结构体
type JsonBody struct {
	JsonRpc string                 `json:"jsonrpc"`
	Result  map[string]interface{} `json:"result"`
	Id      string                 `json:"id"`
}
*/

func GetJsonValue(jsonBody string, key string) interface{} {
	var m = map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonBody), &m)
	if err != nil {
		fmt.Println("获取错误:", err)
	}
	return m[key]
}

// GetHostId 通过ip获取hostid
func GetHostId(c *Client, ip string) interface{} {
	filter := map[string]interface{}{"ip": ip}
	params := map[string]interface{}{"output": "", "filter": filter}
	reqJson := GetJsonStr(c.token, "host.get", params)
	body := RequestJson(reqJson, utils.ApiJsonRpcUrl)
	m := GetJsonValue(body, "result")
	return m
}
