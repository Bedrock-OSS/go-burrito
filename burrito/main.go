package burrito

import (
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
	Properties     map[string]interface{}
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
	result := false
	WalkError(r, func(err *Error) bool {
		for _, t := range err.Tags {
			if t == tag {
				result = true
				return true
			}
		}
		return false
	})
	return result
}

// HasTag returns true if the error is burrito error and has the specified tag.
func HasTag(err error, tag string) bool {
	if e, ok := err.(*Error); ok {
		return e.HasTag(tag)
	}
	return false
}

// GetAllMessages returns all messages of the error and all wrapped errors.
func GetAllMessages(err error) []string {
	var messages []string
	e := err
	for e != nil {
		if err1, ok := e.(*Error); ok {
			if err1.Message != nil {
				messages = append(messages, *err1.Message)
			}
			e = err1.Err
		} else {
			messages = append(messages, e.Error())
			break
		}
	}
	return messages
}

func (r *Error) AddProperty(key string, value interface{}) {
	if r.Properties == nil {
		r.Properties = make(map[string]interface{})
	}
	r.Properties[key] = value
}

func (r *Error) GetProperty(key string) interface{} {
	if r.Properties == nil {
		return nil
	}
	return r.Properties[key]
}

func (r *Error) HasProperty(key string) bool {
	if r.Properties == nil {
		return false
	}
	_, ok := r.Properties[key]
	return ok
}

// GetProperty returns the property value for the specified key.
func GetProperty(err error, key string) interface{} {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e.GetProperty(key)
	}
	return nil
}

// WalkError walks the error tree and calls the specified function for each error.
// If the function returns true, the walk is aborted.
func WalkError(err error, fn func(err *Error) bool) {
	for err != nil && IsBurritoError(err) {
		e := AsBurritoError(err)
		if fn(e) {
			return
		}
		err = e.Err
	}
}

// IsBurritoError returns true if the error is a burrito error.
func IsBurritoError(err error) bool {
	_, ok := err.(*Error)
	return ok
}

// AsBurritoError converts an error to a burrito error.
func AsBurritoError(err error) *Error {
	if e, ok := err.(*Error); ok {
		return e
	}
	return nil
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
func wrapErrorStackTrace(err error, text *string) error {
	pc, fn, line, _ := runtime.Caller(2)
	return &Error{
		Err:      err,
		Message:  text,
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
	return wrapErrorStackTrace(err, nil)
}

// WrappedError creates an error with a stack trace from text.
func WrappedError(text string) error {
	return wrapErrorStackTrace(nil, &text)
}

// WrappedErrorf creates an error with a stack trace from formatted text.
func WrappedErrorf(text string, args ...interface{}) error {
	text = fmt.Sprintf(text, args...)
	return wrapErrorStackTrace(nil, &text)
}

// WrapError wraps an error with a stack trace and adds additional text
// information.
func WrapError(err error, text string) error {
	return wrapErrorStackTrace(err, &text)
}

// WrapErrorf wraps an error with a stack trace and adds additional formatted
// text information.
func WrapErrorf(err error, text string, args ...interface{}) error {
	text = fmt.Sprintf(text, args...)
	return wrapErrorStackTrace(err, &text)
}

// GroupErrors combines two or more errors into one. The first error is
// an error that occurred during the main operation. The other errors are
// errors that occurred during error handling.
func GroupErrors(errs ...error) error {
	return wrapErrorHandlerErrorStackTrace(errs...)
}
