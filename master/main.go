package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"schedule/initial"
	"schedule/master/api"
	"strconv"
	"time"
)

func main() {
	r := gin.Default()
	r.Use()
	r.POST("/task/create", api.CreateTask)
	r.POST("/task/delete", api.DeleteTask)
	r.GET("/task/tasks", api.GetTasks)
	r.POST("/task/kill", api.KillTask)
	r.GET("/task/machines", api.GetMachines)

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(initial.Cfg.Master.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(initial.Cfg.Master.ReadTimeOut) * time.Millisecond,
		WriteTimeout: time.Duration(initial.Cfg.Master.WriteTimeOut) * time.Millisecond,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
