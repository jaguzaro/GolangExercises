package plugins

import (
	inf "PayFlex/interfaces"
	"context"
	"fmt"
	"time"
)

type CardPayment struct {
	Total      float32
	CardNumber string
}

const MINPAYMENT = 0.00
const MAXPAYMENT = 100.00

func (c *CardPayment) Payment(ctx context.Context) error {
	if c.Total > MINPAYMENT || c.Total <= MAXPAYMENT {
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("Simulating card payment..")
			fmt.Println(c)
			return nil
		case <-ctx.Done():
			err := ctx.Err()
			return fmt.Errorf("payment canceled %v", err)
		}
	} else {
		return fmt.Errorf("total is not in range")
	}
}

func (c *CardPayment) Reimbursment(ctx context.Context, data any) error {
	if _, ok := data.(inf.Refundable); !ok {
		return fmt.Errorf("data is not type Refundable")
	}

	paymentData := data.(CardPayment)
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Simulating card reimbursment..")
		fmt.Println(paymentData)
		return nil
	case <-ctx.Done():
		err := ctx.Err()
		return fmt.Errorf("reimbursment canceled %v", err)

	}
}
