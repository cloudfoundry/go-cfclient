package resource

import "time"

type Job struct {
	GUID      string              `json:"guid"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Operation string              `json:"operation"` // Current desired operation of the job on a model
	State     string              `json:"state"`     // State of the job; valid values are PROCESSING, POLLING, COMPLETE, or FAILED
	Errors    []CloudFoundryError `json:"errors"`    // Array of errors that occurred while processing the job
	Warnings  []JobWarning        `json:"warnings"`  // Array of warnings that occurred while processing the job
	Links     map[string]Link     `json:"links"`
}

type JobWarning struct {
	Detail string `json:"detail"`
}
