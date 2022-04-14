package main

import "errors"

var (
	ErrNotNumber = errors.New("not a number")
	ErrInvalidColumn = errors.New("invalid column")
	ErrNoFiles = errors.New("no files")
	ErrInvalidOperation = errors.New("invalid operation")
)