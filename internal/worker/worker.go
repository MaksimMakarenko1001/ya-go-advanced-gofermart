package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type Job func(ctx context.Context) (err error)

type Worker struct {
	config  Config
	locker  Locker
	pid     string
	key     string
	segment string
	job     Job
}

func New(config Config, locker Locker, pid, key, segment string, job Job) *Worker {
	return &Worker{
		config:  config,
		locker:  locker,
		pid:     pid,
		key:     key,
		segment: segment,
		job:     job,
	}
}

func (w *Worker) Start(ctx context.Context) {
	go w.run(ctx)
}

func (w *Worker) run(ctx context.Context) {
	onceCh := make(chan struct{}, 1)
	defer close(onceCh)

	for {
		select {
		case <-ctx.Done():
			return
		case onceCh <- struct{}{}:
			jobCtx, cancel := context.WithTimeout(ctx, w.config.JobTimeout)
			go func() {
				ts := time.Now()
				defer cancel()

				if err := w.doJob(jobCtx); err != nil {
					log.Println(err)
				}
				time.Sleep(w.config.JobInterval - time.Since(ts))
				<-onceCh
			}()
		}
	}
}

func (w *Worker) doJob(ctx context.Context) (err error) {
	ts := time.Now()
	ok, err := w.locker.LockAcquire(ctx, w.key, w.segment, ts.Add(w.config.JobInterval), w.pid)
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !ok {
		log.Printf("lock acquisition is not ok:{key=%s,segment=%s,pid=%s}\n", w.key, w.segment, w.pid)
		return nil
	}

	defer func() {
		ok, releaseErr := w.locker.LockRelease(ctx, w.key, w.segment, w.pid)
		if releaseErr != nil {
			err = errors.Join(
				err,
				fmt.Errorf("failed to release lock: %w", releaseErr),
			)
		}
		if !ok {
			log.Printf("lock releasing is not ok:{key=%s,segment=%s,pid=%s}\n", w.key, w.segment, w.pid)
		}
	}()

	err = w.job(ctx)
	if err != nil {
		return fmt.Errorf("failed to do job: %w", err)
	}

	return nil
}
