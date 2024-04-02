package buzz

import (
	"slices"
	"sync"
)

// Hive contains all the workers and synchronizes a graceful shutdown of all the workers.
type Hive struct {
	colony         []*Worker
	block          sync.WaitGroup
	notifyComplete chan struct{}
	middleware     []MiddleFunc
	closed         chan struct{}
}

// New initializes a new [*Hive].
func New() *Hive {
	hive := &Hive{
		notifyComplete: make(chan struct{}, 1),
		closed:         make(chan struct{}, 1),
	}
	hive.startCleanupWorker()
	return hive
}

func (hive *Hive) startCleanupWorker() {
	go func() {
		for range hive.notifyComplete {
			hive.removeDoneWorkers()
		}
		hive.removeDoneWorkers()
		hive.closed <- struct{}{}
	}()
}

func (hive *Hive) removeDoneWorkers() {
	finished := []int{}
	for i := range hive.colony {
		if hive.colony[i].done.Load() {
			finished = append(finished, i)
		}
	}
	for i := len(finished) - 1; i >= 0; i-- {
		hive.colony = slices.Delete(hive.colony, i, i+1)
	}
}

// Use adds the given MiddleFunc's to the hive as default functions. They will get added
// to each worker that is added to the hive. They are placed as the earliest middleware
// in the chain, in the same order they are added here. So, if you add A, B, C to the hive,
// and add a worker that already has D, and E middleware, you will end up with a middleware
// chain on that worker equivalent to A, B, C, D, E. From that point, it's important to note
// that any closures that are added as middleware to the hive may behave in unexpected ways
// since each worker will get the same closure (unless that is your intent!).
func (hive *Hive) Use(middleFunc ...MiddleFunc) {
	hive.middleware = append(hive.middleware, middleFunc...)
}

// Submit starts the workers running and adds them to the hive.
func (hive *Hive) Submit(workers ...*Worker) {
	for _, worker := range workers {
		hive.block.Add(1)
		worker.notifyComplete = hive.notifyComplete
		worker.middleware = append(hive.middleware, worker.middleware...)
		hive.colony = append(hive.colony, worker)
		go worker.run(&hive.block)
	}
}

// StopAll should only be used when you are completely done with the hive. Internal channels are
// closed and all workers are shutdown.
func (hive *Hive) StopAll() {
	for i := range hive.colony {
		hive.colony[i].Stop()
	}
	hive.block.Wait()
	close(hive.notifyComplete)
	<-hive.closed
}
