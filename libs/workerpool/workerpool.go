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
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
	taskChan   chan *Task[T]
	resultChan chan *Result[T, R]
	finishChan chan error
}

func NewWorkerPool[T, R any](
	ctx context.Context,
	workerCount int,
	serveFun func(ctx context.Context, val *T) (*R, error),
	taskAdder func(add func(val *T) bool),
	resultFun func(task *T, result *R, err error) error,
) *WorkerPool[T, R] {
	ctx, cancel := context.WithCancel(ctx) // создаем новый контекст, чтобы не отменить переданный

	wp := &WorkerPool[T, R]{
		ctx:        ctx,
		cancel:     cancel,
		wg:         &sync.WaitGroup{},
		taskChan:   make(chan *Task[T], 1),
		resultChan: make(chan *Result[T, R], 1),
		finishChan: make(chan error, 1),
	}

	// запускаем воркеры
	for i := 0; i < workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(serveFun)
	}

	// закрываем resultChan после завершения всех воркеров
	go func() {
		wp.wg.Wait()
		close(wp.resultChan)
	}()

	// добавляем задачи
	go func() {
		defer close(wp.taskChan)
		taskAdder(wp.taskAddFun)
	}()

	// обрабатываем результаты
	go wp.resultHandler(resultFun)

	return wp
}

func (w *WorkerPool[T, R]) worker(serveFun func(ctx context.Context, val *T) (*R, error)) {
	defer w.wg.Done()

	var task *Task[T]
	var closed bool

	for {
		select {
		case <-w.ctx.Done(): // сперва проверяем, что контекст не отменен
			return
		default:
			select {
			case <-w.ctx.Done():
				return
			case task, closed = <-w.taskChan:
				if closed {
					return
				}

				// выполняем задачу
				val, err := serveFun(w.ctx, task.Val)

				// здесь уже канал может никто не слушать, поэтому используем select
				select {
				case <-w.ctx.Done():
					return
				case w.resultChan <- &Result[T, R]{Task: task, Val: val, Err: err}:
				}
			}
		}
	}
}

func (w *WorkerPool[T, R]) taskAddFun(val *T) bool {
	select {
	case <-w.ctx.Done(): // сперва проверяем, что контекст не отменен
		return false
	default:
		select {
		case <-w.ctx.Done():
			return false
		case w.taskChan <- &Task[T]{Val: val}:
			return true
		}
	}
}

func (w *WorkerPool[T, R]) resultHandler(resultFun func(task *T, result *R, err error) error) {
	var err error

	defer func() {
		close(w.finishChan)
	}()

	for {
		select {
		case <-w.ctx.Done(): // сперва проверяем, что контекст не отменен
			return
		default:
			select {
			case <-w.ctx.Done():
				return
			case res, closed := <-w.resultChan:
				if closed {
					return
				}
				err = resultFun(res.Task.Val, res.Val, res.Err)
				if err != nil {
					w.cancel()
					w.finishChan <- err
					return
				}
			}
		}
	}
}

func (w *WorkerPool[T, R]) Cancel() {
	w.cancel()
}

func (w *WorkerPool[T, R]) Wait() error {
	return <-w.finishChan
}
