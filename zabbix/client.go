package zabbix

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"os"
	"zabbixGraph/utils"
)

var Browser = New()

type Client struct {
	token   string
	cookies []*http.Cookie
}

func init() {
	err := os.MkdirAll(utils.DownloadDir, 777)
	if err != nil {
		log.Fatal("创建文件夹失败:", err)
		os.Exit(1)
	}
}

func New() *Client {
	var c = new(Client)
	c.getToken()
	c.getCookies()
	fmt.Println("获取Token:", c.token)
	fmt.Println("获取Cookies:", c.cookies)
	return c
}

// getToken 获取token
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
	body := RequestJSON(rj, utils.ApiRpcURL)
	var jb = JSONBody{}
	err := json.Unmarshal([]byte(body), &jb)
	if err != nil {
		fmt.Println("解析失败:", err)
	}
	c.token = jb.Result
}

// zabbix未提供获取图形的api（只能获取图形数据），模仿登录获取数据渲染后的图形
func (c *Client) getCookies() {
	req := gorequest.New()
	res, _, errs := req.Get(utils.LoginURL).
		Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:61.0) Gecko/20100101 Firefox/61.0").
		Query("autologin=1&enter=Sign%20in").
		Query("name="+utils.User).
		Query("password="+utils.Password).
		AppendHeader("Refer", utils.LoginURL).End()
	for _, err := range errs {
		fmt.Println("登录失败:", err)
	}
	c.cookies = res.Request.Cookies()
}
