package initial

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func etcdInit() {
	var err error
	EtcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   Cfg.Etcd.Endpoints,
		DialTimeout: time.Duration(Cfg.Etcd.DialTimeOut) * time.Millisecond,
	})
	if err != nil {
		log.Fatalln(err)
	}
	Kv = clientv3.NewKV(EtcdCli)
	Lease = clientv3.NewLease(EtcdCli)
	Watcher = clientv3.NewWatcher(EtcdCli)
}
