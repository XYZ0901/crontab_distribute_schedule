package data

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"schedule/initial"
	"strings"
	"sync"
)

var (
	taskPrefix = initial.Cfg.Etcd.TaskPrefix
	killPrefix = initial.Cfg.Etcd.KillPrefix
	lockPrefix = initial.Cfg.Etcd.LockPrefix
	sep        = "$!$"
	TaskChan   = make(chan Task, 10)
	DeleteChan = make(chan Task, 10)
	KillChan   = make(chan Task, 10)
	tmap       taskMap
)

type Task struct {
	TaskName string `json:"task_name,omitempty"`
	Command  string `json:"command,omitempty"`
	Crontab  string `json:"crontab,omitempty"`
}

type taskMap struct {
	sync.RWMutex
	tMap map[string]bool
}

// 监听本机任务 抢占公共任务
func WatchTasks() error {
	ipPrefix := taskPrefix + ip + "/"
	op := clientv3.OpGet(ipPrefix, clientv3.WithPrefix())
	opRes, err := initial.Kv.Do(context.Background(), op)
	if err != nil {
		return err
	}
	for _, kv := range opRes.Get().Kvs {
		v := string(kv.Value)
		TaskChan <- Task{
			TaskName: strings.TrimPrefix(string(kv.Key), ipPrefix),
			Command:  strings.Split(v, sep)[0],
			Crontab:  strings.Split(v, sep)[1],
		}
	}

	go func() {
		watcherRevision := opRes.Get().Header.GetRevision()
		watcherChan := initial.Watcher.Watch(context.Background(), ipPrefix,
			clientv3.WithPrefix(), clientv3.WithRev(watcherRevision))
		for {
			select {
			case wRes := <-watcherChan:
				for _, even := range wRes.Events {
					switch even.Type {
					case clientv3.EventTypePut:
						v := string(even.Kv.Value)
						TaskChan <- Task{
							TaskName: strings.TrimPrefix(string(even.Kv.Key), ipPrefix),
							Command:  strings.Split(v, sep)[0],
							Crontab:  strings.Split(v, sep)[1],
						}
					case clientv3.EventTypeDelete:
						DeleteChan <- Task{
							TaskName: strings.TrimPrefix(string(even.Kv.Key), ipPrefix),
						}
					}
				}
			}
		}
	}()
	return nil
}

func WatchKillTasks() {
	ipPrefix := killPrefix + ip + "/"
	go func() {
		watchChan := initial.Watcher.Watch(context.Background(), ipPrefix, clientv3.WithPrefix())
		for {
			select {
			case wRes := <-watchChan:
				for _, even := range wRes.Events {
					switch even.Type {
					case clientv3.EventTypePut:
						KillChan <- Task{
							TaskName: strings.TrimPrefix(string(even.Kv.Key), ipPrefix),
						}
					default:
					}
				}
			}
		}

	}()
}

func WatchGlobalTasks() error {
	globalPrefix := taskPrefix + "/"
	count := 10
	tmap = taskMap{tMap: make(map[string]bool, 10)}

	op := clientv3.OpGet(globalPrefix, clientv3.WithPrefix())
	opRes, err := initial.Kv.Do(context.Background(), op)
	if err != nil {
		return err
	}
	for _, kv := range opRes.Get().Kvs {
		TaskName := strings.TrimPrefix(string(kv.Key), globalPrefix)
		lockKey := lockPrefix + TaskName

		txn := initial.Kv.Txn(context.Background())
		txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
			Then(clientv3.OpPut(lockKey, ip)).
			Else(clientv3.OpGet(lockKey))
		txnRes, err := txn.Commit()
		if err != nil {
			continue
		}
		if !txnRes.Succeeded {
			continue
		}

		tmap.Lock()
		if len(tmap.tMap) >= count {
			tmap.Unlock()
			break
		}
		tmap.tMap[TaskName] = true
		tmap.Unlock()

		v := string(kv.Value)
		TaskChan <- Task{
			TaskName: TaskName + "$global",
			Command:  strings.Split(v, sep)[0],
			Crontab:  strings.Split(v, sep)[1],
		}
	}

	go func() {
		watcherRevision := opRes.Get().Header.GetRevision()
		watcherChan := initial.Watcher.Watch(context.Background(), globalPrefix,
			clientv3.WithPrefix(), clientv3.WithRev(watcherRevision))
		for {
			select {
			case wRes := <-watcherChan:
				for _, even := range wRes.Events {
					TaskName := strings.TrimPrefix(string(even.Kv.Key), globalPrefix)
					switch even.Type {
					case clientv3.EventTypePut:
						tmap.Lock()
						if _, ok := tmap.tMap[TaskName]; !ok && len(tmap.tMap) >= count {
							tmap.Unlock()
							continue
						}
						tmap.tMap[TaskName] = true
						tmap.Lock()
						v := string(even.Kv.Value)
						TaskChan <- Task{
							TaskName: TaskName + "$global",
							Command:  strings.Split(v, sep)[0],
							Crontab:  strings.Split(v, sep)[1],
						}
					case clientv3.EventTypeDelete:
						tmap.Lock()
						if _, ok := tmap.tMap[TaskName]; !ok {
							tmap.Unlock()
							continue
						}
						delete(tmap.tMap, TaskName)
						tmap.Unlock()
						DeleteChan <- Task{
							TaskName: TaskName + "$global",
						}
					}
				}
			}
		}
	}()
	return nil
}

func WatchGlobalKillTasks() {
	Prefix := killPrefix + "/"
	go func() {
		watchChan := initial.Watcher.Watch(context.Background(), Prefix, clientv3.WithPrefix())
		for {
			select {
			case wRes := <-watchChan:
				for _, even := range wRes.Events {
					switch even.Type {
					case clientv3.EventTypePut:
						TaskName := strings.TrimPrefix(string(even.Kv.Key), Prefix)
						KillChan <- Task{
							TaskName: TaskName,
						}
					default:
					}
				}
			}
		}

	}()
}
