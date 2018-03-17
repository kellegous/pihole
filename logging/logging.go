package logging

import (
	"net/http"

	"go.uber.org/zap"
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
	zap.L().Info("web response",
		zap.String("addr", r.r.RemoteAddr),
		zap.Int("code", r.s),
		zap.String("method", r.r.Method),
		zap.String("host", r.r.Host),
		zap.String("uri", r.r.RequestURI))
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

// MustSetup ...
func MustSetup() {
	if err := Setup(); err != nil {
		panic(err)
	}
}

// Setup ...
func Setup() error {
	l, err := zap.NewProduction()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(l)
	return nil
}
