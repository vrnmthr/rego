package rego

import (
	"fmt"
	"strings"
	"bufio"
)

// EvalErr represents error generated during evaluation of a query
type EvalErr struct {
	Message string
}

// NewEvalError creates a new EvalError with message msg
func NewEvalError(msg string) *EvalErr {
	return &EvalErr{
		Message: msg,
	}
}

func (e *EvalErr) Error() string {
	return e.Message
}

// UndefinedErr represents error caused by no results in evaluation
type UndefinedErr struct {
	Message string
}

// NewUndefinedError creates a new UndefinedErr with the given message
func NewUndefinedError(msg string) *UndefinedErr {
	return &UndefinedErr{Message: msg}
}

func (e *UndefinedErr) Error() string {
	return e.Message
}

// Errors represents multiple errors
type Errors []error

func (e *Errors) Add(err error) {
	if err != nil {
		switch v := err.(type) {
		case *Errors:
			*e = append(*e, *v...)
		default:
			*e = append(*e, v)
		}
	}
}

func (e *Errors) Error() string {
	buf := make([]string, 0)
	for _, e := range *e {
		buf = append(buf, e.Error())
	}
	return fmt.Sprintf("%v errors:\n%v", len(buf), strings.Join(buf, "\n"))
}

func (e *Errors) NilIfEmpty() error {
	if len(*e) == 0 {
		return nil
	}
	return e
}

// ErrWriter wraps a writer
type ErrWriter struct {
	err error
	w   *bufio.Writer
}

func (ew *ErrWriter) Write(data interface{}) {

	if ew.err != nil {
		return
	}

	var err error
	switch val := data.(type) {
	case []byte:
		_, err = ew.w.Write(val)
	default:
		_, err = ew.w.Write([]byte(fmt.Sprintf("%v", val)))
	}

	ew.err = err
}

func (ew *ErrWriter) Flush() {
	if ew.err != nil {
		return
	}
	ew.err = ew.w.Flush()
}

func (ew *ErrWriter) Error() error {
	return ew.err
}