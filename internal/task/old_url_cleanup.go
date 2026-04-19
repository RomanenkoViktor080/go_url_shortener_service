package task

import (
	"context"
	"time"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/adapter/sql/store"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/repository"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/env"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
)

type urlTaskHandler struct {
	repo      repository.UrlRepository
	pattern   string
	limit     int32
	daysLimit int
}

func NewUrlTaskHandler(repo repository.UrlRepository) handler {
	return &urlTaskHandler{
		repo:      repo,
		pattern:   "url:cleanup",
		limit:     int32(env.GetInt("TASK_URL_CLEANUP_LIMIT", 10000)),
		daysLimit: env.GetInt("TASK_URL_CLEANUP_AFTER_DAYS", 30),
	}
}

func (h *urlTaskHandler) NewTask() *Task {
	return &Task{
		Name:     env.GetString("TASK_URL_CLEANUP_NAME", "URL Cleanup"),
		Schedule: env.GetString("TASK_URL_CLEANUP_SCHEDULE", "@every 1m"),
		Fn:       asynq.NewTask(h.pattern, nil),
	}
}

func (h *urlTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	deleteBefore := time.Now().UTC().AddDate(0, 0, -h.daysLimit)
	return h.repo.DeleteShortUrlBeforeCreatedAt(ctx, store.DeleteShortUrlBeforeCreatedAtParams{
		CreatedAt: pgtype.Timestamp{
			Time:  deleteBefore,
			Valid: true,
		},
		Limit: h.limit,
	})
}
func (h *urlTaskHandler) GetPattern() string {
	return h.pattern
}
