package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	gworker "github.com/pecigonzalo/Gworker"
)

var (
	exitChan = make(chan int)
)

// Example "Job"
type PrintSomeText struct {
	text string
}

// Should implement Execute
func (p *PrintSomeText) Execute() {
	fmt.Println(p.text)
}

func main() {
	// Make sure to listen for OS interrupts
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
		<-sigChan
		go func() {
			// Handle something before exit
			<-exitChan
			os.Exit(0)
		}()
		// If we signal again, forcefully terminate.
		<-sigChan
		os.Exit(1)
	}()

	dispatcher := gworker.NewDispatcher(
		1,   //maxWorkers int,
		100, //maxQueue int,
		50,  //maxPrioQueue int,
	)
	// starting n number of workers
	for i := 0; i < cap(dispatcher.WorkerPool); i++ {
		worker := gworker.NewWorker(dispatcher.WorkerPool)
		worker.Start()
	}

	go dispatcher.Run()

	for i := 0; i < 100; i++ {
		job := PrintSomeText{text: fmt.Sprintf("Job nr: %v", i)}
		dispatcher.JobQueue <- &job
	}
	for i := 0; i < 50; i++ {
		job := PrintSomeText{text: fmt.Sprintf("PrioJob nr: %v", i)}
		dispatcher.PriJobQueue <- &job
	}

	// Wait before exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan
	os.Exit(0)
}
