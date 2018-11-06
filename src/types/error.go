package types

import (
	"fmt"
	"os"
	"strings"
)

// Error groups together information that defines an error. Should always be used to
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Err     error  `json:"err"`
}

// Error() gives a string representing the error; also, forces the Error type to comply with the error interface
func (e *Error) Error() string {
	mode := strings.Trim(os.Getenv("MODE"), " ")
	if mode == "DEBUG" || mode == "" {
		return fmt.Sprintf("ERROR (%v): %s; \n Inner error: %s", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("ERROR (%v): %s", e.Code, e.Message)
}

// PanicIfError is just a wrapper to a panic call that propagates a custom Error when the err property != nil
func PanicIfError(e Error) {
	if e.Err != nil {
		panic(e)
	}
}

// Panic wraps a panic call propagating the given error parameter
func Panic(e Error) {
	panic(e)
}
