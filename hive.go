package buzz

type Hive struct {
	colony []*Worker
}

// New initializes a new [*Hive].
func New() *Hive {
	return &Hive{}
}

func (hive *Hive) AddWorker(task Task) *Worker {
	bee := &Worker{
		task: task,
	}
	hive.colony = append(hive.colony, bee)
	return bee
}
