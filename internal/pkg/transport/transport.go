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

var (
	// Err0BytesRead is used when 0 bytes are read from connection
	Err0BytesRead = errors.New("0 bytes read")
	// Err0BytesWritten is used when 0 bytes are written to connection
	Err0BytesWritten = errors.New("0 bytes written")
)

func WriteMessage(conn net.Conn, msg string) error {
	n, err := conn.Write([]byte(msg))
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

func ReadMessage(reader *bufio.Reader) (string, error) {
	msg := make([]byte, msgBuffer)
	for {
		n, err := reader.Read(msg)
		if err != nil && errors.Is(err, io.EOF) {
			return "", nil
		}
		if err != nil {
			return "", fmt.Errorf("failed to read message: %w", err)
		} else if n == 0 {
			continue
		}

		break
	}

	msg = bytes.ReplaceAll(msg, []byte("\x00"), []byte{})

	return strings.TrimSuffix(string(msg), "\r\n"), nil
}
