package incoming

import (
	"ecommerce/Order/Domain/models"
	"time"
)

// this is our domain type
// application service is responsible to map from application to domain layer
type PaymentCompleted struct {
	OrderID string
	Amount  models.Money
	Time    time.Time
}

func (p PaymentCompleted) Name() string {
	return "PaymentCompleted"
}
