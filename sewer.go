package sewer

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

type Sewer struct {
	connectionString string
	errorChannel     chan error
	stopChannel      chan time.Time
	Name             string
}

func New(name, connectionString string) *Sewer {
	return &Sewer{
		Name:             name,
		connectionString: connectionString,
		errorChannel:     make(chan error),
		stopChannel:      make(chan time.Time),
	}
}

func (s *Sewer) newConnection() (*amqp.Connection, error) {
	c, err := amqp.Dial(s.connectionString)
	if err != nil {
		return nil, err
	}
	return c, err
}
func (s *Sewer) Subscriber() (chan []byte, error) {
	c, err := s.newConnection()
	if err != nil {
		return nil, err
	}

	rmqch, err := c.Channel()
	if err != nil {
		return nil, err
	}

	q, err := rmqch.QueueDeclare(
		fmt.Sprintf("sewer-%s", s.Name), // name
		true,                            // durable
		false,                           // delete when used
		false,                           // exclusive
		false,                           // no-wait
		nil,                             // arguments
	)
	if err != nil {
		return nil, err
	}

	// fair load balancing messages on queue when there are multiple receivers
	err = rmqch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}

	msgs, err := rmqch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}

	ch := make(chan []byte, 1024)
	go func() {
		defer func() {
			s.errorChannel <- rmqch.Close()
		}()
		for {
			select {
			case d := <-msgs:
				ch <- d.Body
				err := d.Ack(false)
				if err != nil {
					s.errorChannel <- err
				}
			case <-s.stopChannel:
				err := c.Close()
				if err != nil {
					s.errorChannel <- err
				}
				close(ch)
				return
			}
		}
	}()
	return ch, nil
}
func (s *Sewer) Publisher() (chan []byte, error) {
	c, err := s.newConnection()
	if err != nil {
		return nil, err
	}

	rmqch, err := c.Channel()
	if err != nil {
		return nil, err
	}

	q, err := rmqch.QueueDeclare(
		fmt.Sprintf("sewer-%s", s.Name), // name
		true,                            // durable
		false,                           // delete when unused
		false,                           // exclusive
		false,                           // no-wait
		nil,                             // arguments
	)
	if err != nil {
		return nil, err
	}

	ch := make(chan []byte)
	go func() {
		defer func() {
			err := rmqch.Close()
			if err != nil {
				s.errorChannel <- err
			}
		}()
		for {
			select {
			case msg := <-ch:
				err = rmqch.Publish(
					"",     // exchange
					q.Name, // routing key
					false,  // mandatory
					false,  // immediate
					amqp.Publishing{
						DeliveryMode: amqp.Persistent,
						ContentType:  "text/plain",
						Body:         msg,
					})
				if err != nil {
					s.errorChannel <- err
				}
			case <-s.stopChannel:
				err := c.Close()
				if err != nil {
					s.errorChannel <- err
				}
				close(ch)
				return
			}
		}
	}()
	return ch, nil
}
func (s *Sewer) ErrorChannel() <-chan error {
	return s.errorChannel
}

func (s *Sewer) Close() {
	close(s.stopChannel)
}
