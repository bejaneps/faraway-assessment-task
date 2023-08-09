package server

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"testing"

	transportMocks "github.com/bejaneps/faraway-assessment-task/internal/pkg/transport/mocks"
	repoServerMocks "github.com/bejaneps/faraway-assessment-task/internal/repository/server/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestChallenge(t *testing.T) {
	ctx := context.TODO()

	var (
		number      = int64(10)
		hash        = fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d", number))))
		quote       = "random quote"
		randomError = errors.New("random error")

		cleanupFunc = func() {
			randomNumberGenerator = getRandomNumberAndHash
		}
	)

	t.Run("success", func(t *testing.T) {
		t.Cleanup(cleanupFunc)

		randomNumberGenerator = func(max *big.Int) (*big.Int, string, error) {
			return big.NewInt(number), hash, nil
		}

		rw := new(transportMocks.ReadWriter)
		rw.
			On("WriteMessage", hash).
			Return(nil)
		rw.
			On("ReadMessage").
			Return(fmt.Sprintf("%d", number), nil)
		rw.
			On("WriteMessage", quote).
			Return(nil)

		quoter := new(repoServerMocks.Quoter)
		quoter.
			On("Quote", mock.Anything).
			Return(quote, nil)

		service := New(quoter)
		err := service.Challenge(ctx, rw)
		require.NoError(t, err)
	})

	t.Run("fail random number generator fails", func(t *testing.T) {
		t.Cleanup(cleanupFunc)

		randomNumberGenerator = func(max *big.Int) (*big.Int, string, error) {
			return nil, "", randomError
		}

		service := New(nil)
		err := service.Challenge(ctx, nil)
		require.Error(t, err)
		require.Equal(t, err.Error(), "failed to generate random number and calculate it's hash: random error")
	})

	t.Run("fail random number send error", func(t *testing.T) {
		t.Cleanup(cleanupFunc)

		randomNumberGenerator = func(max *big.Int) (*big.Int, string, error) {
			return big.NewInt(number), hash, nil
		}

		rw := new(transportMocks.ReadWriter)
		rw.
			On("WriteMessage", hash).
			Return(randomError)

		service := New(nil)
		err := service.Challenge(ctx, rw)
		require.Error(t, err)
		require.Equal(t, err.Error(), "failed to send random number: random error")
	})

	t.Run("fail guess number read error", func(t *testing.T) {
		t.Cleanup(cleanupFunc)

		randomNumberGenerator = func(max *big.Int) (*big.Int, string, error) {
			return big.NewInt(number), hash, nil
		}

		rw := new(transportMocks.ReadWriter)
		rw.
			On("WriteMessage", hash).
			Return(nil)
		rw.
			On("ReadMessage").
			Return("", randomError)

		service := New(nil)
		err := service.Challenge(ctx, rw)
		require.Error(t, err)
		require.Equal(t, err.Error(), "failed to read guess number: random error")
	})

	t.Run("fail guess number invalid error", func(t *testing.T) {
		t.Cleanup(cleanupFunc)

		randomNumberGenerator = func(max *big.Int) (*big.Int, string, error) {
			return big.NewInt(number), hash, nil
		}

		rw := new(transportMocks.ReadWriter)
		rw.
			On("WriteMessage", hash).
			Return(nil)
		rw.
			On("ReadMessage").
			Return("invalidNumber", nil)

		service := New(nil)
		err := service.Challenge(ctx, rw)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to convert to int guess message:")
	})

	t.Run("fail incorrect guess message send error", func(t *testing.T) {
		t.Cleanup(cleanupFunc)

		randomNumberGenerator = func(max *big.Int) (*big.Int, string, error) {
			return big.NewInt(number), hash, nil
		}

		rw := new(transportMocks.ReadWriter)
		rw.
			On("WriteMessage", hash).
			Return(nil)
		rw.
			On("ReadMessage").
			Return("9", nil)
		rw.
			On("WriteMessage", incorectGuess).
			Return(randomError)

		service := New(nil)
		err := service.Challenge(ctx, rw)
		require.Error(t, err)
		require.Equal(t, err.Error(), "failed to send incorrect guess message: random error")
	})

	t.Run("fail guess number wrong error", func(t *testing.T) {
		t.Cleanup(cleanupFunc)

		randomNumberGenerator = func(max *big.Int) (*big.Int, string, error) {
			return big.NewInt(number), hash, nil
		}

		rw := new(transportMocks.ReadWriter)
		rw.
			On("WriteMessage", hash).
			Return(nil)
		rw.
			On("ReadMessage").
			Return("9", nil)
		rw.
			On("WriteMessage", incorectGuess).
			Return(nil)

		service := New(nil)
		err := service.Challenge(ctx, rw)
		require.Error(t, err)
		require.Equal(t, err, ErrNumberNotFound)
	})

	t.Run("fail get random quote error", func(t *testing.T) {
		t.Cleanup(cleanupFunc)

		randomNumberGenerator = func(max *big.Int) (*big.Int, string, error) {
			return big.NewInt(number), hash, nil
		}

		rw := new(transportMocks.ReadWriter)
		rw.
			On("WriteMessage", hash).
			Return(nil)
		rw.
			On("ReadMessage").
			Return(fmt.Sprintf("%d", number), nil)

		quoter := new(repoServerMocks.Quoter)
		quoter.
			On("Quote", mock.Anything).
			Return("", randomError)

		service := New(quoter)
		err := service.Challenge(ctx, rw)
		require.Error(t, err)
		require.Equal(t, err.Error(), "failed to get random quote: random error")
	})

	t.Run("fail get random quote error", func(t *testing.T) {
		t.Cleanup(cleanupFunc)

		randomNumberGenerator = func(max *big.Int) (*big.Int, string, error) {
			return big.NewInt(number), hash, nil
		}

		rw := new(transportMocks.ReadWriter)
		rw.
			On("WriteMessage", hash).
			Return(nil)
		rw.
			On("ReadMessage").
			Return(fmt.Sprintf("%d", number), nil)
		rw.
			On("WriteMessage", quote).
			Return(randomError)

		quoter := new(repoServerMocks.Quoter)
		quoter.
			On("Quote", mock.Anything).
			Return(quote, nil)

		service := New(quoter)
		err := service.Challenge(ctx, rw)
		require.Error(t, err)
		require.Equal(t, err.Error(), "failed to send quote: random error")
	})
}
