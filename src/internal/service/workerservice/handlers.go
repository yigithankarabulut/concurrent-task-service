package workerservice

import (
	"fmt"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/models"
	"github.com/yigithankarabulut/ConcurrentTaskService/src/internal/service/taskservice/dto"
)

func (t *taskWorker) Submit(f models.TaskJobModel) (any, error) {
	t.ReqChan <- f
	for {
		select {
		case <-f.Context.Done():
			return nil, f.Context.Err()
		case err := <-t.ErrChan:
			return nil, err
		case res := <-t.ResChan:
			return res, nil
		case <-t.done:
			t.logger.Info("worker closed while processing the request", "job: ", f.JOB)
			return nil, fmt.Errorf("worker closed while processing the request")
		}
	}
}

func (w *taskWorker) get(f models.TaskJobModel) {
	req := dto.GetTaskRequest{
		ID: f.ID,
	}
	resp, err := w.service.Get(f.Context, req)
	if err != nil {
		select {
		case <-w.done:
			close(w.ResChan)
			close(w.ErrChan)
			close(w.ReqChan)
			return
		default:
			w.ErrChan <- err
			return
		}
	}
	select {
	case <-w.done:
		close(w.ResChan)
		close(w.ErrChan)
		close(w.ReqChan)
		return
	default:
		w.ResChan <- resp
	}
}

func (w *taskWorker) set(f models.TaskJobModel) {
	req := dto.SetTaskRequest{
		ID:          f.ID,
		Title:       f.Title,
		Description: f.Description,
		Status:      f.Status,
	}
	resp, err := w.service.Set(f.Context, req)
	if err != nil {
		select {
		case <-w.done:
			close(w.ResChan)
			close(w.ErrChan)
			close(w.ReqChan)
			return
		default:
			w.ErrChan <- err
			return
		}
	}
	select {
	case <-w.done:
		close(w.ResChan)
		close(w.ErrChan)
		close(w.ReqChan)
		return
	default:
		w.ResChan <- resp
	}
}

func (w *taskWorker) delete(f models.TaskJobModel) {
	req := dto.DeleteTaskRequest{
		ID: f.ID,
	}
	err := w.service.Delete(f.Context, req)
	if err != nil {
		select {
		case <-w.done:
			close(w.ResChan)
			close(w.ErrChan)
			close(w.ReqChan)
			return
		default:
			w.ErrChan <- err
			return
		}
	}
	select {
	case <-w.done:
		close(w.ResChan)
		close(w.ErrChan)
		close(w.ReqChan)
		return
	default:
		w.ResChan <- nil
	}
}

func (w *taskWorker) update(f models.TaskJobModel) {
	req := dto.UpdateTaskRequest{
		ID:          f.ID,
		Title:       f.Title,
		Description: f.Description,
		Status:      f.Status,
	}
	resp, err := w.service.Update(f.Context, req)
	if err != nil {
		select {
		case <-w.done:
			close(w.ResChan)
			close(w.ErrChan)
			close(w.ReqChan)
			return
		default:
			w.ErrChan <- err
			return
		}
	}
	select {
	case <-w.done:
		close(w.ResChan)
		close(w.ErrChan)
		close(w.ReqChan)
		return
	default:
		w.ResChan <- resp
	}
}

func (w *taskWorker) list(f models.TaskJobModel) {
	req := dto.ListTaskRequest{
		Status: f.Status,
	}
	resp, err := w.service.List(f.Context, req)
	if err != nil {
		select {
		case <-w.done:
			close(w.ResChan)
			close(w.ErrChan)
			close(w.ReqChan)
			return
		default:
			w.ErrChan <- err
			return
		}
	}
	select {
	case <-w.done:
		close(w.ResChan)
		close(w.ErrChan)
		close(w.ReqChan)
		return
	default:
		w.ResChan <- resp
	}
}

func (w *taskWorker) worker() {
	defer w.Wg.Done()

	for {
		select {
		case <-w.done:
			w.mu.Lock()
			w.workerCount--
			if w.workerCount == 0 {
				w.logger.Info("all workers closed.")
			}
			w.mu.Unlock()
			return
		case f := <-w.ReqChan:
			switch f.JOB {
			case "GET":
				w.get(f)
			case "SET":
				w.set(f)
			case "DELETE":
				w.delete(f)
			case "UPDATE":
				w.update(f)
			case "LIST":
				w.list(f)
			}
		}
	}
}
