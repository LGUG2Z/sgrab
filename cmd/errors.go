package cmd

import (
	"errors"
	"fmt"
)

var (
	ErrCouldNotFindSeries = func(series string) error {
		return fmt.Errorf("Could not find series '%s' on seedbox.", series)
	}
	ErrCredentialsRejected            = errors.New("Credentials rejected by seedbox.")
	ErrInterruptReceived              = errors.New("Received an interrupt. Cleanup successful.")
	ErrInterruptReceivedCleanupFailed = errors.New("Received an interrupt. Cleanup failed.")
	ErrInformationMissing             = errors.New("Required information missing. See 'sgrab --help'.")
)
