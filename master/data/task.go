package data

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"schedule/initial"
	"strings"
)

var (
	taskPrefix = initial.Cfg.Etcd.TaskPrefix
	killPrefix = initial.Cfg.Etcd.KillPrefix
	sep        = "$!$"
)

type Task struct {
	TaskName  string `json:"task_name,omitempty"`
	MachineIp string `json:"machine_ip"`
	Command   string `json:"command,omitempty"`
	Crontab   string `json:"crontab,omitempty"`
}

func (t *Task) Create() error {
	op := clientv3.OpPut(taskPrefix+t.MachineIp+"/"+t.TaskName, t.Command+sep+t.Crontab)
	_, err := initial.Kv.Do(context.Background(), op)
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) Delete() error {
	op := clientv3.OpDelete(taskPrefix + t.MachineIp + "/" + t.TaskName)
	_, err := initial.Kv.Do(context.Background(), op)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllTasks() error {
	op := clientv3.OpDelete(taskPrefix, clientv3.WithPrefix())
	_, err := initial.Kv.Do(context.Background(), op)
	if err != nil {
		return err
	}
	return nil
}

func GetTaskList() ([]Task, error) {
	tasks := []Task{}
	op := clientv3.OpGet(taskPrefix, clientv3.WithPrefix())
	res, err := initial.Kv.Do(context.Background(), op)
	if err != nil {
		return tasks, err
	}
	for _, kv := range res.Get().Kvs {
		v := string(kv.Value)
		tasks = append(tasks, Task{
			TaskName: strings.TrimPrefix(string(kv.Key), taskPrefix),
			Command:  strings.Split(v, sep)[0],
			Crontab:  strings.Split(v, sep)[1],
		})
	}

	return tasks, nil
}

func (t *Task) Kill() error {
	lgRes, err := initial.Lease.Grant(context.Background(), 1)
	if err != nil {
		return err
	}
	op := clientv3.OpPut(killPrefix+t.MachineIp+"/"+t.TaskName, "", clientv3.WithLease(lgRes.ID))
	_, err = initial.Kv.Do(context.Background(), op)
	if err != nil {
		return err
	}
	return nil
}
