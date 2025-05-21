package scheduler

import (
	"context"
	"log"
	"sync"
	"time"
)

type Task struct {
	Name     string
	Interval time.Duration
	Function TaskFunc
}

type TaskFunc func(ctx context.Context) error

type Scheduler struct {
	tasks  []Task
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks: make([]Task, 0),
	}
}

func (s *Scheduler) AddTask(name string, interval time.Duration, task TaskFunc) {
	s.tasks = append(s.tasks, Task{
		Name:     name,
		Interval: interval,
		Function: task,
	})
}

func (s *Scheduler) Start(ctx context.Context) {
	ctx, s.cancel = context.WithCancel(ctx)

	for _, task := range s.tasks {
		s.wg.Add(1)
		go s.runTask(ctx, task)
	}
}

func (s *Scheduler) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
	log.Println("Scheduler stopped")
}

func (s *Scheduler) runTask(ctx context.Context, task Task) {
	defer s.wg.Done()

	ticker := time.NewTicker(task.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopping task '%s'...", task.Name)
			return
		case <-ticker.C:
			s.executeTask(ctx, task)
		}
	}
}

func (s *Scheduler) executeTask(ctx context.Context, task Task) {
	err := task.Function(ctx)
	if err != nil {
		log.Printf("Task '%s' failed: %v", task.Name, err)
		return
	}
}
