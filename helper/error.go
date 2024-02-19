package helper

import "errors"

var (
	ErrInternal = errors.New("sorry, there's an error on our side")
)

type ResponseError struct {
	err  error
	code int
}

func NewResponseError(err error, code int) ResponseError {
	return ResponseError{err, code}
}

func (re ResponseError) Error() string {
	return re.err.Error()
}

func (re ResponseError) Code() int {
	return re.code
}
