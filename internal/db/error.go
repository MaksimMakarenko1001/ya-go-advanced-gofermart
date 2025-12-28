package db

import (
	"errors"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/backoff"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	errScanRow   = errors.New("error scan row")
	errNoData    = errors.New("error no data")
	errUnmarshal = errors.New("error unmarshal")

	errHostUndefined   = errors.New("host undefined")
	errPortUndefined   = errors.New("port undefined")
	errDBNameUndefined = errors.New("db name undefined")
	errUserUndefined   = errors.New("user undefined")
)

func ClassifyPgError(err error) backoff.ErrorClassification {
	if err == nil {
		return backoff.NonRetriable
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return backoff.NonRetriable
	}

	switch pgErr.Code {
	case pgerrcode.ConnectionException,
		pgerrcode.ConnectionDoesNotExist,
		pgerrcode.ConnectionFailure,
		pgerrcode.CannotConnectNow:
		return backoff.Retriable
	}

	return backoff.NonRetriable
}
