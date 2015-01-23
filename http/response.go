package http

import (
	"net/http"
)

// an HTTP response object
type Response http.Response

// downcasts the response to the native HTTP interface
func (r *Response) Native() *http.Response {
	return (*http.Response)(r)
}

// creates a new HTTP response object
func NewResponse(response *http.Response) *Response {
	return (*Response)(response)
}
