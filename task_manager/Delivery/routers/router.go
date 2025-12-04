package routers

import (
	"github.com/gin-gonic/gin"
	deliveryControllers "task_manager/Delivery/controllers"
	"task_manager/Infrastructure"
)

func SetupRouter(ctrl *deliveryControllers.Controller, jwtSvc Infrastructure.JWTService) *gin.Engine {
	r := gin.Default()

	// Public
	r.POST("/register", ctrl.Register)
	r.POST("/login", ctrl.Login)

	// Protected
	auth := r.Group("/")
	auth.Use(Infrastructure.AuthRequired(jwtSvc))
	{
		// Read tasks (any authenticated user)
		auth.GET("/tasks", ctrl.GetTasks)
		auth.GET("/tasks/:id", ctrl.GetTaskByID)

		// Admin-only
		auth.POST("/promote", Infrastructure.AdminOnly(), ctrl.Promote)
		auth.POST("/tasks", Infrastructure.AdminOnly(), ctrl.CreateTask)
		auth.PUT("/tasks/:id", Infrastructure.AdminOnly(), ctrl.UpdateTask)
		auth.DELETE("/tasks/:id", Infrastructure.AdminOnly(), ctrl.DeleteTask)
	}

	return r
}
