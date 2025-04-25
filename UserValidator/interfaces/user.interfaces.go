package interfaces

import "context"

type Verifier interface {
	Verify(context.Context, any) (bool, error)
}

type Scorable interface {
	Verifier
	GetScore() float32
	GetDescription() string
}
