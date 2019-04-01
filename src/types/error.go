package types

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Error groups together information that defines an error. Should always be used to
type Error struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Details interface{} `json:"details"`
}

// Error() gives a string representing the error; also, forces the Error type to comply with the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("ERROR (%d): %s", e.Code, e.Message)
}

// BadRequestError create an Error instance with http.StatusBadRequest code
func BadRequestError(message string, details interface{}) *Error {
	return &Error{message, http.StatusBadRequest, details}
}

// BadRequestError create an Error instance with http.StatusNotFound code
func NotFoundError(message string, details interface{}) *Error {
	return &Error{message, http.StatusNotFound, details}
}

// BadRequestError create an Error instance with http.StatusInternalServerError code
func InternalServerError(message string, details interface{}) *Error {
	return &Error{message, http.StatusInternalServerError, details}
}

// PanicIfError is just a wrapper to a panic call that propagates error when it's not nil
func PanicIfError(e error) {
	if e != nil { /**/
		logrus.Errorf(e.Error())
		panic(e)
	}
}

// Panic wraps a panic call propagating the given error parameter
func Panic(e Error) {
	logrus.Errorf(e.Error())
	panic(e)
}
