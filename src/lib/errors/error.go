package errors

import (
	"fmt"
	"runtime"
)

type customError struct {
	message    string
	errMessage string
	filename   string
	line       int
	function   string
	cause      error
}

func (e *customError) Error() string {
	return fmt.Sprint(e.message)
}

func NewError(msg string, errMsg string, val ...interface{}) error {
	return create(nil, msg, errMsg, val...)
}

func GetCaller(err error) (string, int, string, error) {
	ce, ok := err.(*customError)
	if !ok {
		fmt.Println("not ok")
		return "", 0, "", create(nil, "something wrong", "failed to cast to custom error")
	}

	return ce.filename, ce.line, ce.errMessage, nil
}

func create(cause error, msg string, errMsg string, val ...interface{}) error {
	err := &customError{
		message:    msg,
		errMessage: fmt.Sprintf(errMsg, val...),
		cause:      cause,
	}

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return err
	}
	err.filename, err.line = file, line

	f := runtime.FuncForPC(pc)
	if f == nil {
		return err
	}
	err.function = f.Name()

	return err
}
