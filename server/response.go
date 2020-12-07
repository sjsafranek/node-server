package server

import (
	"encoding/json"
	"fmt"
	"io"
)

func NewResponse() *Response {
	return &Response{Status: "ok", StatusCode: 200}
}

type Response struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	Error      string      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func (self *Response) Marshal() (string, error) {
	b, err := json.Marshal(self)
	if nil != err {
		return "", err
	}
	return string(b), err
}

func (self *Response) SetData(data interface{}) {
	self.Data = data
}

func (self *Response) SetError(err error) {
	self.Status = "error"
	self.Error = err.Error()
}

func (self *Response) SetStatusCode(statusCode int) {
	self.StatusCode = statusCode
}

func (self *Response) Write(w io.Writer) error {
	payload, err := self.Marshal()
	if nil != err {
		return err
	}
	_, err = fmt.Fprintln(w, payload)
	return err
}

type Request struct {
	Method string        `json:"method"`
	Params RequestParams `json:"params"`
}

type RequestParams struct {
	BucketId string `json:"bucket_id"`
	Key      string `json:"key"`
	Value    []byte `json:"value"`
}
