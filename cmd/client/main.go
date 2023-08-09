package main

import (
	"context"
	"time"

	_ "github.com/bejaneps/faraway-assessment-task/internal/pkg/debug" // prints debug info

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/log"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/runtime"
	serviceClient "github.com/bejaneps/faraway-assessment-task/internal/service/client"
	"github.com/bejaneps/faraway-assessment-task/internal/transport/client"
	"github.com/caarlos0/env/v9"
)

type config struct {
	ServerPort        string `env:"SERVER_PORT" envDefault:"5252"`
	ServerHost        string `env:"SERVER_HOST" envDefault:"server"`
	StopClientTimeout int    `env:"STOP_CLIENT_TIMEOUT" envDefault:"10"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err.Error())
	}

	serverAddr := cfg.ServerHost + ":" + cfg.ServerPort

	ctx, cancelFunc := context.WithCancel(context.Background())
	ctx = log.ContextWithAttributes(ctx, log.Attributes{"serverAddr": serverAddr})

	service := serviceClient.New(0)
	c := client.New(service)
	err := runtime.RunUntilSignal(
		func() error { // start func
			return c.Dial(ctx, serverAddr)
		},
		func(ctx context.Context) error { // stop func
			cancelFunc()
			return nil
		}, time.Duration(cfg.StopClientTimeout)*time.Second,
	)
	if err != nil {
		log.Fatal(err.Error())
	}
}
