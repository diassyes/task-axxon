package view

import (
	"net/http"
)

type Response struct {
	UUID    string      `json:"uuid,omitempty"`
	Status  int         `json:"status,omitempty"`
	Headers http.Header `json:"headers,omitempty"`
	Length  int64       `json:"length,omitempty"`
}

func NewResponse(r *http.Response, uuid string) Response {
	return Response{
		UUID:    uuid,
		Headers: r.Header,
		Status:  r.StatusCode,
		Length:  r.ContentLength,
	}
}
