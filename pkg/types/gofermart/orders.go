package gofermart

const (
	OrderStatusNone       OrderStatusType = "NONE"
	OrderStatusNew        OrderStatusType = "NEW"
	OrderStatusInvalid    OrderStatusType = "INVALID"
	OrderStatusProcessing OrderStatusType = "PROCESSING"
	OrderStatusProcessed  OrderStatusType = "PROCESSED"
)

type OrderStatusType string

func (ost OrderStatusType) String() string {
	return string(ost)
}
