package job

type Job interface {
	Execute()
}

type SumJob struct {
	a   int64
	b   int64
	sum int64
}

func (p *SumJob) Execute() {
	p.sum = p.a + p.b
}

func NewSumJob(a, b int64) *SumJob {
	return &SumJob{
		a: a,
		b: b,
	}
}
