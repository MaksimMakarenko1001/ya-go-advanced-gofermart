package backoff

import "context"

type ErrorClassification int
type retried func(ctx context.Context) (err error)

const (
	// NonRetriable - операцию не следует повторять
	NonRetriable ErrorClassification = iota

	// Retriable - операцию можно повторить
	Retriable
)
