package Gworker

// Job represents the job to be run
type Job interface {
	Execute()
}
