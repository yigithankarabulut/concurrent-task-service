package models

import "context"

type Task struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type TaskJobModel struct {
	ID          uint            `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	JOB         string          `json:"-"`
	Context     context.Context `json:"-"`
}
