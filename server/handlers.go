package server

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/sjsafranek/lemur"
)

func httpError(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	var resp Response
	resp.SetError(errors.New(http.StatusText(http.StatusNotFound)))
	resp.Write(w)
}

func (self *Server) Do(request *Request) *Response {
	var response Response

	switch request.Method {

	case "get":
		bucket, err := self.db.Get(request.Params.BucketId)
		if nil != err {
			response.SetStatusCode(http.StatusNotFound)
			response.SetError(errors.New(http.StatusText(http.StatusNotFound)))
			return &response
		}

		value, err := bucket.Get(request.Params.Key)
		if nil != err {
			response.SetStatusCode(http.StatusInternalServerError)
			response.SetError(errors.New(http.StatusText(http.StatusInternalServerError)))
			return &response
		}

		response.SetData(value)

	case "set":
		bucket, err := self.db.Get(request.Params.BucketId)
		if nil != err {
			response.SetStatusCode(http.StatusNotFound)
			response.SetError(errors.New(http.StatusText(http.StatusNotFound)))
			return &response
		}

		bucket.Set(request.Params.Key, request.Params.Value)

	default:
		response.SetStatusCode(http.StatusMethodNotAllowed)
		response.SetError(errors.New(http.StatusText(http.StatusMethodNotAllowed)))
	}

	return &response
}

func (self *Server) HttpGetKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := lemur.Vars(r)

	response := self.Do(&Request{
		Method: "get",
		Params: RequestParams{
			BucketId: vars["bucketId"],
			Key:      vars["key"],
		},
	})

	if 200 != response.StatusCode {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		response.Write(w)
		return
	}

	w.WriteHeader(response.StatusCode)
	w.Write(response.Data.([]byte))
}

func (self *Server) HttpSetKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := lemur.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if nil != err {
		var response Response
		response.SetStatusCode(http.StatusBadRequest)
		response.SetError(errors.New(http.StatusText(http.StatusInternalServerError)))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		response.Write(w)
		return
	}
	defer r.Body.Close()

	response := self.Do(&Request{
		Method: "set",
		Params: RequestParams{
			BucketId: vars["bucketId"],
			Key:      vars["key"],
			Value:    body,
		},
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	response.Write(w)
}
