package data

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"schedule/initial"
	"strings"
)

type Machine struct {
	Name string // value
	Ip   string // key
}

var (
	machinePrefix = initial.Cfg.Etcd.MachinePrefix
)

func GetMachines() ([]Machine, error) {
	machines := []Machine{}
	op := clientv3.OpGet(machinePrefix, clientv3.WithPrefix())
	opRes, err := initial.Kv.Do(context.Background(), op)
	if err != nil {
		return machines, err
	}
	for _, kv := range opRes.Get().Kvs {
		machines = append(machines, Machine{
			Name: string(kv.Value),
			Ip:   strings.TrimPrefix(string(kv.Key), machinePrefix),
		})
	}
	return machines, nil
}
