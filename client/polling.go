package client

import (
	"errors"
	"fmt"
	"time"
)

var ErrAsyncProcessTimeout = errors.New("timed out after waiting for async process")

type PollingOptions struct {
	Timeout       time.Duration
	CheckInterval time.Duration
	FailedState   string
}

func NewPollingOptions() *PollingOptions {
	return &PollingOptions{
		FailedState:   "FAILED",
		Timeout:       time.Minute * 5,
		CheckInterval: time.Second,
	}
}

type getStateFunc func() (string, string, error)

func PollForStateOrTimeout(getState getStateFunc, successState string, opts *PollingOptions) error {
	if opts == nil {
		opts = NewPollingOptions()
	}

	timeout := time.After(opts.Timeout)
	ticker := time.NewTicker(opts.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return ErrAsyncProcessTimeout
		case <-ticker.C:
			state, cfError, err := getState()
			if err != nil {
				return err
			}
			switch state {
			case successState:
				return nil
			case opts.FailedState:
				return fmt.Errorf("received state %s while waiting for async process: %s", opts.FailedState, cfError)
			}
		}
	}
}
