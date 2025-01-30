package incoming

import (
	domain "ecommerce/Order/Domain"
	"ecommerce/Order/Domain/events/contracts"
	"time"
)

// this is our domain type
type PaymentCompleted struct {
	OrderID string
	Amount  domain.Money
	Time    time.Time
}

// maps the shared contract to our bounded context's content
func PaymentCompletedFromContract(c contracts.PaymentCompletedV1) PaymentCompleted {
	return PaymentCompleted{
		OrderID: c.OrderID,
		Amount:  domain.Money(c.Amount),
		Time:    c.TimeStamp,
	}
}
