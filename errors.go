package eval

// Represents an error generated during evaluation of a query
type EvalErr struct {
	Message string
}

func NewEvalError(msg string) *EvalErr {
	return &EvalErr{
		Message: msg,
	}
}

func (e *EvalErr) Error() string {
	return e.Message
}

// Represents error caused by no results in evaluation
type UndefinedErr struct {
	Message string
}

func NewUndefinedError(msg string) *UndefinedErr {
	return &UndefinedErr{Message: msg}
}

func (e *UndefinedErr) Error() string {
	return e.Message
}
