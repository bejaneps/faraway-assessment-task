package client

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/log"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/transport"
)

const (
	msgNumbersLen       = 2
	msgNumbersSeparator = "-"

	incorectGuess = "incorrect guess"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Respond(ctx context.Context, conn net.Conn) error {
	reader := bufio.NewReader(conn)

	// read first message and get number range
	msg, err := transport.ReadMessage(reader)
	if err != nil {
		return fmt.Errorf("failed to read first message: %w", err)
	}

	minMax := strings.Split(msg, msgNumbersSeparator)
	if len(minMax) < msgNumbersLen {
		return fmt.Errorf("incorrect message: %s", msg)
	}

	minStr, maxStr := minMax[0], minMax[1]
	min, err := strconv.Atoi(minStr)
	if err != nil {
		return fmt.Errorf("failed to parse number range min %s: %w", minStr, err)
	}
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		return fmt.Errorf("failed to parse number range max %s: %w", maxStr, err)
	}

	log.FromContext(ctx).Info(
		"received number range from server",
		log.Int("min", min), log.Int("max", max),
	)

	for i := min; i < max; i++ {
		if err := transport.WriteMessage(conn, fmt.Sprintf("%d", i)); err != nil {
			return fmt.Errorf("failed to send guess number: %w", err)
		}

		msg, err := transport.ReadMessage(reader)
		if err != nil {
			return err
		}
		if msg != "" && msg != incorectGuess {
			fmt.Println(msg)
			return nil
		}
	}

	return nil
}
