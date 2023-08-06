package server

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net"
	"strconv"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/transport"
)

const (
	numberRangeUpperBound = math.MaxInt64
	numberRangeMax        = 5 // could be dynamically changed depending on server load

	incorectGuess = "incorrect guess"
)

// ErrNumberNotFound is used after client fails to guess number in max attempts
var ErrNumberNotFound = errors.New("number wasn't found in max attempts")

type Service struct {
	numberRangeUpperBoundBig *big.Int
	numberRangeMaxBig        *big.Int
}

func New() *Service {
	return &Service{
		numberRangeUpperBoundBig: big.NewInt(numberRangeUpperBound),
		numberRangeMaxBig:        big.NewInt(numberRangeMax),
	}
}

func (s *Service) Challenge(conn net.Conn) error {
	// generate random number and send it to client, so he starts guessing
	number, min, max, err := s.getRandomNumberAndRange()
	if err != nil {
		return fmt.Errorf("failed to generate random number: %w", err)
	}
	numberChallenge := prepareNumberRangeChallenge(min, max)
	err = transport.WriteMessage(conn, numberChallenge)
	if err != nil {
		return fmt.Errorf("failed to send initial number range challenge: %w", err)
	}

	// read guess number messages, until number is guessed correctly
	reader := bufio.NewReader(conn)
	for i := 0; i < numberRangeMax; i++ {
		msg, err := transport.ReadMessage(reader)
		if err != nil {
			return fmt.Errorf("failed to read guess number: %w", err)
		}

		guess, err := strconv.Atoi(msg)
		if err != nil {
			return fmt.Errorf("failed to convert to int guess message: %w", err)
		}

		if int64(guess) == number.Int64() {
			return nil
		}

		err = transport.WriteMessage(conn, incorectGuess)
		if err != nil {
			return fmt.Errorf("failed to send incorrect message: %w", err)
		}
	}

	return ErrNumberNotFound
}

// getRandomNumberAndRange generates a random number in range of n and n+5
// where n is any number between 0 and max.Int64, it also returns min and max values
// which are n and n+5
func (s *Service) getRandomNumberAndRange() (number *big.Int, min, max int64, err error) {
	number, err = rand.Int(rand.Reader, s.numberRangeUpperBoundBig)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to generate random number: %w", err)
	}

	min, max = number.Int64(), number.Int64()+numberRangeMax

	numberInRange, err := rand.Int(rand.Reader, s.numberRangeMaxBig)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to generate random number range: %w", err)
	}

	number.Add(number, numberInRange)

	return number, min, max, nil
}

func prepareNumberRangeChallenge(min, max int64) string {
	return fmt.Sprintf("%d-%d", min, max)
}
