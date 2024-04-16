package kafka

type pool struct {
	sem  chan struct{}
	work chan func()
}

func newPool(size int) *pool {
	p := &pool{
		sem:  make(chan struct{}, size),
		work: make(chan func()),
	}

	return p
}

func (p *pool) Schedule(task func()) {
	select {
	case p.work <- task:
		return
	case p.sem <- struct{}{}:
		go p.worker(task)
	}
}

func (p *pool) worker(task func()) {
	defer func() { <-p.sem }()

	for {
		task()
		task = <-p.work
	}
}
