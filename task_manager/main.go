package main
import "fmt"


import (
	"log"

	"task_manager/controllers"
	"task_manager/data"
	"task_manager/router"
)

func main() {
	// Initialize MongoDB and get the tasks collection
	collection, err := data.ConnectMongo()

	fmt.Print(err)
	if err != nil {
		log.Fatal("failed to initialize MongoDB:", err)
	}

	// Initialize service and controller
	taskService := data.NewTaskService(collection)
	taskController := controllers.NewTaskController(taskService)

	// router
	r := router.SetupRouter(taskController)

	// run server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
