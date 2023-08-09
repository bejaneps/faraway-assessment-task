package client

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/avast/retry-go"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/log"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/transport"
	svcClient "github.com/bejaneps/faraway-assessment-task/internal/service/client"
)

const (
	connectionRetryAttempts = 10
	connectionRetyDelay     = 1 * time.Second
)

type Client struct {
	responder svcClient.Responder
}

func New(responder svcClient.Responder) *Client {
	return &Client{
		responder: responder,
	}
}

func (c *Client) Dial(ctx context.Context, address string) error {
	var (
		conn net.Conn
		err  error
	)
	err = retry.Do(func() error {
		conn, err = net.Dial("tcp", address)
		if err != nil {
			return err
		}

		return nil
	}, retry.Attempts(connectionRetryAttempts), retry.Delay(connectionRetyDelay), retry.DelayType(retry.FixedDelay))
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}
	defer conn.Close()

	log.FromContext(ctx).Info("connected successfully")

	if err := c.responder.Respond(ctx, transport.New(conn)); err != nil {
		return fmt.Errorf("failed to solve server challenge: %w", err)
	}

	return nil
}
