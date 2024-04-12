package dto

import "github.com/yigithankarabulut/ConcurrentTaskService/internal/models"

type SetTaskRequest struct {
	ID          uint   `json:"id" validate:"required,numeric"`
	Title       string `json:"title" validate:"required,max=255,min=3"`
	Description string `json:"description" validate:"required,max=255,min=3"`
	Status      string `json:"status" validate:"required,max=255,min=3"`
}

type GetTaskRequest struct {
	ID uint `json:"id" validate:"required"`
}

type ListTaskRequest struct {
	Status string `json:"status" validate:"required"`
}

type UpdateTaskRequest struct {
	ID          uint   `json:"id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required"`
}

type DeleteTaskRequest struct {
	ID uint `json:"id" validate:"required"`
}

func (l ListTaskRequest) TaskJobMapper(model *models.TaskJobModel) models.TaskJobModel {
	model.Status = l.Status
	return *model
}

func (u UpdateTaskRequest) TaskJobMapper(model *models.TaskJobModel) models.TaskJobModel {
	model.ID = u.ID
	model.Title = u.Title
	model.Description = u.Description
	model.Status = u.Status
	return *model
}

func (s SetTaskRequest) TaskJobMapper(model *models.TaskJobModel) models.TaskJobModel {
	model.ID = s.ID
	model.Title = s.Title
	model.Description = s.Description
	model.Status = s.Status
	return *model
}
