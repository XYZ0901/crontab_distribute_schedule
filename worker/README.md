## Worker
### Task Sync
watch the path of tasks in etcd
### Task Schedule
based on the crontab to touch of task
### Task Execute
the distributed lock based on etcd
### log store
store log to db_log
### struct
#### the struct of store to etcd
manage
```
/crontab/tasks/${taskName} -> "${jsonType}"
```
control
```
/crontab/kill/${taskName} -> ""
```
#### the struct of store to db_log
```
{
	taskName,
	command,
	err,
	output,
	startTime,
	endTime,
}
```