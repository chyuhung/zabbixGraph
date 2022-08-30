package zabbix

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"os"
	"strings"
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

func setCookie(name string, value string) *http.Cookie {
	var cookie = &http.Cookie{
		Name:  replaceNR(name),
		Value: replaceNR(value),
		//Path:     "/zabbix",
		//HttpOnly: true,
	}
	return cookie
}

// zabbix未提供获取图形的api（只能获取图形数据），模仿登录获取数据渲染后的图形
func (c *Client) getCookies() {
	req := gorequest.New()
	res, _, errs := req.Get(utils.LoginURL).
		Set("User-Agent", utils.UserAgent).
		Query("autologin=1&enter=Sign%20in").
		Query("name="+utils.User).
		Query("password="+utils.Password).
		AppendHeader("Refer", utils.LoginURL).End()
	for _, err := range errs {
		fmt.Println("登录失败:", err)
	}
	fmt.Println(res.Header)
	cookies := res.Request.Cookies()
	//部分系统编译后无法正常获取cookies值，需手动构造
	if len(cookies) == 0 {
		fmt.Println("未成功获取到cookies,请根据头文件进行构造")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Set-Cookie Name:")
		name, _ := reader.ReadString('\n')
		fmt.Print("Set-Cookie Value:")
		value, _ := reader.ReadString('\n')
		fmt.Println(setCookie(name, value))
		c.cookies = append(c.cookies, setCookie(name, value))
		fmt.Println("构造cookies:", c.cookies)
	}
	c.cookies = cookies
}
