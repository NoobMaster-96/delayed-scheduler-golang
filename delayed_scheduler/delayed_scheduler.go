package delayed_scheduler

import (
	"github.com/NoobMaster-96/delayed-scheduler-golang/job"
	"time"
)

type Scheduler interface {
	// Schedule takes a job and executes it after the given duration of time in sec
	Schedule(job job.Job, duration time.Duration)
}

type delayedScheduler struct {
}

func NewScheduler() *delayedScheduler {
	return &delayedScheduler{}
}

func (d *delayedScheduler) Schedule(job job.Job, duration time.Duration) {
	go func() {
		time.Sleep(duration)
		job.Execute()
	}()
}
