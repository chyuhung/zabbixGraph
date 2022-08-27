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
	params := map[string]interface{}{"output": []string{"host"}, "filter": filter}
	rj := GetJSONStr(c.token, "host.get", params)
	body := RequestJSON(rj, utils.ApiRpcURL)
	var jb = JSONBody{}
	err := json.Unmarshal([]byte(body), &jb)
	if err != nil {
		fmt.Println("解析错误:", err)
	}
	return jb.Result[0]["hostid"]
}
func GetGraphID(c *Client, hostID string, graphList []string) string {
	type JSONBody struct {
		JSONRpc string
		Result  []map[string]string
		ID      string
	}
	filter := map[string]interface{}{"name": graphList}
	params := map[string]interface{}{"output": []string{"graphid"}, "hostids": hostID, "filter": filter}
	rj := GetJSONStr(c.token, "graph.get", params)
	body := RequestJSON(rj, utils.ApiRpcURL)
	var jb = JSONBody{}
	err := json.Unmarshal([]byte(body), &jb)
	if err != nil {
		fmt.Println("解析失败:", err)
	}
	return jb.Result[0]["graphid"]
}

// GetGraph 下载图片
func GetGraph(c *Client, filename string, graphID string, timeFrom string, timeTo string, width string, height string) {
	//var v = map[string]interface{}{"graphid": graphID, "from": timeFrom, "to": timeTo, "width": width, "height": height, "profileIdx": "web.charts.filter"}
	request := gorequest.New()
	_, body, errs := request.Get(utils.GraphURL).
		Query("graphid=" + graphID).
		Query("from=" + timeFrom).
		Query("to=" + timeTo).
		Query("width=" + width).
		Query("height=" + height).
		Query("profileIdx=" + "web.charts.filter").
		AddCookies(c.cookies).
		End()
	if errs != nil {
		fmt.Println("请求失败:", errs)
		os.Exit(1)
	}
	f, err := os.Create(utils.DownloadDir + filename + ".png")
	if err != nil {
		fmt.Println("创建文件失败:", err)
	}
	defer f.Close()
	_, err = f.Write([]byte(body))
	if err != nil {
		fmt.Println("写文件失败:", err)
	}

}
