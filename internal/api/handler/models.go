package handler

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

type ResponseInfo struct {
	Size   int
	Status int
	Body   bytes.Buffer
}

type ResponseWriter struct {
	http.ResponseWriter
	Response *ResponseInfo
}

func (r *ResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.Response.Size += size
	r.Response.Body.Write(b)
	return size, err
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.Response.Status = statusCode
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
