package tasks_transport

import (
	"time"

	"github.com/fffeehh/practice_backend_todoapp/internal/core/domain"
)

type TaskDTOResponse struct {
	ID int `json:"id"`
	Version int `json:"version"`
	Title string `json:"title"`
	Description *string `json:"description"`
	Completed bool `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	CompletedAt *time.Time `jsom:"completed_at"`
	AuthorUserID int `json:"author_user_id"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID: task.ID,
		Version: task.Version,
		Title: task.Title,
		Description: task.Description,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt,
		CompletedAt: task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func taskDTOsFromDomains(tasks []domain.Task) []TaskDTOResponse {
	dtos := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		dtos[i] = taskDTOFromDomain(task)
	}

	return dtos
}
