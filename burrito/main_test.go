package burrito

import (
	"errors"
	"fmt"
	"testing"
)

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
	AssertError(t, err, "Outer error\n   [github.com/Bedrock-OSS/go-burrito/burrito.TestStackTrace] main_test.go:26\n[+]: test error")
	PrintStackTrace = false
}

func TestForceStackTrace(t *testing.T) {
	err := WrapError(errors.New("test error"), "Outer error").(*Error)
	err.ForceStackTrace(true)
	AssertError(t, err, "Outer error\n   [github.com/Bedrock-OSS/go-burrito/burrito.TestForceStackTrace] main_test.go:32\n[+]: test error")
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
	AssertError(t, err, "Outer error\n   [github.com/Bedrock-OSS/go-burrito/burrito.TestDoubleWrapErrorWithStackTrace] main_test.go:56\n[+]: Middle error\n   [github.com/Bedrock-OSS/go-burrito/burrito.TestDoubleWrapErrorWithStackTrace] main_test.go:56\n[+]: test error")
}

func TestWrapErrorf(t *testing.T) {
	err := WrapErrorf(errors.New("test error"), "Outer error %d", 1).(*Error)
	AssertError(t, err, "Outer error 1\n[+]: test error")
}

func TestWrappedErrorf(t *testing.T) {
	err := WrappedErrorf("Error %d", 1).(*Error)
	AssertError(t, err, "Error 1")
}

func TestTags(t *testing.T) {
	err := WrapError(errors.New("test error"), "Outer error").(*Error)
	err.AddTag("test")
	AssertError(t, err, "Outer error\n[+]: test error")
	AssertTags(t, err, []string{"test"})
	if err.HasTag("test2") {
		t.Fatalf("Expected no tag: test2")
	}
}

func TestNestedTags(t *testing.T) {
	err := WrapError(errors.New("test error"), "Middle error").(*Error)
	err.AddTag("test")
	err2 := WrapError(err, "Outer error").(*Error)
	err2.AddTag("test2")
	AssertError(t, err2, "Outer error\n[+]: Middle error\n[+]: test error")
	AssertTags(t, err2, []string{"test", "test2"})
	if err2.HasTag("test3") {
		t.Fatalf("Expected no tag: test3")
	}
}

func TestReadmeExample(t *testing.T) {
	err := WrappedError("This is a root error")
	err = WrapErrorf(err, "We failed to do Y and we can provide the cause")
	err = WrapErrorf(err, "We failed to do X and we can provide the cause")
	AssertError(t, err, "We failed to do X and we can provide the cause\n[+]: We failed to do Y and we can provide the cause\n[+]: This is a root error")
}

func TestReadmeExample2(t *testing.T) {
	PrintStackTrace = true
	err := WrappedError("This is a root error")
	err = WrapErrorf(err, "We failed to do Y and we can provide the cause")
	err = WrapErrorf(err, "We failed to do X and we can provide the cause")
	AssertError(t, err, "We failed to do X and we can provide the cause\n   [github.com/Bedrock-OSS/go-burrito/burrito.TestReadmeExample2] main_test.go:104\n[+]: We failed to do Y and we can provide the cause\n   [github.com/Bedrock-OSS/go-burrito/burrito.TestReadmeExample2] main_test.go:103\n[+]: This is a root error\n   [github.com/Bedrock-OSS/go-burrito/burrito.TestReadmeExample2] main_test.go:102")
	PrintStackTrace = false
}

func TestReadmeExample3(t *testing.T) {
	const ErrNotFound = "not_found"
	err := WrappedError("File not found")
	err.(*Error).AddTag(ErrNotFound)
	err = WrapErrorf(err, "We failed to do Y and we can provide the cause")
	err = WrapErrorf(err, "We failed to do X and we can provide the cause")
	AssertError(t, err, "We failed to do X and we can provide the cause\n[+]: We failed to do Y and we can provide the cause\n[+]: File not found")
	AssertTags(t, err.(*Error), []string{ErrNotFound})
}

func TestGroupErrors(t *testing.T) {
	err := GroupErrors(
		WrappedError("Error 1"),
		WrappedError("Error 2"),
		WrappedError("Error 3"),
	)
	AssertError(t, err, "Error 1\nAdditionally the following errors occurred:\n\nError 2\n\nError 3")
}

func TestWrappedGroupErrors(t *testing.T) {
	err := GroupErrors(
		WrappedError("Error 1"),
		WrappedError("Error 2"),
		WrappedError("Error 3"),
	)
	AssertError(t, WrapError(err, "Failed to do something"), "Failed to do something\n[+]: Error 1\n  >> Additionally the following errors occurred:\n  >> \n  >> Error 2\n  >> \n  >> Error 3")
}

func TestMixedErrors(t *testing.T) {
	err := WrapError(PassError(WrapError(WrappedError("test error"), "Inner error")), "Outer error")
	AssertError(t, err, "Outer error\n[+]: Inner error\n[+]: test error")
}

func TestWalkErrorStop(t *testing.T) {
	err := WrapError(PassError(WrapError(WrappedError("test error"), "Inner error")), "Outer error")
	var count int
	WalkError(err, func(err *Error) bool {
		count++
		return true
	})
	if count != 1 {
		t.Fatalf("Expected count: %d, got: %d", 1, count)
	}
}

func TestWalkError(t *testing.T) {
	err := WrapError(PassError(WrapError(WrappedError("test error"), "Inner error")), "Outer error")
	var count int
	WalkError(err, func(err *Error) bool {
		count++
		return false
	})
	if count != 4 {
		t.Fatalf("Expected count: %d, got: %d", 4, count)
	}
}

func TestProperty(t *testing.T) {
	err := WrappedError("test error")
	err.(*Error).AddProperty("test", "test value")
	if !err.(*Error).HasProperty("test") {
		t.Fatalf("Expected property: %s", "test")
	}
	if err.(*Error).GetProperty("test") != "test value" {
		t.Fatalf("Expected property value: %s", "test value")
	}
	if err.(*Error).HasProperty("test2") {
		t.Fatalf("Unexpected property: %s", "test2")
	}
	if err.(*Error).GetProperty("test2") != nil {
		t.Fatalf("Unexpected property value: %s", err.(*Error).GetProperty("test2"))
	}
}

func AssertTags(t *testing.T, err *Error, strings []string) {
	for _, tag := range strings {
		if !err.HasTag(tag) {
			t.Fatalf("Expected tag: %s", tag)
		}
	}
}

func AssertError(t *testing.T, err error, expected string) {
	if err == nil {
		t.Fatalf("Expected error: %s", expected)
	} else if err.Error() != expected {
		fmt.Printf("%s\n", err)
		t.Fatalf("Expected error: %s, got: %s", expected, err.Error())
	}
}
