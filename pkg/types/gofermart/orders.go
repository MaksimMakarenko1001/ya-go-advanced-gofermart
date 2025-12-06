package gofermart

const (
	OrderStatusNew        OrderStatusType = "NEW"
	OrderStatusInvalid    OrderStatusType = "INVALID"
	OrderStatusProcessing OrderStatusType = "PROCESSING"
	OrderStatusProcessed  OrderStatusType = "PROCESSED"
)

type OrderStatusType string

func (ost OrderStatusType) String() string {
	return string(ost)
}
