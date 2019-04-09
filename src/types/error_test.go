package types

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestError(t *testing.T) {
	err := Error{Code: 1, Message: "Test error"}

	msg := err.Error()

	if msg != fmt.Sprintf("ERROR (%v): %s; \n Inner error: %s", err.Code, err.Message, err.Err) {
		t.Error("Wrong format error message")
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

func TestErrorConstructorFunctions(t *testing.T) {
	type args struct {
		message string
		err     error
		details []string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name        string
		args        args
		want        want
		createError func(message string, err error, details ...string) *Error
	}{
		{
			name: "internal server error",
			args: args{
				message: "internal",
				err:     errors.New("nil"),
				details: []string{"detail1", "detail2"},
			},
			want:        want{code: http.StatusInternalServerError},
			createError: InternalServerError,
		},
		{
			name: "internal server error with nil err and details",
			args: args{
				message: "internal",
				err:     nil,
				details: nil,
			},
			want:        want{code: http.StatusInternalServerError},
			createError: InternalServerError,
		},
		{
			name: "bad request",
			args: args{
				message: "bad request",
				err:     errors.New("400"),
				details: []string{"bad", "request"},
			},
			want:        want{code: http.StatusBadRequest},
			createError: BadRequestError,
		},
		{
			name: "bad request with nil err and details",
			args: args{
				message: "bad request",
				err:     nil,
				details: nil,
			},
			want:        want{code: http.StatusBadRequest},
			createError: BadRequestError,
		},
		{
			name: "not found",
			args: args{
				message: "not found",
				err:     errors.New("404"),
				details: []string{"not", "found"},
			},
			want:        want{code: http.StatusNotFound},
			createError: NotFoundError,
		},
		{
			name: "not found with nil err and details",
			args: args{
				message: "not found",
				err:     nil,
				details: nil,
			},
			want:        want{code: http.StatusNotFound},
			createError: NotFoundError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *Error
			if tt.args.details == nil {
				got = tt.createError(tt.args.message, tt.args.err)
			} else {
				got = tt.createError(tt.args.message, tt.args.err, tt.args.details...)
			}

			if got.Message != tt.args.message {
				t.Errorf("expect message '%s' but got `%s`", tt.args.message, got.Message)
			}
			if !reflect.DeepEqual(got.Details, tt.args.details) {
				t.Errorf("expect detail '%v' but got `%v`", tt.args.details, got.Details)
			}
			if !reflect.DeepEqual(got.Err, tt.args.err) {
				t.Errorf("expect Err '%#v' but got `%#v`", tt.args.err, got.Err)
			}
			if got.Code != tt.want.code {
				t.Errorf("error code = %v, want %v", got.Code, tt.want.code)
			}
		})
	}
}
