package hash

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type Repository struct {
	key string
}

func New(cfg Config) *Repository {
	return &Repository{
		key: cfg.Key,
	}
}

func (r *Repository) Hash(ctx context.Context, message []byte) (string, error) {
	if r.key == "" {
		return "", nil
	}

	h := hmac.New(sha256.New, []byte(r.key))
	if _, err := h.Write(message); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
