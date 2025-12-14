package accrual

import (
	"errors"
	"net/url"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/backoff"
)

var (
	errAccuralInternalServer  = errors.New("accural internal server error")
	errAccuralTooManyRequests = errors.New("accural too many requests")
)

func ClassifyHTTPError(err error) backoff.ErrorClassification {
	if err == nil {
		return backoff.NonRetriable
	}

	var urlError *url.Error
	if errors.As(err, &urlError) {
		return backoff.Retriable
	}

	if errors.Is(err, errAccuralInternalServer) {
		return backoff.Retriable
	}

	return backoff.NonRetriable
}
