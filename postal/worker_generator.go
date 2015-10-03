package postal

type WorkerGenerator struct {
	InstanceIndex int
	Count         int
}

type Worker interface {
	Work()
}

func (w WorkerGenerator) Work(workerFunc func(id int) Worker) {
	firstID := w.InstanceIndex*w.Count + 1
	for i := 0; i < w.Count; i++ {
		workerFunc(firstID + i).Work()
	}
}
