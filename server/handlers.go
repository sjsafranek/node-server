package server

import (
	"errors"
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

func (self *Server) HttpGetKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := lemur.Vars(r)

	bucketId := vars["bucketId"]
	bucket, err := self.db.Get(bucketId)
	if nil != err {
		httpError(w, http.StatusNotFound)
		return
	}

	key := vars["key"]
	value, err := bucket.Get(key)
	if nil != err {
		httpError(w, http.StatusInternalServerError)
		return
	}

	w.Write(value)
}

func (self *Server) HttpSetKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := lemur.Vars(r)

	bucketId := vars["bucketId"]
	bucket, err := self.db.Get(bucketId)
	if nil != err {
		httpError(w, http.StatusNotFound)
		return
	}

	err = lemur.Body(r, func(body []byte) error {
		key := vars["key"]
		return bucket.Set(key, body)
	})

	if nil != err {
		httpError(w, http.StatusInternalServerError)
		return
	}

}
