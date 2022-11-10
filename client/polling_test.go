package client

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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

	failedFn := func() (string, error) {
		return "FAILED", nil
	}
	successFn := func() (string, error) {
		return "SUCCESS", nil
	}
	timeoutFn := func() (string, error) {
		return "PROCESSING", nil
	}

	err := PollForStateOrTimeout(failedFn, "NOPE", noWaitOpts)
	require.Equal(t, AsyncProcessFailedError, err)

	err = PollForStateOrTimeout(successFn, "SUCCESS", noWaitOpts)
	require.NoError(t, err)

	err = PollForStateOrTimeout(timeoutFn, "SUCCESS", noWaitOpts)
	require.Equal(t, AsyncProcessTimeoutError, err)
}
