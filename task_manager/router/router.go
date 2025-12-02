package router

import (
	"github.com/gin-gonic/gin"
	"task_manager/controllers"
	"task_manager/middleware"
)

func SetupRouter(userController *controllers.UserController, taskController *controllers.TaskController) *gin.Engine {
	r := gin.Default()

	// Public
	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)

	// Authenticated
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		// Any authenticated user can read tasks
		auth.GET("/tasks", taskController.GetTasks)
		auth.GET("/tasks/:id", taskController.GetTaskByID)

		// Admin-only actions
		auth.POST("/promote", middleware.AdminOnly(), userController.Promote)

		auth.POST("/tasks", middleware.AdminOnly(), taskController.CreateTask)
		auth.PUT("/tasks/:id", middleware.AdminOnly(), taskController.UpdateTask)
		auth.DELETE("/tasks/:id", middleware.AdminOnly(), taskController.DeleteTask)
	}

	return r
}
