package workerpool

import (
	"context"
	"sync"
)

type Task[T any] struct {
	Val *T
}

type Result[T, R any] struct {
	Task *Task[T]
	Val  *R
	Err  error
}

type WorkerPool[T, R any] struct {
	wg         *sync.WaitGroup
	taskChan   chan *Task[T]
	resultChan chan *Result[T, R]
}

func NewWorkerPool[T, R any](ctx context.Context, workerCount int,
	hFun func(ctx context.Context, val *T) (*R, context.CancelFunc, error),
) *WorkerPool[T, R] {
	internalCtx, internalCtxCancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	taskChan := make(chan *Task[T])
	resultChan := make(chan *Result[T, R])

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range taskChan {
				if ctx.Err() != nil {
					resultChan <- &Result[T, R]{
						Task: task,
						Err:  ctx.Err(),
					}
					continue // мы должны дочитать все задачи из taskChan
				}
				if internalCtx.Err() != nil {
					resultChan <- &Result[T, R]{
						Task: task,
						Err:  internalCtx.Err(),
					}
					continue // мы должны дочитать все задачи из taskChan
				}

				val, err := hFun(ctx, task.Val)
				resultChan <- &Result[T, R]{
					Task: task,
					Val:  val,
					Err:  err,
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return &WorkerPool[T, R]{
		wg:         wg,
		taskChan:   taskChan,
		resultChan: resultChan,
	}
}

func (w *WorkerPool[T, R]) AddTask(ctx context.Context, val *T) {
	w.taskChan <- &Task[T]{
		Val: val,
	}
}

func (w *WorkerPool[T, R]) ResultChan() <-chan *Result[T, R] {
	return w.resultChan
}
