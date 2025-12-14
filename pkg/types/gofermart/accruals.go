package gofermart

const (
	AccrualStatusNone       AccrualStatusType = "NONE"
	AccrualStatusRegistered AccrualStatusType = "REGISTERED"
	AccrualStatusInvalid    AccrualStatusType = "INVALID"
	AccrualStatusProcessing AccrualStatusType = "PROCESSING"
	AccrualStatusProcessed  AccrualStatusType = "PROCESSED"
)

type AccrualStatusType string

func (ast AccrualStatusType) String() string {
	return string(ast)
}
