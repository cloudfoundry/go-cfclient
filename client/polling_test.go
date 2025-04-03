package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var CustomStagingErr = "StagingError - Staging error: Start command not specified"

func TestNewPollingOptions(t *testing.T) {
	opts := NewPollingOptions()
	require.Equal(t, "FAILED", opts.FailedState)
	require.Equal(t, time.Minute*5, opts.Timeout)
	require.Equal(t, time.Second, opts.CheckInterval)
}

func TestPollForStateOrTimeout(t *testing.T) {
	noWaitOpts := NewPollingOptions()
	noWaitOpts.Timeout = time.Second
	noWaitOpts.CheckInterval = time.Millisecond

	failedFn := func() (string, string, error) {
		return "FAILED", CustomStagingErr, nil
	}
	successFn := func() (string, string, error) {
		return "SUCCESS", "", nil
	}
	timeoutFn := func() (string, string, error) {
		return "PROCESSING", "", nil
	}

	err := PollForStateOrTimeout(failedFn, "NOPE", noWaitOpts)
	require.Error(t, err)
	require.Equal(t, "received state FAILED while waiting for async process: "+CustomStagingErr, err.Error())

	err = PollForStateOrTimeout(successFn, "SUCCESS", noWaitOpts)
	require.NoError(t, err)

	err = PollForStateOrTimeout(timeoutFn, "SUCCESS", noWaitOpts)
	require.Equal(t, ErrAsyncProcessTimeout, err)
}
