package buzz

import (
	"slices"
	"sync"
)

type Hive struct {
	colony         []*Worker
	block          sync.WaitGroup
	notifyComplete chan *Worker
}

// New initializes a new [*Hive].
func New() *Hive {
	hive := &Hive{
		notifyComplete: make(chan *Worker, 1),
	}
	hive.startCleanupWorker()
	return hive
}

func (hive *Hive) startCleanupWorker() {
	go func() {
		for completed := range hive.notifyComplete {
			index := -1
			for i, w := range hive.colony {
				if completed == w {
					index = i
					break
				}
			}
			if index >= 0 {
				hive.colony = slices.Delete(hive.colony, index, index+1)
			}
		}
	}()
}

func (hive *Hive) Submit(worker *Worker) {
	hive.block.Add(1)
	worker.notifyComplete = hive.notifyComplete
	hive.colony = append(hive.colony, worker)
	go worker.run(&hive.block)
}

// StopAll should only be used when you are completely done with the hive. Internal channels are
// closed and all workers are shutdown.
func (hive *Hive) StopAll() {
	for i := range hive.colony {
		hive.colony[i].Stop()
	}
	hive.block.Wait()
	close(hive.notifyComplete)
}
