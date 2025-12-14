package handler

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

type responseHashWriter struct {
	http.ResponseWriter
	body     bytes.Buffer
	hashFunc func(message []byte) (string, error)
}

func (rh *responseHashWriter) Write(b []byte) (int, error) {
	rh.ResponseWriter.Write(b)
	return rh.body.Write(b)
}

func (rh *responseHashWriter) WriteHeader(statusCode int) {
	rh.ResponseWriter.WriteHeader(statusCode)
	if statusCode == http.StatusOK {
		hash, err := rh.hashFunc(rh.body.Bytes())
		if err != nil {
			WriteError(rh.ResponseWriter, err)
			return
		}
		rh.ResponseWriter.Header().Set("HashSHA256", hash)
	}
}

type ResponseInfo struct {
	Size   int
	Status int
	Body   bytes.Buffer
}

type responseWriter struct {
	http.ResponseWriter
	response *ResponseInfo
}

func (r *responseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.response.Size += size
	r.response.Body.Write(b)
	return size, err
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.response.Status = statusCode
}

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(b []byte) (int, error) {
	return c.zw.Write(b)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusOK {
		c.w.Header().Set("Content-Encoding", "gzip")
	}

	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (c *compressReader) Read(b []byte) (int, error) {
	return c.zr.Read(b)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
