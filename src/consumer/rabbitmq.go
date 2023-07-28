package main

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type routesService interface {
	AddRoute(userID, routeID string) error
}

type routesEventsListener struct {
	service routesService
	conn    *amqp.Connection
	ch      *amqp.Channel
	msgs    <-chan amqp.Delivery
}

func NewRoutesEventsListener(s routesService, rabbitUrl string) (*routesEventsListener, error) {
	conn, err := amqp.Dial(rabbitUrl)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	err = ch.QueueBind(
		"payments",
		"",
		"payments_exchange",
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}
	msgs, err := ch.Consume(
		"payments", // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}

	return &routesEventsListener{service: s, conn: conn, ch: ch, msgs: msgs}, nil
}

func (r routesEventsListener) Listen() {
	forever := make(chan bool)

	go func() {
		for d := range r.msgs {
			var data map[string]interface{}

			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				fmt.Printf("Error parsing JSON: %v\n", err)
				return
			}

			routeID := data["routeId"].(string)
			passengerID := data["passanger"].(map[string]interface{})["id"].(string)

			r.service.AddRoute(passengerID, routeID)
		}
	}()

	fmt.Println("Consumer is waiting for messages. To exit press CTRL+C")
	<-forever
}

func (r routesEventsListener) Close() {
	_ = r.ch.Close()
	_ = r.conn.Close()
}

type testRoutesService struct{}

func (t testRoutesService) AddRoute(userID, routeID string) error {
	fmt.Println(userID, routeID)
	return nil
}

func main() {
	server, err := NewRoutesEventsListener(testRoutesService{}, "amqp://user:password@localhost:5672/")
	if err != nil {
		panic(err.Error())
	}

	defer server.Close()

	server.Listen()
}
