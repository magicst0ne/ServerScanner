package dispatcher

import (
	"sync"
)

type callFunc func(v interface{})

type Dispatcher struct {
	pool    chan *worker
	wg      sync.WaitGroup
	workers []*worker
	queue   chan interface{}
	quit    chan struct{}
	fn      callFunc
	flag    int
}

// Start starts the specified dispatcher but does not wait for it to complete.
func (d *Dispatcher) Start() {
	for _, w := range d.workers {
		w.start()
	}

	go func() {
		for {
			select {
			case v := <-d.queue:
				(<-d.pool).data <- v

			case <-d.quit:
				return
			}
		}
	}()
}

// Add adds a given value to the queue of the dispatcher.
func (d *Dispatcher) Add(v interface{}) {
	d.wg.Add(1)
	d.queue <- v
}

func (d *Dispatcher) Wait() {
	d.wg.Wait()
}

func (d *Dispatcher) SetFlag(flag int) {
	d.flag = flag
}

func (d *Dispatcher) GetFlag(flag int) int {
	return d.flag
}

// NewDispatcher returns a pointer of Dispatcher.
func NewDispatcher(maxWorkers int, maxQueues int, fn callFunc) *Dispatcher {
	d := &Dispatcher{
		pool:  make(chan *worker, maxWorkers),
		queue: make(chan interface{}, maxQueues),
		quit:  make(chan struct{}),
		fn:    fn,
		flag:  0,
	}

	// worker init
	d.workers = make([]*worker, cap(d.pool))
	for i := 0; i < cap(d.pool); i++ {
		w := worker{
			dispatcher: d,
			data:       make(chan interface{}),
			quit:       make(chan struct{}),
		}
		d.workers[i] = &w
	}
	return d
}

type worker struct {
	dispatcher *Dispatcher
	data       chan interface{}
	quit       chan struct{}
}

func (w *worker) start() {
	go func() {
		for {
			w.dispatcher.pool <- w

			select {
			case v := <-w.data:
				w.dispatcher.fn(v)
				w.dispatcher.wg.Done()

			case <-w.quit:
				return
			}
		}
	}()
}
