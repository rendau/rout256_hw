package workerpool

import (
	"context"
	"sync"
)

type TaskVal struct {
	Val int
}

type ResultVal struct {
	Val int
}

// --------------------------------------------

type Task struct {
	Ctx context.Context
	Val *TaskVal
}

type Result struct {
	Task *Task
	Val  *ResultVal
	Err  error
}

type WorkerPool struct {
	wg         *sync.WaitGroup
	taskChan   chan *Task
	resultChan chan *Result
}

func NewWorkerPool(globalCtx context.Context, workerCount int,
	hFun func(ctx context.Context, val *TaskVal) (*ResultVal, error),
) *WorkerPool {
	wg := &sync.WaitGroup{}
	taskChan := make(chan *Task)
	resultChan := make(chan *Result)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range taskChan {
				if globalCtx.Err() != nil {
					resultChan <- &Result{
						Task: task,
						Err:  globalCtx.Err(),
					}
					continue // мы должны дочитать все задачи из taskChan
				}
				if task.Ctx.Err() != nil {
					resultChan <- &Result{
						Task: task,
						Err:  task.Ctx.Err(),
					}
					continue // мы должны дочитать все задачи из taskChan
				}

				val, err := hFun(task.Ctx, task.Val)
				resultChan <- &Result{
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

	return &WorkerPool{
		wg:         wg,
		taskChan:   taskChan,
		resultChan: resultChan,
	}
}

func (w *WorkerPool) AddTask(ctx context.Context, val *TaskVal) {
	w.taskChan <- &Task{
		Ctx: ctx,
		Val: val,
	}
}

func (w *WorkerPool) ResultChan() <-chan *Result {
	return w.resultChan
}
