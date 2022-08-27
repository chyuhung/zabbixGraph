package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	Server      string
	Port        string
	User        string
	Password    string
	ApiRpcURL   string
	GraphURL    string
	LoginURL    string
	DownloadDir string
)

func loadData(file *ini.File) {
	Server = file.Section("zabbix").Key("Server").MustString("10.191.101.101")
	Port = file.Section("zabbix").Key("Port").MustString(":80")
	User = file.Section("zabbix").Key("User").MustString("Admin")
	Password = file.Section("zabbix").Key("Password").MustString("zabbix")
	ApiRpcURL = "http://" + Server + Port + "/zabbix/api_jsonrpc.php"
	GraphURL = "http://" + Server + Port + "/zabbix/chart2.php"
	LoginURL = "http://" + Server + Port + "/zabbix/index.php"
	DownloadDir = file.Section("config").Key("DownloadDir").MustString("images")
}
func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误:", err)
	}
	loadData(file)
}
