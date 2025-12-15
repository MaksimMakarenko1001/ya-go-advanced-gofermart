package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg"
)

func WriteJSONResult(w http.ResponseWriter, response []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(response); err != nil {
		WriteError(w, fmt.Errorf("write json not ok, %w", err))
	}
}

func WriteResult(w http.ResponseWriter, res string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, res)
}

func WriteOK(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
}

func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain")
	if err == nil {
		err = pkg.ErrInternalServer
	}

	var errE *pkg.Error

	if !errors.As(err, &errE) {
		errE = pkg.ErrInternalServer
	}
	http.Error(w, errE.Error(), errE.HTTPStatus())
}
