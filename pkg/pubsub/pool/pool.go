package pool

type Pool struct {
	sem  chan struct{}
	work chan func()
}

func New(size int) *Pool {
	p := &Pool{
		sem:  make(chan struct{}, size),
		work: make(chan func()),
	}

	return p
}

func (p *Pool) Schedule(task func()) {
	select {
	case p.work <- task:
		return
	case p.sem <- struct{}{}:
		go p.worker(task)
	}
}

func (p *Pool) worker(task func()) {
	defer func() { <-p.sem }()

	for {
		task()
		task = <-p.work
	}
}
