package zabbix

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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
	var c = &Client{}
	c.getToken()
	c.getCookies()
	fmt.Println("Token:", c.token)
	fmt.Println("Cookies:", c.cookies)
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
	loginURL := utils.LoginURL
	user := utils.User
	password := utils.Password
	v := url.Values{
		"name":      []string{user},
		"password":  []string{password},
		"autologin": []string{"1"},
		"enter":     []string{"Sign in"},
	}

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&options)
	client := &http.Client{
		Jar: jar,
	}
	resp, _ := client.PostForm(loginURL, v)
	c.cookies = resp.Cookies()
}
