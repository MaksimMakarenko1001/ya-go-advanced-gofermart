package handler

import (
	"compress/gzip"
	"net/http"
	"slices"
	"strings"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg"
)

const (
	TypeContentTextPlain       = "text/plain"
	TypeContentApplicationJSON = "application/json"
)

type Middleware func(next http.Handler) http.Handler

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func MiddlewareTypeContentTextPlain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		if rq.Header.Get("Content-Type") != TypeContentTextPlain {
			WriteError(w, pkg.ErrBadRequest.SetInfo("not supported Content-Type"))
			return
		}

		next.ServeHTTP(w, rq)
	})
}

func MiddlewareTypeContentApplicationJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		if rq.Header.Get("Content-Type") != TypeContentApplicationJSON {
			WriteError(w, pkg.ErrBadRequest.SetInfo("not supported Content-Type"))
			return
		}

		next.ServeHTTP(w, rq)
	})
}

func MiddlewareCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		w := rw

		supportsGzip := slices.Contains(r.Header.Values("Accept-Encoding"), "gzip")
		if supportsGzip {
			cw := &compressWriter{
				w:  rw,
				zw: gzip.NewWriter(rw),
			}
			w = cw
			defer cw.Close()
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			zr, err := gzip.NewReader(r.Body)
			if err != nil {
				WriteError(w, err)
				return
			}

			cr := &compressReader{
				r:  r.Body,
				zr: zr,
			}

			r.Body = cr
			defer cr.Close()
		}
		next.ServeHTTP(w, r)
	})
}
