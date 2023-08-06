package main

import (
	"context"
	"time"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/log"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/runtime"
	"github.com/bejaneps/faraway-assessment-task/internal/server"
	serviceServer "github.com/bejaneps/faraway-assessment-task/internal/service/server"
	"github.com/caarlos0/env/v9"
)

type config struct {
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

	service := serviceServer.New()
	s := server.New(ctx, service, "0.0.0.0:"+cfg.ServerPort)
	err := runtime.RunUntilSignal(
		func() error { // start func
			return s.ListenAndAccept(ctx)
		},
		func(ctx context.Context) error { // stop func
			cancelFunc()
			return s.Close()
		}, time.Duration(cfg.StopServerTimeout)*time.Second,
	)
	if err != nil {
		log.Fatal(err.Error())
	}
}
