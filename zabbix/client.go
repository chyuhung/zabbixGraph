package zabbix

import (
	"encoding/json"
	"fmt"
	"zabbixGraph/utils"
)

type Client struct {
	token string
}

func New() *Client {
	var c = new(Client)
	c.getToken()
	fmt.Println("获取Token:", c.token)
	return c
}

/*
// JSONBody 响应json结构体
type JSONBody struct {
	JSONRpc string            `json:"jsonrpc"`
	Result  map[string]string `json:"result"`
	ID      string            `json:"id"`
}
*/
/*
func GetJSONResult(body string) map[string]string {
	var jb = JSONBody{}
	err := json.Unmarshal([]byte(body), &jb)
	if err != nil {
		fmt.Println("解析失败:", err)
	}
	return jb.Result
}
*/

// Login 获取token
func (c *Client) getToken() {
	type JSONBody struct {
		JSONRpc string `json:"JSONRpc,omitempty"`
		Result  string `json:"result,omitempty"`
		ID      string `json:"ID,omitempty"`
	}
	params := map[string]interface{}{
		"user":     utils.User,
		"password": utils.Password,
	}
	//login info
	rj := GetJSONStr(c.token, "user.login", params)
	body := RequestJSON(rj, utils.ApiJSONRpcURL)
	var jb = JSONBody{}
	err := json.Unmarshal([]byte(body), &jb)
	if err != nil {
		fmt.Println("解析失败:", err)
	}
	c.token = jb.Result
}
