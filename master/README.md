## Master
### Task Manage
HTTP Api of task's CURD
### Task Log
search task's history log
### Task Control
forced end of task

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