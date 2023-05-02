package utility

import (
	"fmt"
	"runtime"
	"strings"
)

type errorImpl struct {
	errors  []error
	stack   string
	message string
}

type Error interface {
	AddMessage(message string) Error
	AddError(err error) Error
	String() string
	Is(target error) bool
}

func NewError(err error) Error {
	data := make([]byte, 10000)
	runtime.Stack(data, false)
	return &errorImpl{
		stack:   strings.TrimRight(string(data), "\x00"),
		errors:  []error{err},
		message: err.Error(),
	}
}

func (e *errorImpl) String() string {
	return fmt.Sprintf("%s\nStack trace: %s", e.message, e.stack)
}

func (e *errorImpl) AddMessage(message string) Error {
	e.message = fmt.Sprintf("%s: %s", message, e.message)
	return e
}

func (e *errorImpl) AddError(err error) Error {
	e.message = fmt.Sprintf("%s: %s", err.Error(), e.message)
	e.errors = append(e.errors, err)
	return e
}

func (e *errorImpl) Is(target error) bool {
	for i := range e.errors {
		if e.errors[i] == target {
			return true
		}
	}

	return false
}
