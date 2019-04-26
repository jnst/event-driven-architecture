package pubsub

import "fmt"

type Usecase struct {}

func NewUsecase() *Usecase {
	return &Usecase{}
}

func (u *Usecase) DoAnyProcess(eventName string) {
	fmt.Printf("  event: %s\n", eventName)

	// Write your business domain code here
}
