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

// delayedJob is the wrapper of Job which contains extra information required for the priority queue implementation
type delayedJob struct {
	job job.Job
	// priority is the epoch time at which this job needs to be executed
	priority int64
	// index is the position of this Job in the priorityQueue
	index int
}

// priorityJobQueue is the priority queue implementation which will store the delayedJob in their order of execution
type priorityJobQueue []*delayedJob

func (pq priorityJobQueue) Len() int { return len(pq) }

func (pq priorityJobQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq priorityJobQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityJobQueue) Push(x interface{}) {
	n := len(*pq)
	delayedJob := x.(*delayedJob)
	delayedJob.index = n
	*pq = append(*pq, delayedJob)
}

func (pq *priorityJobQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	delayedJob := old[n-1]
	old[n-1] = nil        // avoid memory leak
	delayedJob.index = -1 // for safety
	*pq = old[0 : n-1]
	return delayedJob
}

type delayedScheduler struct {
	ticker           *time.Ticker
	priorityJobQueue priorityJobQueue
	sync.Mutex
	done chan bool
}

func NewScheduler() *delayedScheduler {
	delayedScheduler := &delayedScheduler{
		priorityJobQueue: make(priorityJobQueue, 0),
		done:             make(chan bool),
		ticker:           time.NewTicker(1 * time.Second),
	}
	heap.Init(&delayedScheduler.priorityJobQueue)
	go delayedScheduler.start()
	return delayedScheduler
}

func (ds *delayedScheduler) start() {
	go func() {
		for {
			select {
			case <-ds.done:
				ds.ticker.Stop()
				return
			case t := <-ds.ticker.C:
				fmt.Println("Tick at ", t)
				ds.Lock()
				if len(ds.priorityJobQueue) != 0 {
					job := heap.Pop(&ds.priorityJobQueue).(*delayedJob)
					now := time.Now().Unix()
					if job.priority-now > 0 {
						heap.Push(&ds.priorityJobQueue, job)
					} else {
						go job.job.Execute()
					}
				}
				ds.Unlock()
			}
		}
	}()
}

func (ds *delayedScheduler) Stop() {
	ds.done <- true
}

func (ds *delayedScheduler) Schedule(job job.Job, duration time.Duration) {
	ds.Lock()
	now := time.Now().Unix()
	timeOfExecution := now + int64(duration.Seconds())
	delayedJob := &delayedJob{
		job:      job,
		priority: timeOfExecution,
	}
	heap.Push(&ds.priorityJobQueue, delayedJob)
	ds.Unlock()
}
