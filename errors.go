package rego

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
