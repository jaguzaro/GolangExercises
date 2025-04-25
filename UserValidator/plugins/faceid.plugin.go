package plugins

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type FaceIDData struct {
	BinaryData []byte
}

func (f *FaceIDData) Verify(ctx context.Context, data any) (bool, error) {
	FaceID(ctx)
	switch data.(type) {
	case FaceIDData:
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("FaceID completado")
			fmt.Println(data.(FaceIDData).BinaryData)
			return true, nil
		case <-ctx.Done():
			err := ctx.Err()
			return false, err
		}
	}

	return false, errors.New("non type")
}

func FaceID(ctx context.Context) {
	fmt.Println("FaceID Running...")
}

func (f *FaceIDData) GetScore() float32 {
	return 0.9
}

func (f *FaceIDData) GetDescription() string {
	return "HolaMundo"
}
