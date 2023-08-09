package server

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/log"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/transport"
	repoServer "github.com/bejaneps/faraway-assessment-task/internal/repository/server"
)

const (
	maxRandomNumber = 2_000_000 // it takes quite some time to find a number in this range

	incorectGuess = "incorrect guess"
)

// ErrNumberNotFound is used after client fails to guess number
var ErrNumberNotFound = errors.New("number wasn't found")

// Challenger is wrapper for POW challenge-respond algorithm
type Challenger interface {
	// Challenge challenges client with a specific algorithm to solve
	Challenge(ctx context.Context, rw transport.ReadWriter) error
}

type Service struct {
	maxRandomNumberBig *big.Int
	quoter             repoServer.Quoter
}

func New(quoter repoServer.Quoter) *Service {
	return &Service{
		maxRandomNumberBig: big.NewInt(maxRandomNumber),
		quoter:             quoter,
	}
}

func (s *Service) Challenge(ctx context.Context, rw transport.ReadWriter) error {
	// generate random number and send it to client, so he starts guessing
	number, hash, err := randomNumberGenerator(s.maxRandomNumberBig)
	if err != nil {
		return fmt.Errorf("failed to generate random number and calculate it's hash: %w", err)
	}
	log.FromContext(ctx).Debug("sending random number to client", log.Int64("number", number.Int64()), log.String("hash", hash))
	if err := rw.WriteMessage(hash); err != nil {
		return fmt.Errorf("failed to send random number: %w", err)
	}

	// read guess number
	msg, err := rw.ReadMessage()
	if err != nil {
		return fmt.Errorf("failed to read guess number: %w", err)
	}
	log.FromContext(ctx).Debug("read guess number from client", log.String("number", msg))
	guess, err := strconv.Atoi(msg)
	if err != nil {
		return fmt.Errorf("failed to convert to int guess message: %w", err)
	}

	// if value was guessed wrong, send incorect guess message
	if int64(guess) != number.Int64() {
		if err := rw.WriteMessage(incorectGuess); err != nil {
			return fmt.Errorf("failed to send incorrect guess message: %w", err)
		}

		return ErrNumberNotFound
	}

	// if value was guessed correct, get value from repo and send it
	quote, err := s.quoter.Quote(ctx)
	if err != nil {
		return fmt.Errorf("failed to get random quote: %w", err)
	}
	if err := rw.WriteMessage(quote); err != nil {
		return fmt.Errorf("failed to send quote: %w", err)
	}

	return nil
}

var randomNumberGenerator = getRandomNumberAndHash

func getRandomNumberAndHash(max *big.Int) (*big.Int, string, error) {
	number, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate random number: %w", err)
	}

	sum := sha256.Sum256([]byte(fmt.Sprintf("%d", number.Int64())))

	return number, fmt.Sprintf("%x", sum[:]), nil
}
