package burrito

import (
	"errors"
	"fmt"
	"testing"
)

func AssertError(t *testing.T, err error, expected string) {
	if err == nil {
		t.Fatalf("Expected error: %s", expected)
	} else if err.Error() != expected {
		fmt.Printf("%s\n", err)
		t.Fatalf("Expected error: %s, got: %s", expected, err.Error())
	}
}

func TestWrapError(t *testing.T) {
	err := WrapError(errors.New("test error"), "Outer error")
	AssertError(t, err, "Outer error\n[+]: test error")
}

func TestDoubleWrapError(t *testing.T) {
	err := WrapError(WrapError(errors.New("test error"), "Middle error"), "Outer error")
	AssertError(t, err, "Outer error\n[+]: Middle error\n[+]: test error")
}

func TestMultilineWrapError(t *testing.T) {
	err := WrapError(errors.New("test error"), "Outer\nerror")
	AssertError(t, err, "Outer\n   >> error\n[+]: test error")
}

func TestStackTrace(t *testing.T) {
	PrintStackTrace = true
	err := WrapError(errors.New("test error"), "Outer error")
	AssertError(t, err, "Outer error\n   [github.com/Bedrock-OSS/go-burrito/burrito.TestStackTrace] main_test.go:35\n[+]: test error")
	PrintStackTrace = false
}

func TestForceStackTrace(t *testing.T) {
	err := WrapError(errors.New("test error"), "Outer error").(*Error)
	err.ForceStackTrace(true)
	AssertError(t, err, "Outer error\n   [github.com/Bedrock-OSS/go-burrito/burrito.TestForceStackTrace] main_test.go:41\n[+]: test error")
}

func TestDisableStackTrace(t *testing.T) {
	PrintStackTrace = true
	err := WrapError(errors.New("test error"), "Outer error").(*Error)
	err.ForceStackTrace(false)
	AssertError(t, err, "Outer error\n[+]: test error")
	PrintStackTrace = false
}

func TestPassError(t *testing.T) {
	err := WrapError(PassError(errors.New("test error")), "Outer error")
	AssertError(t, err, "Outer error\n[+]: test error")
}

func TestWrappedError(t *testing.T) {
	err := WrapError(WrappedError("test error"), "Outer error")
	AssertError(t, err, "Outer error\n[+]: test error")
}

func TestDoubleWrapErrorWithStackTrace(t *testing.T) {
	err := WrapError(WrapError(errors.New("test error"), "Middle error"), "Outer error").(*Error)
	err.ForceStackTrace(true)
	AssertError(t, err, "Outer error\n   [github.com/Bedrock-OSS/go-burrito/burrito.TestDoubleWrapErrorWithStackTrace] main_test.go:65\n[+]: Middle error\n   [github.com/Bedrock-OSS/go-burrito/burrito.TestDoubleWrapErrorWithStackTrace] main_test.go:65\n[+]: test error")
}

func TestWrapErrorf(t *testing.T) {
	err := WrapErrorf(errors.New("test error"), "Outer error %d", 1).(*Error)
	AssertError(t, err, "Outer error 1\n[+]: test error")
}

func TestWrappedErrorf(t *testing.T) {
	err := WrappedErrorf("Error %d", 1).(*Error)
	AssertError(t, err, "Error 1")
}
