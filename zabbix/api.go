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
	JSONRpc string                 `json:"jsonrpc,omitempty"`
	Method  string                 `json:"method,omitempty"`
	Params  map[string]interface{} `json:"params,omitempty"`
	ID      string                 `json:"id,omitempty"`
	Auth    string                 `json:"auth,omitempty"`
}

// GetJSONStr 构造json字符串
func GetJSONStr(token string, method string, params map[string]interface{}) string {
	var m = JsonData{
		JSONRpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      "1",
		Auth:    token,
	}
	js, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	return string(js)
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

// RequestJSON 通过json请求数据，结果为string类型
func RequestJSON(json string, url string) string {
	request := gorequest.New()
	_, body, errs := request.Post(url).
		Set("Content-Type", "application/json").
		Send(json).
		Timeout(3 * time.Second).
		End()
	if errs != nil {
		fmt.Println("请求失败:", errs)
		os.Exit(1)
	}
	return body
}

// GetHostID 通过ip获取hostid
func GetHostID(c *Client, ip string) string {
	type JSONBody struct {
		JSONRpc string              `json:"JSONRpc,omitempty"`
		Result  []map[string]string `json:"result,omitempty"`
		ID      string              `json:"ID,omitempty"`
	}
	filter := map[string]interface{}{"ip": ip}
	host := []string{"hostid"}
	params := map[string]interface{}{"output": host, "filter": filter}
	rj := GetJSONStr(c.token, "host.get", params)
	body := RequestJSON(rj, utils.ApiJSONRpcURL)
	var jb = JSONBody{}
	err := json.Unmarshal([]byte(body), &jb)
	if err != nil {
		fmt.Println("解析错误:", err)
	}
	return jb.Result[0]["hostid"]
}
