package data

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net"
	"os"
	"schedule/initial"
)

var (
	ip = "127.0.0.1"
)

func init() {
	setLocalIp()
	registerToEtcd()
	_ = WatchTasks()
	WatchKillTasks()
	_ = WatchGlobalTasks()
}

func setLocalIp() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
				ip = ipnet.IP.String()
			}
		}
	}
}

func registerToEtcd() {
	op := clientv3.OpPut(initial.Cfg.Etcd.MachinePrefix+ip, initial.Cfg.Worker.Name)
	_, err := initial.Kv.Do(context.Background(), op)
	if err != nil {
		log.Fatalln("[ERROR] [data.registerToEtcd] Kv.do Error", err)
		return
	}
}
