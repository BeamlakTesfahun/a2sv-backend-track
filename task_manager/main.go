package main

import (
	"log"

	"task_manager/controllers"
	"task_manager/data"
	"task_manager/router"
)

func main() {
	// in memory storage
	taskService := data.NewTaskService()

	// controller
	taskController := controllers.NewTaskController(taskService)

	// router
	r := router.SetupRouter(taskController)

	// run server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
