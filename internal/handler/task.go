package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ahmednurovic/task-manager-api/internal/model"
)

type TaskService interface {
	CreateTask(ctx *gin.Context, task *model.Task) error
	GetTasks(ctx *gin.Context, userID int64) ([]*model.Task, error)
	UpdateTask(ctx *gin.Context, task *model.Task) error
	DeleteTask(ctx *gin.Context, taskID int64, userID int64) error
}

type TaskHandler struct {
	service TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task for the given user
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body model.Task true "Task object"
// @Success 201 {object} model.Task
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task model.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateTask(c, &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTasks godoc
// @Summary Get tasks by user ID
// @Description Get a list of tasks for the given user
// @Tags tasks
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {array} model.Task
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	tasks, err := h.service.GetTasks(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body model.Task true "Updated task object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var task model.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateTask(c, &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param user_id query int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	userIDStr := c.Query("user_id")

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.service.DeleteTask(c, taskID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
