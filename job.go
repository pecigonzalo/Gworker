package Gworker

// Job represents the job to be run
// A Job should implement an Execute method
type Job interface {
	Execute()
}
