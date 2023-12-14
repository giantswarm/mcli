package cmd

import (
	"errors"
	"fmt"
)

var ErrInvalidFlag = errors.New("invalid flag")

func invalidFlagError(invalidFlag string) error {
	return fmt.Errorf("required flag %s not set.\n%w", invalidFlag, ErrInvalidFlag)
}
