package delayed_scheduler

import (
	"container/heap"
	"fmt"
	"github.com/NoobMaster-96/delayed-scheduler-golang/job"
	"sync"
	"time"
)

type Scheduler interface {
	// Schedule takes a job and executes it after the given duration of time in sec
	Schedule(job job.Job, duration time.Duration)
}

type DelayedJob struct {
	job      job.Job
	priority int64
	index    int
}

type PriorityJobQueue []*DelayedJob

func (pq PriorityJobQueue) Len() int { return len(pq) }

func (pq PriorityJobQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityJobQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityJobQueue) Push(x interface{}) {
	n := len(*pq)
	delayedJob := x.(*DelayedJob)
	delayedJob.index = n
	*pq = append(*pq, delayedJob)
}

func (pq *PriorityJobQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	delayedJob := old[n-1]
	old[n-1] = nil        // avoid memory leak
	delayedJob.index = -1 // for safety
	*pq = old[0 : n-1]
	return delayedJob
}

type DelayedScheduler struct {
	priorityJobQueue PriorityJobQueue
	sync.Mutex
	done chan bool
}

func NewDelayedScheduler() *DelayedScheduler {
	delayedScheduler := &DelayedScheduler{
		priorityJobQueue: make(PriorityJobQueue, 0),
		done:             make(chan bool),
	}
	heap.Init(&delayedScheduler.priorityJobQueue)
	go delayedScheduler.Start()
	return delayedScheduler
}

func (ds *DelayedScheduler) Start() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ds.done:
				ticker.Stop()
				return
			case t := <-ticker.C:
				fmt.Println("Tick at ", t)
				ds.Lock()
				if len(ds.priorityJobQueue) != 0 {
					delayedJob := heap.Pop(&ds.priorityJobQueue).(*DelayedJob)
					now := time.Now().Unix()
					if delayedJob.priority-now > 0 {
						heap.Push(&ds.priorityJobQueue, delayedJob)
					} else {
						delayedJob.job.Execute()
					}
				}
				ds.Unlock()
			}
		}
	}()
}

func (ds *DelayedScheduler) Stop() {
	ds.done <- true
}

func (ds *DelayedScheduler) Schedule(job job.Job, duration time.Duration) {
	ds.Lock()
	now := time.Now().Unix()
	runTime := now + int64(duration.Seconds())
	delayedJob := &DelayedJob{
		job:      job,
		priority: runTime,
	}
	heap.Push(&ds.priorityJobQueue, delayedJob)
	ds.Unlock()
}
