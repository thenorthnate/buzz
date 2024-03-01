package buzz

import "sync"

type Hive struct {
	colony []*Worker
	block  sync.WaitGroup
}

// New initializes a new [*Hive].
func New() *Hive {
	return &Hive{}
}

func (hive *Hive) Submit(worker *Worker) {
	hive.block.Add(1)
	hive.colony = append(hive.colony, worker)
	go worker.run(&hive.block)
}

func (hive *Hive) StopAll() {
	for i := range hive.colony {
		hive.colony[i].Stop()
	}
	hive.block.Wait()
}
