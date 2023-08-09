package client

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/log"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/transport"
)

const (
	maxRandomNumber = 2_000_000 // it takes quite some time to find a number in this range

	incorectGuess = "incorrect guess"
)

// ErrNumberNotFound is used after client fails to guess number
var ErrNumberNotFound = errors.New("number wasn't found")

// Responder is wrapper for POW challenge-respond algorithm
type Responder interface {
	// Respond tries to solve challenge from server
	Respond(ctx context.Context, rw transport.ReadWriter) error
}

type Service struct {
	upperNumberSearchLimit int64
}

func New(upperNumberSearchLimit int64) *Service {
	if upperNumberSearchLimit == 0 {
		upperNumberSearchLimit = maxRandomNumber
	}

	return &Service{
		upperNumberSearchLimit: upperNumberSearchLimit,
	}
}

func (s *Service) Respond(ctx context.Context, rw transport.ReadWriter) error {
	// read hash of number
	msg, err := rw.ReadMessage()
	if err != nil {
		return fmt.Errorf("failed to read hash: %w", err)
	}
	log.FromContext(ctx).Debug("received hash from server", log.String("hash", msg))

	// iterate over all numbers until max and check if hash is equal to server one
	number := 0
	for i := 0; i < int(s.upperNumberSearchLimit); i++ {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d", i))))

		if hash == msg {
			number = i
			break
		}
	}
	if number == 0 {
		return ErrNumberNotFound
	}

	// if hash was equal, then send the number to server
	log.FromContext(ctx).Debug("sending guess number to server", log.Int("number", number))
	if err := rw.WriteMessage(fmt.Sprintf("%d", number)); err != nil {
		return fmt.Errorf("failed to send guess number: %w", err)
	}

	// check if number was guessed correctly
	msg, err = rw.ReadMessage()
	if err != nil {
		return fmt.Errorf("failed to read quote message: %w", err)
	}
	if msg == "" || msg == incorectGuess {
		return fmt.Errorf("failed to guess number, msg from server: %s", msg)
	}
	log.FromContext(ctx).Info("received quote from server", log.String("quote", msg))

	return nil
}
