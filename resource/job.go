package resource

type JobState string

// The 3 lifecycle states
const (
	JobStateProcessing JobState = "PROCESSING"
	JobStatePolling    JobState = "POLLING"
	JobStateComplete   JobState = "COMPLETE"
	JobStateFailed     JobState = "FAILED"
)

type Job struct {
	Operation string              `json:"operation"` // Current desired operation of the job on a model
	State     JobState            `json:"state"`     // State of the job; valid values are PROCESSING, POLLING, COMPLETE, or FAILED
	Errors    []CloudFoundryError `json:"errors"`    // Array of errors that occurred while processing the job
	Warnings  []JobWarning        `json:"warnings"`  // Array of warnings that occurred while processing the job
	Resource  `json:",inline"`
}

type JobWarning struct {
	Detail string `json:"detail"`
}
