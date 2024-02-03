package job

import "fmt"

type Job interface {
	Execute()
}

type PrintJob struct {
	message string
}

func (p *PrintJob) Execute() {
	fmt.Println(p.message)
}

func NewPrintJob(message string) *PrintJob {
	return &PrintJob{
		message: message,
	}
}
