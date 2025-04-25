package main

import (
	"UserValidator/interfaces"
	"UserValidator/plugins"
	str "UserValidator/structs"
	"context"
	"fmt"
	"time"
)

const RANDOMSTRING = "1234"

type Identity[T plugins.FaceIDData | str.SMSData] struct {
	Input  T
	Plugin interfaces.Verifier
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	faceid := Identity[plugins.FaceIDData]{
		Input: plugins.FaceIDData{
			BinaryData: []byte(RANDOMSTRING),
		},
		Plugin: &plugins.FaceIDData{},
	}
	faceid.Run(ctx)
	defer cancel()
}

func (i Identity[T]) Run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Unexpected error", r)
		}
	}()
	ok, err := i.Plugin.Verify(ctx, i.Input)
	if ok {
		switch i.Plugin.(type) {
		case interfaces.Scorable:
			score := i.Plugin.(interfaces.Scorable).GetScore()
			fmt.Printf("Estadisticas encontradas %v", score)
		default:
			fmt.Println("Plugin sin estadisticas")
		}
	} else {
		panic(err)
	}
}
