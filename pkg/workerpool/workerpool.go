package workerpool

type WorkerPool[T any] struct {
	num      int
	Work     chan T
	f        func(T)
	Vacation chan bool
}

func NewWorkerPool[T any](num int, f func(args T)) *WorkerPool[T] {
	wp := &WorkerPool[T]{
		num:      num,
		Work:     make(chan T, num),
		Vacation: make(chan bool),
	}

	for i := 0; i < num; i++ {
		go func() {
			for {
				select {
				case args := <-wp.Work:
					f(args)
				case <-wp.Vacation:
					return
				}
			}
		}()
	}

	return wp
}

func (wp *WorkerPool[T]) Submit(args T) {
	wp.Work <- args
}

func (wp *WorkerPool[T]) Concurrency() int {
	return wp.num
}

func (wp *WorkerPool[T]) Tune(n int) {
	for n > 0 {
		go func() {
			for {
				select {
				case args := <-wp.Work:
					wp.f(args)
				case <-wp.Vacation:
					return
				}
			}
		}()
		n--
	}
	wp.num += n
}

func (wp *WorkerPool[T]) TuneTo(n int) {
	i := n - wp.num
	for i > 0 {
		go func() {
			for {
				select {
				case args := <-wp.Work:
					wp.f(args)
				case <-wp.Vacation:
					return
				}
			}
		}()
		i--
	}
	wp.num = n
}

func (wp *WorkerPool[T]) Fire(n int) {
	for i := 0; i < n; i++ {
		wp.Vacation <- true
	}
	wp.num -= n
}

func (wp *WorkerPool[T]) FireTo(n int) {
	i := wp.num - n
	for i > 0 {
		wp.Vacation <- true
		i--
	}
	wp.num = n
}

func (wp *WorkerPool[T]) Release() {
	wp.Vacation <- true
}

func (wp *WorkerPool[T]) Len() int {
	return len(wp.Work)
}

func (wp *WorkerPool[T]) Cap() int {
	return cap(wp.Work)
}
