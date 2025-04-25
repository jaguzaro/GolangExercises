package interfaces

import "context"

type PaymentMethod interface {
	Payment(context.Context) error
}

type Refundable interface {
	PaymentMethod
	Reimbursment(context.Context, any) error
}
