package main

import (
	inf "PayFlex/interfaces"
	pl "PayFlex/plugins"
	"context"
	"fmt"
	"time"
)

type PaymentProccessor[T pl.CardPayment] struct {
	PaymentMethod inf.PaymentMethod
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cardPayment := PaymentProccessor[pl.CardPayment]{
		PaymentMethod: &pl.CardPayment{
			Total:      59.00,
			CardNumber: "1234567890123456",
		},
	}
	go ExecutePayment(ctx, cardPayment)
	func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()
}

func ExecutePayment[T pl.CardPayment](ctx context.Context, method PaymentProccessor[T]) {
	func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from error", r)
		}
	}()
	err := method.PaymentMethod.Payment(ctx)
	if err != nil {
		panic(err)
	}

	if _, ok := method.PaymentMethod.(inf.Refundable); ok {
		fmt.Println("Refund policies")
	}
}
