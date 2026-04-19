package main

import (
	"context"
	"log"
	"log/slog"

	_ "github.com/RomanenkoViktor080/url_shortener_service/docs"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/adapter/sql/store"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/cache"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/config"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/encoder"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/generator"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/validator"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/repository"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	ctx := context.Background()

	logger := config.InitLogger()
	cfx := config.Mount()

	//validator init
	validator.InitValidatorTranslator()

	//sql connection init
	pool, err := pgxpool.New(ctx, cfx.Dns)
	if err != nil {
		log.Fatal(err)
	}
	dbStore := store.NewStore(pool)
	defer pool.Close()

	slog.Info("connected to the database")

	//redis connection init
	redis := config.NewRedisClient()
	if err = redis.Ping(ctx).Err(); err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}
	defer redis.Close()

	base62Encoder := encoder.NewBase62Encoder()
	hashGenerator := generator.NewGenerator(dbStore, base62Encoder)

	//caches init
	urlCache := cache.NewUrlCache(*redis)
	hashCache := cache.NewHashCache(dbStore, hashGenerator)

	//repositories init
	urlRepository := repository.NewUrlRepository(dbStore, hashCache, urlCache)

	//services init
	urlService := service.NewUrlService(urlRepository)

	//router init
	handler := config.RouterMount(urlService)

	config.InitScheduler(
		urlRepository,
		hashGenerator,
		logger,
	)

	if err = cfx.Run(handler); err != nil {
		log.Fatal(err)
	}
}
