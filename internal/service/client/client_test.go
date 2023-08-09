package client

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"testing"

	transportMocks "github.com/bejaneps/faraway-assessment-task/internal/pkg/transport/mocks"
	"github.com/stretchr/testify/require"
)

func TestRespond(t *testing.T) {
	ctx := context.TODO()

	var (
		number      = int64(10)
		hash        = fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d", number))))
		quote       = "random quote"
		randomError = errors.New("random error")
	)

	t.Run("success", func(t *testing.T) {
		rw := new(transportMocks.ReadWriter)

		rw.
			On("ReadMessage").
			Return(hash, nil)
		rw.
			On("WriteMessage", fmt.Sprintf("%d", number)).
			Return(nil)
		rw.
			On("ReadMessage").
			Return(quote, nil)

		service := New(20)
		err := service.Respond(ctx, rw)
		require.NoError(t, err)
	})

	t.Run("fail read hash error", func(t *testing.T) {
		rw := new(transportMocks.ReadWriter)

		rw.
			On("ReadMessage").
			Return("", randomError)

		service := New(20)
		err := service.Respond(ctx, rw)
		require.Error(t, err)
		require.Equal(t, err.Error(), "failed to read hash: random error")
	})

	t.Run("fail number not found error", func(t *testing.T) {
		rw := new(transportMocks.ReadWriter)

		rw.
			On("ReadMessage").
			Return(hash, nil)

		service := New(5)
		err := service.Respond(ctx, rw)
		require.Error(t, err)
		require.Equal(t, err, ErrNumberNotFound)
	})

	t.Run("fail send guess number error", func(t *testing.T) {
		rw := new(transportMocks.ReadWriter)

		rw.
			On("ReadMessage").
			Return(hash, nil)
		rw.
			On("WriteMessage", fmt.Sprintf("%d", number)).
			Return(randomError)

		service := New(20)
		err := service.Respond(ctx, rw)
		require.Error(t, err)
		require.Equal(t, err.Error(), "failed to send guess number: random error")
	})

	t.Run("fail read quote message error", func(t *testing.T) {
		rw := new(transportMocks.ReadWriter)

		rw.
			On("ReadMessage").
			Return(hash, nil).
			Once()
		rw.
			On("WriteMessage", fmt.Sprintf("%d", number)).
			Return(nil)
		rw.
			On("ReadMessage").
			Return("", randomError).
			Once()

		service := New(20)
		err := service.Respond(ctx, rw)
		require.Error(t, err)
		require.Equal(t, err.Error(), "failed to read quote message: random error")
	})

	t.Run("fail incorrect guess error", func(t *testing.T) {
		rw := new(transportMocks.ReadWriter)

		rw.
			On("ReadMessage").
			Return(hash, nil).
			Once()
		rw.
			On("WriteMessage", fmt.Sprintf("%d", number)).
			Return(nil)
		rw.
			On("ReadMessage").
			Return(incorectGuess, nil).
			Once()

		service := New(20)
		err := service.Respond(ctx, rw)
		require.Error(t, err)
		require.Equal(t, err.Error(), "failed to guess number, msg from server: incorrect guess")
	})
}
