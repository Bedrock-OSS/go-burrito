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
					text = fmt.Sprintf("%s\n[%s]: %s", text, color.RedString("+"), strings.Replace(*e.Message, "\n", color.YellowString("\n   >> "), -1))
				}
			}
			if shouldPrintStackTrace {
				text = fmt.Sprintf("%s\n   [%s] %s:%d", text, e.FuncName, filepath.Base(e.File), e.Line)
			}
			err = e.Err
		} else {
			text = fmt.Sprintf("%s\n[%s]: %s", text, color.RedString("+"), err.Error())
			err = nil
		}
	}
	return text
}

func (r *Error) ForceStackTrace(enabled bool) {
	r.ShowStackTrace = &enabled
}

var PrintStackTrace = false

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
//
// - mainErr - the error that is being handled. The message of this error
//   should be properly formatted and have the stack trace if debug mode is
//   enabled.
// - handlerErr - the error that occurred during handling. This value can be
//   nil. In this case only the errorHandlerText is used for printing the
//   part of the message related to the error handler.
// - connectorText - text used to connect both errors. For example:
//   "Another error occurred while handling the previous error:". This text can
//   be empty. In this case the errors are separated by two new lines.
// - errorHandlerText - additional text to be added to the error message. This
//   text can be empty. IN this case only the handlerErr is used for printing
//   the part of the message related to the error handler.
func wrapErrorHandlerErrorStackTrace(
	mainErr, handlerErr error, connectorText, errorHandlerText string,
) error {
	// Add header (the main message)
	text := mainErr.Error() + "\n\n"
	// Add connector text (optional)
	if connectorText != "" {
		text = text + connectorText + "\n\n"
	}
	// Format and add the error handler error
	errorHandlerText = strings.Replace(
		errorHandlerText, "\n", color.YellowString("\n   >> "), -1)
	redPlus := color.RedString("+")
	if handlerErr == nil {
		if errorHandlerText != "" {
			errorHandlerText = fmt.Sprintf(
				"[%s]: %s", redPlus, errorHandlerText)
		}
		// else: no error, but this function shouldn't be used like this
		// no extra text. But it's possible that it will leave connector text
		// at the end.
	} else {
		if errorHandlerText != "" {
			errorHandlerText = fmt.Sprintf(
				"[%s]: %s\n[%s]: %s", redPlus, errorHandlerText, redPlus,
				handlerErr.Error())
		} else {
			errorHandlerText = fmt.Sprintf(
				"[%s]: %s\n[%s]: %s", redPlus, errorHandlerText, redPlus,
				handlerErr.Error())
		}
	}
	text = text + errorHandlerText
	// Add stack trace (optional)
	if PrintStackTrace {
		pc, fn, line, _ := runtime.Caller(2)
		text = fmt.Sprintf(
			"%s\n   [%s] %s:%d", text, runtime.FuncForPC(pc).Name(),
			filepath.Base(fn), line)
	}
	return errors.New(text)
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

// WrapErrorHandlerError combines two errors into one. The first error is
// an error that occurred during the main operation. The second error is an
// error that occurred during error handling. Errors are combined using
// connectorText. Additional text can be added to the handler error message
// using errorHandlerText.
func WrapErrorHandlerError(
	mainErr, handlerErr error, connectorText, errorHandlerText string,
) error {
	return wrapErrorHandlerErrorStackTrace(
		mainErr, handlerErr, connectorText, errorHandlerText)
}

// PassErrorHandlerError combines mainErr and handlerError similar to
// WrapErrorHandlerError, but it doesn't provide any additional text
// (analogous to PassError).
func PassErrorHandlerError(mainErr, handlerErr error, connectorText string) error {
	return wrapErrorHandlerErrorStackTrace(mainErr, handlerErr, connectorText, "")
}
