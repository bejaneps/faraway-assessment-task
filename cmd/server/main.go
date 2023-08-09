package main

import (
	"context"
	"time"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/db"
	_ "github.com/bejaneps/faraway-assessment-task/internal/pkg/debug" // prints debug info

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/cache"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/log"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/runtime"
	repoServer "github.com/bejaneps/faraway-assessment-task/internal/repository/server"
	serviceServer "github.com/bejaneps/faraway-assessment-task/internal/service/server"
	"github.com/bejaneps/faraway-assessment-task/internal/transport/server"
	"github.com/caarlos0/env/v9"
)

type config struct {
	CacheConfig cache.Config
	DBConfig    db.Config

	ServerPort        string `env:"SERVER_PORT" envDefault:"5252"`
	StopServerTimeout int    `env:"STOP_SERVER_TIMEOUT" envDefault:"10"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	ctx = log.ContextWithAttributes(ctx, log.Attributes{"serverPort": cfg.ServerPort})

	cacheClient, err := cache.New(cfg.CacheConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	dbClient, err := db.New(cfg.DBConfig)
	if err != nil {
		cacheClient.Close()
		log.Fatal(err.Error())
	}

	repo := repoServer.New(cacheClient, dbClient)
	service := serviceServer.New(repo)
	s, err := server.New(ctx, service, "0.0.0.0:"+cfg.ServerPort)
	if err != nil {
		logIfError(dbClient.Close())
		logIfError(cacheClient.Close())
		log.Fatal(err.Error())
	}

	if err := runtime.RunUntilSignal(
		func() error { // start func
			return s.ListenAndAccept(ctx)
		},
		func(ctx context.Context) error { // stop func
			cancelFunc()
			logIfError(cacheClient.Close())
			logIfError(dbClient.Close())
			return s.Close()
		}, time.Duration(cfg.StopServerTimeout)*time.Second,
	); err != nil {
		log.Fatal(err.Error())
	}
}

func logIfError(err error) {
	if err != nil {
		log.Error(err.Error())
	}
}
