package burrito

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"path/filepath"
	"runtime"
	"strings"
)

type Error struct {
	Err            error
	Message        *string
	File           string
	Line           int
	FuncName       string
	ShowStackTrace *bool
	Tags           []string
}

type ErrorGroup struct {
	Errors []error
}

func (r *ErrorGroup) Error() string {
	text := fmt.Sprintf("%s\n%s", r.Errors[0], GroupErrorText)
	for i := 1; i < len(r.Errors); i++ {
		text = fmt.Sprintf("%s\n\n%s", text, r.Errors[i].Error())
	}
	return text
}

func (r *Error) Error() string {
	shouldPrintStackTrace := PrintStackTrace
	if r.ShowStackTrace != nil {
		shouldPrintStackTrace = *r.ShowStackTrace
	}
	text := ""
	var err error = r
	for err != nil {
		if e, ok := err.(*Error); ok {
			if e.Message != nil {
				if e == r {
					text = fmt.Sprintf("%s%s", text, strings.Replace(*e.Message, "\n", color.YellowString("\n   >> "), -1))
				} else {
					text = fmt.Sprintf("%s\n[%s]: %s", text, color.RedString("+"), strings.Replace(*e.Message, "\n", color.YellowString("\n  >> "), -1))
				}
			}
			if shouldPrintStackTrace {
				text = fmt.Sprintf("%s\n   [%s] %s:%d", text, e.FuncName, filepath.Base(e.File), e.Line)
			}
			err = e.Err
		} else {
			text = fmt.Sprintf("%s\n[%s]: %s", text, color.RedString("+"), strings.Replace(err.Error(), "\n", color.YellowString("\n  >> "), -1))
			break
		}
	}
	return text
}

// AddTag adds a tag to the error.
func (r *Error) AddTag(tag string) {
	r.Tags = append(r.Tags, tag)
}

// HasTag returns true if the error has the specified tag.
func (r *Error) HasTag(tag string) bool {
	e := r
	for e != nil {
		for _, t := range e.Tags {
			if t == tag {
				return true
			}
		}
		if e.Err != nil {
			if err, ok := e.Err.(*Error); ok {
				e = err
				continue
			}
		}
		break
	}
	return false
}

// ForceStackTrace overrides the global PrintStackTrace setting.
func (r *Error) ForceStackTrace(enabled bool) {
	r.ShowStackTrace = &enabled
}

// PrintStackTrace is a global variable that controls whether stack traces are printed or not.
var PrintStackTrace = false

// GroupErrorText is a string that is used to indicate a group of errors.
var GroupErrorText = color.RedString("Additionally the following errors occurred:")

// wrapErrorStackTrace is used by other wrapped error functions to add a stack
// trace to the error message.
func wrapErrorStackTrace(err error, text string) error {
	pc, fn, line, _ := runtime.Caller(2)
	return &Error{
		Err:      err,
		Message:  &text,
		File:     filepath.Base(fn),
		Line:     line,
		FuncName: runtime.FuncForPC(pc).Name(),
	}
}

// wrapErrorHandlerErrorStackTrace is a helper function for wrapping errors
// that occurred during error handling.
func wrapErrorHandlerErrorStackTrace(
	errs ...error,
) error {
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}
	return &ErrorGroup{Errors: errs}
}

// PassError adds stack trace to an error without any additional text.
func PassError(err error) error {
	text := err.Error()
	if PrintStackTrace {
		pc, fn, line, _ := runtime.Caller(1)
		text = fmt.Sprintf(
			"%s\n   [%s] %s:%d", text, runtime.FuncForPC(pc).Name(),
			filepath.Base(fn), line)
	}
	return errors.New(text)
}

// WrappedError creates an error with a stack trace from text.
func WrappedError(text string) error {
	return wrapErrorStackTrace(nil, text)
}

// WrappedErrorf creates an error with a stack trace from formatted text.
func WrappedErrorf(text string, args ...interface{}) error {
	text = fmt.Sprintf(text, args...)
	return wrapErrorStackTrace(nil, text)
}

// WrapError wraps an error with a stack trace and adds additional text
// information.
func WrapError(err error, text string) error {
	return wrapErrorStackTrace(err, text)
}

// WrapErrorf wraps an error with a stack trace and adds additional formatted
// text information.
func WrapErrorf(err error, text string, args ...interface{}) error {
	return wrapErrorStackTrace(err, fmt.Sprintf(text, args...))
}

// GroupErrors combines two or more errors into one. The first error is
// an error that occurred during the main operation. The other errors are
// errors that occurred during error handling.
func GroupErrors(errs ...error) error {
	return wrapErrorHandlerErrorStackTrace(errs...)
}
