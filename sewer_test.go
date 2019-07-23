package sewer

import (
	"fmt"
	"testing"
	"time"
)

const DefaultAddr = "amqp://guest:guest@localhost:5672/"

func TestSewer_SubscribeAndPublish(t *testing.T) {

	sewer := New("test", DefaultAddr)

	p, err := sewer.Publisher()
	if err != nil {
		t.Fatalf("Error occurred creating publisher: %s", err.Error())
	}
	s, err := sewer.Subscriber()
	if err != nil {
		t.Fatalf("Error occurred creating subscriber: %s", err.Error())
	}
	go func() {
		for i := 0; i < 3; i++ {
			p <- []byte(fmt.Sprintf("The time is %s", time.Now().String()))
			time.Sleep(time.Second)
		}
		sewer.Close()
	}()
	for {
		select {
		case data := <-s:
			println(fmt.Sprintf("DATA: %s", data))

		case err := <-sewer.ErrorChannel():
			t.Fatalf("Error occurred: %s", err.Error())

		case <-sewer.stopChannel:
			return
		}
	}
}
