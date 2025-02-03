package incoming

import (
	"ecommerce/Order/Domain/IntegrationEvents/contracts"
	"ecommerce/Order/Domain/models"
	"time"
)

// this is our domain type
type PaymentCompleted struct {
	OrderID string
	Amount  models.Money
	Time    time.Time
}

func (p PaymentCompleted) Name() string {
	return "PaymentCompleted"
}

// maps the shared contract to our bounded context's content
func PaymentCompletedFromContract(c contracts.PaymentCompletedV1) PaymentCompleted {
	return PaymentCompleted{
		OrderID: c.OrderID,
		Amount:  models.Money(c.Amount),
		Time:    c.TimeStamp,
	}
}
