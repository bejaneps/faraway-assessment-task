package runtime

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/log"
)

// WaitSignal method is a runtime utility function that blocks the runtime until
// a signal is received
func WaitSignal() os.Signal {
	return <-getSignalChan()
}

func getSignalChan() chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	return sig
}

func RunUntilSignal(run func() error, stop func(context.Context) error, timeout time.Duration) error {
	sigChan := getSignalChan()
	errSig := make(chan error)

	go func() {
		errSig <- run()
	}()

	select {
	case err := <-errSig:
		return err
	case sig := <-sigChan:
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		log.Info("received signal", log.String("signal", sig.String()))
		if stop != nil {
			err := stop(ctx)
			if err != nil {
				log.Error("could not stop server:", log.StdError(err))
			}
		}
	}

	return nil
}
