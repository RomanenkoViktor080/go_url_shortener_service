package task

import (
	"context"

	"github.com/hibiken/asynq"
)

type handler interface {
	NewTask() *Task
	ProcessTask(ctx context.Context, t *asynq.Task) error
	GetPattern() string
}

type Handler struct {
	Pattern string
	Handler func(context.Context, *asynq.Task) error
	Task    Task
}
type Task struct {
	Name     string
	Schedule string
	Fn       *asynq.Task
}
