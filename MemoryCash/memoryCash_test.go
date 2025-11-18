package MemoryCash

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemo(t *testing.T) {
	t.Run("Double Get from cache", func(t *testing.T) {
		ctx := context.Background()
		count := 0

		memo := NewMemo(func(context.Context, string) (int, error) {
			count++
			return count, nil
		}, 10)
		get1, err := memo.Get(ctx, "first get")
		assert.NoError(t, err)
		get2, err := memo.Get(ctx, "first get")
		assert.NoError(t, err)
		assert.Equal(t, get1, get2)
		assert.Equal(t, 1, count, "Функция должна вызваться только 1 раз")
	})

	t.Run("Different keys", func(t *testing.T) {
		ctx := context.Background()
		count := 0
		memo := NewMemo(func(context.Context, string) (int, error) {
			count++
			return count, nil
		}, 10)
		get1, err := memo.Get(ctx, "first get")
		assert.NoError(t, err)
		get2, err := memo.Get(ctx, "second get")
		assert.NoError(t, err)

		assert.NotEqual(t, get1, get2, "Different keys should give different values")
		assert.Equal(t, count, 2, "Get should call 2 times")
	})

	t.Run("Context with timeout", func(t *testing.T) {
		memo := NewMemo(func(ctx context.Context, key string) (string, error) {
			select {
			case <-time.After(2 * time.Second):
				return "get done", nil
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}, 10)

		ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second)
		defer cancel1()

		ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel2()

		get1, err := memo.Get(ctx1, "timeout get")
		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Equal(t, get1, "")

		get2, err := memo.Get(ctx2, "success get")
		assert.NoError(t, err)
		assert.Equal(t, get2, "get done")
	})

	t.Run("Concurrency test", func(t *testing.T) {
		var callCounter int32
		memo := NewMemo(func(ctx context.Context, key string) (int, error) {
			atomic.AddInt32(&callCounter, 1)
			return 42, nil
		}, 4)

		ctx := context.Background()
		result := make(chan int, 10)

		for i := 0; i < 10; i++ {
			go func() {
				res, err := memo.Get(ctx, "test")
				assert.NoError(t, err)
				result <- res
			}()
		}

		for i := 0; i < 10; i++ {
			assert.Equal(t, <-result, 42)
		}

		assert.Equal(t, callCounter, int32(1))
	})

	t.Run("NonCached errors", func(t *testing.T) {
		var callCounter int32
		memo := NewMemo(func(ctx context.Context, key string) (int, error) {
			atomic.AddInt32(&callCounter, 1)
			return 0, errors.New("some error")
		}, 4)
		ctx := context.Background()
		_, err1 := memo.Get(ctx, "test")
		assert.Error(t, err1)
		_, err2 := memo.Get(ctx, "test")
		assert.Error(t, err2)

		assert.Equal(t, callCounter, int32(1))
	})
}
