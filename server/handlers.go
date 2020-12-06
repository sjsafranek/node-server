package server

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/sjsafranek/node-server/httprouter"
)

func httpError(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	var resp Response
	resp.SetError(errors.New(http.StatusText(http.StatusNotFound)))
	resp.Write(w)
}

func (self *Server) HttpGetKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := httprouter.Vars(r)

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
	vars := httprouter.Vars(r)

	bucketId := vars["bucketId"]
	bucket, err := self.db.Get(bucketId)
	if nil != err {
		httpError(w, http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if nil != err {
		httpError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	key := vars["key"]
	err = bucket.Set(key, body)
	if nil != err {
		httpError(w, http.StatusInternalServerError)
		return
	}

}
