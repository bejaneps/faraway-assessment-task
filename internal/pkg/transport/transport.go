package transport

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

const msgBuffer = 100

// Err0BytesWritten is used when 0 bytes are written to connection
var Err0BytesWritten = errors.New("0 bytes written")

// ReadWriter is convenience wrapper for net.Conn
type ReadWriter interface {
	// WriteMessage writes message to underlying connection
	WriteMessage(string) error
	// ReadMessage reads message from underlying connection
	ReadMessage() (string, error)
}

type transport struct {
	conn net.Conn
}

func New(conn net.Conn) ReadWriter {
	return &transport{
		conn: conn,
	}
}

func (t *transport) WriteMessage(msg string) error {
	n, err := t.conn.Write([]byte(msg))
	if err != nil && errors.Is(err, net.ErrClosed) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to write to server: %w", err)
	} else if n == 0 {
		return fmt.Errorf("failed to write to server: %w", Err0BytesWritten)
	}

	return nil
}

func (t *transport) ReadMessage() (string, error) {
	reader := bufio.NewReader(t.conn)

	msg := make([]byte, msgBuffer)

	for {
		n, err := reader.Read(msg)
		if err != nil && errors.Is(err, io.EOF) {
			return "", nil
		}
		if err != nil {
			return "", fmt.Errorf("failed to read message: %w", err)
		} else if n == 0 || len(msg) == 0 {
			continue
		}

		break
	}

	msg = bytes.ReplaceAll(msg, []byte("\x00"), []byte{})

	return strings.TrimSuffix(string(msg), "\r\n"), nil
}
