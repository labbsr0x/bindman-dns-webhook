package types

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestCheckDNSRecord(t *testing.T) {
	r := DNSRecord{Name: "t.test.com", Value: "0.0.0.0", Type: "A", ClusterID: "c1"}
	result, errs := r.Check()

	if !result && len(errs) == 0 {
		t.Errorf("Got non-sucess execution but no errors reported for r := '%v'", r)
	}
	if !result {
		t.Errorf("Expecting success. Got the otherwise for r := '%v'; errs := '%v'", r, errs)
	}

	rs := []DNSRecord{DNSRecord{ClusterID: "c1", Value: "0.0.0.0", Type: "A"}, DNSRecord{Value: "0.0.0.0", Type: "A"}, DNSRecord{Type: "A"}, DNSRecord{}}

	for i, r := range rs {
		result, errs = r.Check()

		if !result && len(errs) == 0 {
			t.Errorf("Got non-sucess execution but no errors reported for r := '%v'", r)
		}

		if result {
			t.Errorf("Expecting error. Got success for r := '%v", r)
		}

		if len(errs) > i+1 {
			t.Errorf("Expecting exactly 1 error for '%v'. Got %v. errs := '%v", r, len(errs), errs)
		}
	}
}

func TestError(t *testing.T) {
	err := Error{Code: 1, Err: fmt.Errorf("Test"), Message: "Test error"}

	msg := err.Error()

	if strings.Contains(msg, "Inner error") {
		t.Errorf("When BINDMAN_MODE not defined, expecting the inner error to be hidden")
	}

	os.Setenv("BINDMAN_MODE", "DEBUG")

	msg = err.Error()

	if !strings.Contains(msg, "Inner error") {
		t.Errorf("When BINDMAN_MODE is set to DEBUG, expecting the inner error to be shown")
	}
}

func TestPanicIfErrorSucc(t *testing.T) {
	err := Error{Code: 1, Err: fmt.Errorf("Test"), Message: "Test error"}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expecting the code to panic")
		}
	}()

	PanicIfError(err)
}

func TestPanicIfErrorErr(t *testing.T) {
	err := Error{Code: 1, Message: "Test Error"}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expecting the code to not panic.")
		}
	}()

	PanicIfError(err)
}

func TestPanic(t *testing.T) {
	err := Error{}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expecting the code to panic")
		}
	}()

	Panic(err)
}
