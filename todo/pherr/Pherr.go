package pherr

import (
	"bytes"
	"net/http"
	"runtime/debug"
)

type baseError struct {
	Message    string
	StackTrace string
	ErrorType  ErrorType
}

type KnownError struct {
	*baseError
	FriendlyMessage string
}

func NewKnown(err error, friendlyMsg string, opts ...Option) *KnownError {
	options := &options{}

	for _, val := range opts {
		val(options)
	}

	buf := new(bytes.Buffer)
	buf.Write(debug.Stack())
	retVal := &KnownError{
		baseError: &baseError{
			Message:    err.Error(),
			StackTrace: buf.String(),
			ErrorType:  options.errorType,
		},
		FriendlyMessage: friendlyMsg,
	}

	return retVal
}

func (ke *KnownError) WriteHttpResponse(w http.ResponseWriter) {
	http.Error(w, ke.FriendlyMessage, ke.ErrorType.ToHttpStatusCode())
}

func (be *baseError) Error() string {
	return be.Message
}
