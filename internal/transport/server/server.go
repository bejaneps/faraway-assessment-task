package server

import (
	"context"
	"fmt"
	"net"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/log"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/transport"
	svcServer "github.com/bejaneps/faraway-assessment-task/internal/service/server"
)

const protocolTCP = "tcp"

// Server is a simple struct to hold dependencies of server
type Server struct {
	listener   net.Listener
	challenger svcServer.Challenger
}

// New is a constructor for Server object, New always panicks if any error happens
func New(ctx context.Context, challenger svcServer.Challenger, address string) (*Server, error) {
	listener, err := new(net.ListenConfig).Listen(ctx, protocolTCP, address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	return &Server{
		listener:   listener,
		challenger: challenger,
	}, nil
}

// ListenAndAccept listens tcp and accepts/handles all incoming connections concurrently
func (s *Server) ListenAndAccept(ctx context.Context) error {
	for {
		log.FromContext(ctx).Info("accepting new connection ...")

		conn, err := s.listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}

		ctx = log.ContextWithAttributes(ctx, log.Attributes{"clientAddress": conn.RemoteAddr().String()})
		go s.handleConnection(ctx, conn)
	}
}

// handleConnection tries pow challenge on active connection,
// if challenge is solved then quote is sent to client
func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	log.FromContext(ctx).Info("handling new connection ...")

	if err := s.challenger.Challenge(ctx, transport.New(conn)); err != nil {
		log.FromContext(ctx).Error("failed to challenge client", log.StdError(err))
		return
	}
}

// Close closes the tcp listener so it stops acccepts new connection
func (s *Server) Close() error {
	return s.listener.Close()
}
