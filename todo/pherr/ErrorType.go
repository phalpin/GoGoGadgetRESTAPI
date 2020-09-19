package pherr

import "net/http"

type ErrorType int64

const (
	Unknown ErrorType = 0 + iota
	NotFound
	InvalidArgument
)

func (et ErrorType) ToHttpStatusCode() int {
	switch et {
	case NotFound:
		return http.StatusNotFound
	case InvalidArgument:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
