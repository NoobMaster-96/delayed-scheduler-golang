package main

import (
	"github.com/NoobMaster-96/delayed-scheduler-golang/delayed_scheduler"
	"github.com/NoobMaster-96/delayed-scheduler-golang/job"
	"time"
)

func main() {
	delayedScheduler := delayed_scheduler.NewDelayedScheduler()

	printJob1 := job.NewPrintJob("Hey!!")
	printJob2 := job.NewPrintJob("There!!")
	printJob3 := job.NewPrintJob("User!!")

	delayedScheduler.Schedule(printJob1, 5*time.Second)
	time.Sleep(3 * time.Second)
	delayedScheduler.Schedule(printJob2, 1*time.Second)
	delayedScheduler.Schedule(printJob3, 3*time.Second)

	time.Sleep(10 * time.Second)
	delayedScheduler.Stop()
}
