package config

import (
	"log"
	"log/slog"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/generator"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/repository"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/task"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/env"
	"github.com/hibiken/asynq"
)

func InitScheduler(
	urlRepository repository.UrlRepository,
	gen generator.Generator,
	logger *slog.Logger) {
	redisOpt := asynq.RedisClientOpt{
		Addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
		Username: env.GetString("REDIS_USERNAME", "user"),
		Password: env.GetString("REDIS_PASSWORD", "password"),
	}
	asynqLogger := NewAsynqLoggerAdapter(logger)
	srv := asynq.NewServer(redisOpt, asynq.Config{
		Logger: asynqLogger,
	})

	mux := asynq.NewServeMux()

	//handlers init
	registry := GetTaskRegistry(urlRepository, gen)
	for _, rItem := range registry {
		mux.HandleFunc(rItem.Pattern, rItem.Handler)
	}

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run asynq server: %v", err)
		}
	}()

	//scheduler init
	sch := asynq.NewScheduler(redisOpt, &asynq.SchedulerOpts{
		Logger: asynqLogger,
	})

	for _, rItem := range registry {
		if _, err := sch.Register(rItem.Task.Schedule, rItem.Task.Fn); err != nil {
			log.Fatalf("could not register task %s: %v", rItem.Task.Name, err)
		}
	}
	go func() {
		if err := sch.Run(); err != nil {
			log.Fatal(err)
		}
	}()
}

func GetTaskRegistry(repo repository.UrlRepository, gen generator.Generator) []task.Handler {
	urlTaskHandler := task.NewUrlTaskHandler(repo)
	urlTask := urlTaskHandler.NewTask()

	hashTaskHandler := task.NewCacheTaskHandler(repo, gen)
	hashTask := hashTaskHandler.NewTask()

	return []task.Handler{
		{
			Pattern: urlTaskHandler.GetPattern(),
			Handler: urlTaskHandler.ProcessTask,
			Task: task.Task{
				Name:     urlTask.Name,
				Schedule: urlTask.Schedule,
				Fn:       urlTask.Fn,
			},
		},
		{
			Pattern: hashTaskHandler.GetPattern(),
			Handler: hashTaskHandler.ProcessTask,
			Task: task.Task{
				Name:     hashTask.Name,
				Schedule: hashTask.Schedule,
				Fn:       hashTask.Fn,
			},
		},
	}
}
