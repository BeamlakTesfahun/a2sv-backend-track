package router

import (
	"github.com/gin-gonic/gin"
	"task_manager/controllers"
)

func SetupRouter(taskController *controllers.TaskController) *gin.Engine {
	r := gin.Default()

	// grouped routes
	taskRoutes := r.Group("/tasks")
	{
		taskRoutes.GET("", taskController.GetTasks)
		taskRoutes.GET("/:id", taskController.GetTaskByID)
		taskRoutes.POST("", taskController.CreateTask)
		taskRoutes.PUT("/:id", taskController.UpdateTask)
		taskRoutes.DELETE("/:id", taskController.DeleteTask)
	}

	return r
}
