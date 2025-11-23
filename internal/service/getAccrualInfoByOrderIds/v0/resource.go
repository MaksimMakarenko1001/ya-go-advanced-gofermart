package v0

import "strings"

type AccrualResponse struct {
	Err     error
	Payload *AccrualPayload
}

type AccrualPayload struct {
	Order   string            `json:"order"`
	Status  AccrualStatusType `json:"status"`
	Accrual float64           `json:"accrual"`
}

const (
	AccrualStatusRegistered AccrualStatusType = "REGISTERED"
	AccrualStatusInvalid    AccrualStatusType = "INVALID"
	AccrualStatusProcessing AccrualStatusType = "PROCESSING"
	AccrualStatusProcessed  AccrualStatusType = "PROCESSED"
)

type AccrualStatusType string

func (ast AccrualStatusType) String() string {
	return strings.ToLower(string(ast))
}
