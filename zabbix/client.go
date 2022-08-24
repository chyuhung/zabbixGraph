package zabbix

import (
	"fmt"
	"zabbixGraph/utils"
)

type Client struct {
	token string
}

func New() *Client {
	var c = new(Client)
	c.getToken()
	fmt.Println("获取Token成功:", c.token)
	return c
}

// Login 获取token
func (c *Client) getToken() {
	params := map[string]interface{}{
		"user":     utils.User,
		"password": utils.Password,
	}
	//login info
	reqJson := GetJsonStr(c.token, "user.login", params)
	body := RequestJson(reqJson, utils.ApiJsonRpcUrl)
	c.token = GetJsonValue(body, "result").(string)
}
