package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"task_manager/data"
	"task_manager/models"
)

type TaskController struct {
	Service *data.TaskService
}

// creates a new controller instance
func NewTaskController(service *data.TaskService) *TaskController {
	return &TaskController{Service: service}
}

// GET /tasks
func (c *TaskController) GetTasks(ctx *gin.Context) {
	tasks, err := c.Service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch tasks",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": tasks,
	})
}

// GET /tasks/:id
func (c *TaskController) GetTaskByID(ctx *gin.Context) {
	// ID is now a string (MongoDB _id)
	id := ctx.Param("id")

	task, err := c.Service.GetByID(id)
	if err != nil {
		if err == data.ErrNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch task",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": task,
	})
}

// POST /tasks
func (c *TaskController) CreateTask(ctx *gin.Context) {
	var input models.Task

	// Bind JSON and validate required fields (Title)
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input: " + err.Error(),
		})
		return
	}

	created, err := c.Service.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create task",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": created,
	})
}

// PUT /tasks/:id
func (c *TaskController) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")

	var input models.Task
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input: " + err.Error(),
		})
		return
	}

	updated, err := c.Service.Update(id, input)
	if err != nil {
		if err == data.ErrNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update task",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": updated,
	})
}

// DELETE /tasks/:id
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.Service.Delete(id); err != nil {
		if err == data.ErrNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "task not found",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete task",
			})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
