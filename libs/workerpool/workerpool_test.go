package workerpool

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWorkerPool(t *testing.T) {
	workerCount := 2
	taskCount := 5

	// Define a worker function that sleeps for 100ms and returns the input value
	workerFun := func(ctx context.Context, val *int) (*int, error) {
		time.Sleep(100 * time.Millisecond)
		res := *val + 1
		return &res, nil
	}

	// Define a task adder function that adds 5 tasks to the worker pool
	taskAdder := func(add func(val *int) bool) {
		for i := 0; i < taskCount; i++ {
			val := i
			add(&val)
		}
	}

	// Define a result function that checks if the result value is equal to the input value
	resultFun := func(task *int, result *int, err error) error {
		if err != nil {
			return err
		}
		if *task+1 != *result {
			return errors.New("result value does not match input value")
		}
		return nil
	}

	wp := NewWorkerPool(context.Background(), workerCount, workerFun, taskAdder, resultFun)

	err := wp.Wait()
	require.NoError(t, err, "worker pool failed with error")
	require.True(t, wp.GetDuration() > 290*time.Millisecond, "worker pool finished too fast")
	require.True(t, wp.GetDuration() < 310*time.Millisecond, "worker pool finished too slow")

	wp = NewWorkerPool(context.Background(), 1, workerFun, taskAdder, resultFun)

	err = wp.Wait()
	require.NoError(t, err, "worker pool failed with error")
	require.True(t, wp.GetDuration() > 490*time.Millisecond, "worker pool finished too fast")
	require.True(t, wp.GetDuration() < 510*time.Millisecond, "worker pool finished too slow")

	wp = NewWorkerPool(context.Background(), workerCount, workerFun, taskAdder, resultFun)

	time.Sleep(150 * time.Millisecond)
	wp.Cancel()

	err = wp.Wait()
	require.NoError(t, err, "worker pool failed with error")
	require.True(t, wp.GetDuration() > 150*time.Millisecond, "worker pool finished too fast")
	require.True(t, wp.GetDuration() < 160*time.Millisecond, "worker pool finished too slow")
}
