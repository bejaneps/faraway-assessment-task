package server

import (
	"context"
	"fmt"
	"net"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/log"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/quote"
	"github.com/bejaneps/faraway-assessment-task/internal/pkg/transport"
)

const protocolTCP = "tcp"

type Challenger interface {
	// Challenge challenges client with a number range algorithm to solve
	Challenge(conn net.Conn) error
}

// Server is a simple struct to hold dependencies of server
type Server struct {
	listener   net.Listener
	challenger Challenger
}

// New is a constructor for Server object, New always panicks if any error happens
func New(ctx context.Context, challenger Challenger, address string) *Server {
	listener, err := new(net.ListenConfig).Listen(ctx, protocolTCP, address)
	if err != nil {
		panic(fmt.Errorf("failed to listen: %w", err))
	}

	return &Server{
		listener:   listener,
		challenger: challenger,
	}
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

	if err := s.challenger.Challenge(conn); err != nil {
		log.FromContext(ctx).Error("failed to challenge client", log.StdError(err))
		return
	}

	q := quote.Random()
	log.FromContext(ctx).Info("challenge was solved succesfully, sending quote", log.String("quote", q))
	if err := transport.WriteMessage(conn, q); err != nil {
		log.FromContext(ctx).Error("failed to send quote", log.StdError(err))
		return
	}
}

// Close closes the tcp listener so it stops acccepts new connection
func (s *Server) Close() error {
	return s.listener.Close()
}
