package logging

import (
	"net/http"

	"github.com/golang/glog"
)

type response struct {
	http.ResponseWriter
	r *http.Request
	s int
}

func (r *response) WriteHeader(status int) {
	r.s = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *response) log() {
	glog.Infof("addr=%s code=%d method=%s host=%s uri=%s",
		r.r.RemoteAddr,
		r.s,
		r.r.Method,
		r.r.Host,
		r.r.RequestURI)
}

// WithLog ...
func WithLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response{
			ResponseWriter: w,
			r:              r,
		}
		defer res.log()
		h.ServeHTTP(&res, r)
	})
}
