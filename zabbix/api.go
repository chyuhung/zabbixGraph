package zabbix

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"os"
	"strconv"
	"strings"
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
	filter := map[string]interface{}{"host": ip}
	params := map[string]interface{}{"output": []string{"host"}, "filter": filter}
	rj := GetJSONStr(c.token, "host.get", params)
	body := RequestJSON(rj, utils.ApiRpcURL)
	var jb = JSONBody{}
	err := json.Unmarshal([]byte(body), &jb)
	if err != nil {
		fmt.Println("解析错误:", err)
	}
	if len(jb.Result) == 0 {
		fmt.Println("获取hostid失败,ip:", ip)
		return ""
	} else {
		return jb.Result[0]["hostid"]
	}
}
func GetGraphID(c *Client, hostID string, graphList []string) []map[string]string {
	type JSONBody struct {
		JSONRpc string
		Result  []map[string]string
		ID      string
	}
	filter := map[string]interface{}{"name": graphList}
	params := map[string]interface{}{"output": []string{"graphid", "name"}, "hostids": hostID, "filter": filter}
	rj := GetJSONStr(c.token, "graph.get", params)
	body := RequestJSON(rj, utils.ApiRpcURL)
	var jb = JSONBody{}
	err := json.Unmarshal([]byte(body), &jb)
	if err != nil {
		fmt.Println("解析失败:", err)
	}
	if len(jb.Result) == 0 {
		fmt.Println("获取graphid失败,hostid:", hostID)
		return nil
	} else {
		return jb.Result
	}
}

func GetFilename(name string) string {
	filename := strings.ReplaceAll(name, ".", "-")
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "/", "~")
	filename = strings.ReplaceAll(filename, ":", "")
	filename = strings.ReplaceAll(filename, "\n", "")
	return filename
}

func createFile(filePath string, data []byte) {
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println("创建文件失败:", err)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		fmt.Println("写文件失败:", err)
	}
}

// GetGraph 下载图片
func GetGraph(c *Client, ip string, graphID string, graphName string) {
	request := gorequest.New()
	_, body, errs := request.Get(utils.GraphURL).
		Query("graphid=" + graphID).
		Query("from=" + utils.TimeFrom).
		Query("to=" + utils.TimeTo).
		Query("width=" + utils.Width).
		Query("height=" + utils.Height).
		Query("profileIdx=" + "web.charts.filter").
		AddCookies(c.cookies).
		End()
	if errs != nil {
		fmt.Println("请求失败:", errs)
		os.Exit(1)
	}
	createFile(utils.DownloadDir+"/"+ip+"-"+GetFilename(graphName)+".png", []byte(body))

}

func replaceNR(s string) string {
	for {
		if strings.HasSuffix(s, "\n") || strings.HasSuffix(s, "\r") {
			s = strings.ReplaceAll(s, "\n", "")
			s = strings.ReplaceAll(s, "\r", "")
		} else {
			return s
		}
	}
	return s
}
func GetHostList(filename string) []string {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("获取主机失败:", err)
	}
	if len(data) == 0 {
		fmt.Println("主机列表为空")
		return nil
	} else {
		hostsStr := string(data)
		hostsList := strings.SplitAfter(hostsStr, "\n")
		for i, _ := range hostsList {
			hostsList[i] = replaceNR(hostsList[i])
		}
		return hostsList
	}
}

func DownloadGraph() {
	t, err := strconv.Atoi(utils.DownloadSpeed)
	if err != nil {
		fmt.Println(err)
	}
	hostList := GetHostList(utils.HostsFile)
	for _, ip := range hostList {
		hostID := GetHostID(Browser, ip)
		if hostID != "" {
			fmt.Printf("------------------------------\n%s\nhostid:%s\n", ip, hostID)
			graphMaps := GetGraphID(Browser, hostID, utils.GraphNameList)
			for i, _ := range graphMaps {
				graphID := graphMaps[i]["graphid"]
				graphName := graphMaps[i]["name"]
				if graphID != "" {
					fmt.Printf("graphid:%s %s\n", graphID, graphName)
					time.Sleep(time.Duration(t) * time.Second)
					GetGraph(Browser, ip, graphID, graphName)
				}
			}
		}
	}
}
