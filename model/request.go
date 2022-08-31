package model

import (
	"task-axxon/view"
)

// Request - implemented in order to save data about request and response to DB
type Request struct {
	Request  *view.Request  `json:"request"`
	Response *view.Response `json:"response"`
}

func NewRequest(req *view.Request, resp *view.Response) *Request {
	return &Request{
		Request:  req,
		Response: resp,
	}
}
