package controllers

import (
	"net/http"
	"strconv"

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
	tasks := c.Service.GetAll()
	ctx.JSON(http.StatusOK, gin.H{
		"data": tasks,
	})
}

// GET /tasks/:id
func (c *TaskController) GetTaskByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task ID",
		})
		return
	}

	task, err := c.Service.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "task not found",
		})
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


	created := c.Service.Create(input)
	ctx.JSON(http.StatusCreated, gin.H{
		"data": created,
	})
}

// PUT /tasks/:id
func (c *TaskController) UpdateTask(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task ID",
		})
		return
	}

	var input models.Task
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input: " + err.Error(),
		})
		return
	}

	updated, err := c.Service.Update(id, input)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "task not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": updated,
	})
}

// DELETE /tasks/:id
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task ID",
		})
		return
	}

	if err := c.Service.Delete(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "task not found",
		})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
