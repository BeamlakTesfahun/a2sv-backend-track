package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"task_manager/Domain"
	"task_manager/Repositories"
	"task_manager/Usecases"
)

type Controller struct {
	TaskUC Usecases.TaskUsecases
	UserUC Usecases.UserUsecases
}

func NewController(taskUC Usecases.TaskUsecases, userUC Usecases.UserUsecases) *Controller {
	return &Controller{TaskUC: taskUC, UserUC: userUC}
}

type authRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type promoteRequest struct {
	Username string `json:"username" binding:"required"`
}

func (ctr *Controller) Register(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	user, err := ctr.UserUC.Register(req.Username, req.Password)
	if err != nil {
		if err == Repositories.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func (ctr *Controller) Login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	user, token, err := ctr.UserUC.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"access_token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"role":     user.Role,
			},
		},
	})
}

func (ctr *Controller) Promote(c *gin.Context) {
	var req promoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	user, err := ctr.UserUC.Promote(req.Username)
	if err != nil {
		if err == Repositories.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to promote user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"username": user.Username, "role": user.Role}})
}

func (ctr *Controller) GetTasks(c *gin.Context) {
	tasks, err := ctr.TaskUC.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func (ctr *Controller) GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	task, err := ctr.TaskUC.GetTaskByID(id)
	if err != nil {
		if err == Repositories.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func (ctr *Controller) CreateTask(c *gin.Context) {
	var input Domain.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	created, err := ctr.TaskUC.CreateTask(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

func (ctr *Controller) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var input Domain.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	updated, err := ctr.TaskUC.UpdateTask(id, input)
	if err != nil {
		if err == Repositories.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

func (ctr *Controller) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	if err := ctr.TaskUC.DeleteTask(id); err != nil {
		if err == Repositories.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}

	c.Status(http.StatusNoContent)
}
