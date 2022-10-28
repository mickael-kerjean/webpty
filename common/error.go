package common

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("Page Not Found")
	ErrNotAuthorized = errors.New("Not Authorized")
	ErrNotAvailable  = errors.New("Not Available")
	ErrNotValid      = errors.New("Not Valid")
)
