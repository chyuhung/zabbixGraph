package main

import (
	"fmt"
	"zabbixGraph/zabbix"
)

func main() {
	client := zabbix.New()
	hostId := zabbix.GetHostID(client, "127.0.0.1")
	fmt.Println(hostId)
}
