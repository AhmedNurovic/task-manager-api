package service

import (
	"errors"

	"github.com/ahmednurovic/task-manager-api/internal/model"
	"github.com/ahmednurovic/task-manager-api/internal/repository"
	"github.com/gin-gonic/gin"
)

type TaskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) CreateTask(ctx *gin.Context, task *model.Task) error {
	if task.Title == "" {
		return errors.New("title is required")
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return err
	}

	return nil
}

func (s *TaskService) GetTasks(ctx *gin.Context, userID int64) ([]*model.Task, error) {
	tasks, err := s.taskRepo.GetAllForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx *gin.Context, task *model.Task) error {
	existingTask, err := s.taskRepo.GetByID(ctx, int64(task.ID))
	if err != nil {
		return ErrTaskNotFound
	}

	existingTask.Title = task.Title
	existingTask.Status = task.Status

	if err := s.taskRepo.Update(ctx, existingTask); err != nil {
		return err
	}

	return nil
}

func (s *TaskService) DeleteTask(ctx *gin.Context, taskID int64, userID int64) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return ErrTaskNotFound
	}

	if task.UserID != uint(userID) {
		return errors.New("unauthorized to delete this task")
	}

	if err := s.taskRepo.Delete(ctx, taskID, userID); err != nil {
		return err
	}

	return nil
}
