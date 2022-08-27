package main

import (
	"fmt"
	"zabbixGraph/zabbix"
)

func main() {
	client := zabbix.New()
	hostId := zabbix.GetHostID(client, "127.0.0.1")
	fmt.Println(hostId)
	graphID := zabbix.GetGraphID(client, hostId, []string{"CPU utilization"})
	fmt.Println(graphID)
	zabbix.GetGraph(client, "cpuinfo-127.0.0.1", "525", "now-1h", "now", "1000", "800")
}
