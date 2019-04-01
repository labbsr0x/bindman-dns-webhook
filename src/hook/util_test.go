package hook

import (
	"net/http"
	"testing"
)

func Test_writeJSONResponse(t *testing.T) {
	type args struct {
		payload    interface{}
		statusCode int
		w          http.ResponseWriter
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeJSONResponse(tt.args.payload, tt.args.statusCode, tt.args.w)
		})
	}
}

func Test_handleError(t *testing.T) {
	type args struct {
		w http.ResponseWriter
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleError(tt.args.w)
		})
	}
}
