package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
	"strings"
)

var (
	Server        string
	Port          string
	User          string
	Password      string
	ApiRpcURL     string
	GraphURL      string
	LoginURL      string
	DownloadDir   string
	Width         string
	Height        string
	TimeFrom      string
	TimeTo        string
	HostsFile     string
	GraphNameList []string
	DownloadSpeed string
)

func loadData(file *ini.File) {
	Server = file.Section("zabbix").Key("Server").MustString("10.191.101.101")
	Port = ":" + file.Section("zabbix").Key("Port").MustString("80")
	User = file.Section("zabbix").Key("User").MustString("Admin")
	Password = file.Section("zabbix").Key("Password").MustString("zabbix")
	ApiRpcURL = "http://" + Server + Port + file.Section("zabbix").Key("ApiURL").MustString("/zabbix/api_jsonrpc.php")
	GraphURL = "http://" + Server + Port + file.Section("zabbix").Key("GraphURL").MustString("/zabbix/chart2.php")
	LoginURL = "http://" + Server + Port + file.Section("zabbix").Key("LoginURL").MustString("/zabbix/index.php")
	DownloadDir = file.Section("config").Key("DownloadDir").MustString("img")
	HostsFile = file.Section("config").Key("HostsFile").MustString("hosts.txt")
	DownloadSpeed = file.Section("config").Key("DownloadSpeed").MustString("0")
	Width = file.Section("graph").Key("Width").MustString("600")
	Height = file.Section("graph").Key("Height").MustString("200")
	TimeFrom = file.Section("graph").Key("TimeFrom").MustString("now-1h")
	TimeTo = file.Section("graph").Key("TimeTo").MustString("now")
	GraphNameList = strings.Split(file.Section("graph").Key("GraphNameList").MustString("CPU utilization"), ",")
}
func printConfig(s []string) {
	fmt.Printf("%s\n", s)
}
func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误:", err)
	}
	loadData(file)
	printConfig([]string{Server, Port, User, Password, ApiRpcURL, GraphURL, LoginURL, DownloadDir, Width, Height, TimeFrom, TimeTo, HostsFile, DownloadSpeed})
	printConfig(GraphNameList)
}
