package initial

import clientv3 "go.etcd.io/etcd/client/v3"

type masterCfg struct {
	Port         int `yaml:"port"`
	ReadTimeOut  int `yaml:"read_time_out"`
	WriteTimeOut int `yaml:"write_time_out"`
}

type workerCfg struct {
	Name string `yaml:"name"`
}

type etcdCfg struct {
	Endpoints     []string `yaml:"endpoints"`
	DialTimeOut   int      `yaml:"dial_time_out"`
	TaskPrefix    string   `yaml:"task_prefix"`
	KillPrefix    string   `yaml:"kill_prefix"`
	MachinePrefix string   `yaml:"machine_prefix"`
	LockPrefix    string   `yaml:"lock_prefix"`
}

type config struct {
	Master masterCfg `yaml:"master"`
	Worker workerCfg `yaml:"worker"`
	Etcd   etcdCfg   `yaml:"etcd"`
}

var (
	Cfg     config
	EtcdCli *clientv3.Client
	Kv      clientv3.KV
	Lease   clientv3.Lease
	Watcher clientv3.Watcher
)

func init() {
	configInit()
	etcdInit()
}
