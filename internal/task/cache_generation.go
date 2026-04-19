package task

import (
	"context"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/generator"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/repository"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/env"
	"github.com/hibiken/asynq"
)

type cacheTaskHandler struct {
	repo     repository.UrlRepository
	gen      generator.Generator
	pattern  string
	quantity int64
	limit    int64
}

func NewCacheTaskHandler(
	repo repository.UrlRepository,
	gen generator.Generator,
) handler {
	return &cacheTaskHandler{
		repo:     repo,
		gen:      gen,
		pattern:  "cache:generation",
		quantity: int64(env.GetInt("TASK_CACHE_GENERATION_QUANTITY", 10000)),
		limit:    int64(env.GetInt("TASK_CACHE_GENERATION_LIMIT", 100000)),
	}
}

func (h *cacheTaskHandler) NewTask() *Task {
	return &Task{
		Name:     env.GetString("TASK_CACHE_GENERATION_NAME", "Cache Generation"),
		Schedule: env.GetString("TASK_CACHE_GENERATION_SCHEDULE", "@every 1m"),
		Fn:       asynq.NewTask(h.pattern, nil),
	}
}

func (h *cacheTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	count, err := h.repo.CountCaches(ctx)
	if err != nil {
		return err
	}
	if count < h.limit {
		_, err = h.gen.GenerateBatchAndSave(ctx, h.quantity)
	}
	return err
}

func (h *cacheTaskHandler) GetPattern() string {
	return h.pattern
}
