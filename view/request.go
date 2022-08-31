package view

import "net/http"

type Request struct {
	Method  string      `json:"method"`
	URL     string      `json:"url"`
	Headers http.Header `json:"headers"`
	Body    interface{} `json:"body"`
}
