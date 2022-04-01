package main

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
	"os/exec"
	"schedule/worker/data"
)

func main() {
	ch := make(chan struct{})

	taskMap := make(map[string]struct {
		id     int
		cancel context.CancelFunc
	})

	c := cron.New(cron.WithSeconds())
	go c.Run()

	go func() {
		for {
			select {
			case task := <-data.TaskChan:
				oldtask, ok := taskMap[task.TaskName]
				if ok {
					c.Remove(cron.EntryID(oldtask.id))
				}
				ctx, cancel := context.WithCancel(context.Background())
				id, err := c.AddFunc(task.Crontab, func() {
					cmd := exec.CommandContext(ctx, "sh", "-c", task.Command)
					out, err := cmd.Output()

					if err != nil {
						log.Println(
							"[ERROR] task name:", task.TaskName,
							" task command:", task.Command, " err:", err)
					} else {
						log.Println(
							"[INFO] task name:", task.TaskName,
							" task command:", task.Command, " res:", string(out))
					}
				})
				if err != nil {
					break
				}
				taskMap[task.TaskName] = struct {
					id     int
					cancel context.CancelFunc
				}{id: int(id), cancel: cancel}
			case task := <-data.DeleteChan:
				oldtask, ok := taskMap[task.TaskName]
				if ok {
					log.Println(
						"[INFO] Delete task name:", task.TaskName)
					c.Remove(cron.EntryID(oldtask.id))
					delete(taskMap, task.TaskName)
				}
			case task := <-data.KillChan:
				oldtask, ok := taskMap[task.TaskName]
				if ok {
					oldtask.cancel()
				}
			}
		}
	}()

	<-ch
}
