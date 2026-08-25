package dberr

import "errors"

var (
	ErrNotFound      = errors.New("goorm: record not found")
	ErrAlreadyExists = errors.New("goorm: record already exists")
)

var (
	Is = errors.Is
	As = errors.As
)
