package types

import (
	"fmt"
	"net/http"
	"testing"
)

func TestError(t *testing.T) {
	err := Error{Code: 1, Message: "Test error"}

	msg := err.Error()

	if msg != fmt.Sprintf("ERROR (%d): %s", err.Code, err.Message) {
		t.Error("Wrong format error message")
	}
}

func TestBadRequestError(t *testing.T) {
	err := BadRequestError("test", nil)
	if err.Code != http.StatusBadRequest {
		t.Errorf("Expecting 400 but got %d", err.Code)
	}
	if err.Message != "test" {
		t.Errorf("The error message must the same as passed to constructor function. Expecting message 'test' but got %s", err.Message)
	}
	if err.Details != nil {
		t.Errorf("The error details must the same as passed to constructor function. Expecting 'nil' but got %s", err.Details)
	}
}
func TestNotFoundError(t *testing.T) {
	err := NotFoundError("test", nil)
	if err.Code != http.StatusNotFound {
		t.Errorf("Expecting 404 but got %d", err.Code)
	}
	if err.Message != "test" {
		t.Errorf("The error message must the same as passed to constructor function. Expecting message 'test' but got %s", err.Message)
	}
	if err.Details != nil {
		t.Errorf("The error details must the same as passed to constructor function. Expecting 'nil' but got %s", err.Details)
	}
}
func TestInternalServerError(t *testing.T) {
	err := InternalServerError("test", nil)
	if err.Code != http.StatusInternalServerError {
		t.Errorf("Expecting 500 but got %d", err.Code)
	}
	if err.Message != "test" {
		t.Errorf("The error message must the same as passed to constructor function. Expecting message 'test' but got %s", err.Message)
	}
	if err.Details != nil {
		t.Errorf("The error details must the same as passed to constructor function. Expecting 'nil' but got %s", err.Details)
	}
}

func TestPanicIfErrorSucc(t *testing.T) {
	err := fmt.Errorf("Test")
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expecting the code to panic")
		}
		if r != err {
			t.Error("Expecting the error to be the same as passed to the PanicIfError function")
		}
	}()
	PanicIfError(err)
}

func TestPanicIfErrorErr(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Expecting the code to not panic.")
		}
	}()
	PanicIfError(nil)
}

func TestPanic(t *testing.T) {
	err := Error{}
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expecting the code to panic")
		}
	}()
	Panic(err)
}
