package github

import "errors"

var ErrNotFound = errors.New("not found")
var ErrInvalidFormat = errors.New("invalid format")

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
