package benchmarking

import (
	"github.com/NoobMaster-96/delayed-scheduler-golang/delayed_scheduler"
	"github.com/NoobMaster-96/delayed-scheduler-golang/job"
	"runtime/debug"
	"testing"
	"time"
)

func BenchmarkScheduler(b *testing.B) {
	debug.SetGCPercent(500)
	delayedScheduler := delayed_scheduler.NewScheduler()

	sumJob1 := job.NewSumJob(2, 3)
	sumJob2 := job.NewSumJob(3, 4)
	sumJob3 := job.NewSumJob(4, 5)

	for i := 0; i < b.N; i++ {
		delayedScheduler.Schedule(sumJob1, 5*time.Second)
		delayedScheduler.Schedule(sumJob2, 1*time.Second)
		delayedScheduler.Schedule(sumJob3, 3*time.Second)
	}
}
