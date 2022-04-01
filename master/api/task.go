package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"schedule/master/data"
)

func CreateTask(c *gin.Context) {
	task := &data.Task{}
	err := c.ShouldBind(task)
	if err != nil {
		log.Println("[ERROR] [api.CreateTask] c.ShouldBind error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}
	err = task.Create()
	if err != nil {
		log.Println("[ERROR] [api.CreateTask] task.Create error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Create task success",
	})
}

func DeleteTask(c *gin.Context) {
	task := &data.Task{}
	err := c.ShouldBind(task)
	if err != nil {
		log.Println("[ERROR] [api.DeleteTask] c.ShouldBind error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}
	err = task.Delete()
	if err != nil {
		log.Println("[ERROR] [api.DeleteTask] task.Delete error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete task success",
	})
}

func GetTasks(c *gin.Context) {
	//pn_s := c.Query("pn")
	//pSize_s := c.Query("pSize")
	//pn, err := strconv.Atoi(pn_s)
	//if err != nil {
	//	log.Println("[ERROR] [api.GetTasks] pn is not int")
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"message": err.Error(),
	//	})
	//	c.Abort()
	//	return
	//}
	//pSize, err := strconv.Atoi(pSize_s)
	//if err != nil {
	//	log.Println("[ERROR] [api.GetTasks] pSize is not int")
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"message": err.Error(),
	//	})
	//	c.Abort()
	//	return
	//}
	tasks, err := data.GetTaskList()
	if err != nil {
		log.Println("[ERROR] [api.GetTasks] data.GetTaskList error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func KillTask(c *gin.Context) {
	task := &data.Task{}
	err := c.ShouldBind(task)
	if err != nil {
		log.Println("[ERROR] [api.KillTask] c.ShouldBind error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}
	err = task.Kill()
	if err != nil {
		log.Println("[ERROR] [api.KillTask] task.Kill error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Kill task success",
	})
}

func GetMachines(c *gin.Context) {
	machines, err := data.GetMachines()
	if err != nil {
		log.Println("[ERROR] [api.GetMachines] data.GetMachines error:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, machines)
}
