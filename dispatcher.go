package Gworker

// Dispatcher is a work pool handler
type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool   chan chan Job
	maxWorkers   int
	maxQueue     int
	maxPrioQueue int
	JobQueue     chan Job // JobQueue is a buffered channel that we can send work requests on
	PriJobQueue  chan Job // PrioJobQueue is a buffered channel that we can send work requests on
}

// NewDispatcher Create a new dispatcher
func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		WorkerPool: pool,
	}
}

// Run Worker pool handler
func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	sendToWorker := func(job Job) {
		// a job request has been received
		// try to obtain a worker job channel that is available.
		// this will block until a worker is idle
		jobChannel := <-d.WorkerPool

		// dispatch the job to the worker job channel
		jobChannel <- job
	}

	for {
		select {
		case job := <-d.PriJobQueue:
			sendToWorker(job)
		case job := <-d.JobQueue:
			sendToWorker(job)
		}
	}
}
