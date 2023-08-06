package event_handler

import (
	"auth/src/domain/entities"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type routesService interface {
	AddRoute(userID string, path entities.Path) error
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
		"passenger_payments",
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
		"passenger_payments", // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
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
			}

			routeID, ok := data["routeId"].(string)
			if !ok {
				fmt.Printf("Error parsing JSON: no 'routeId' field\n")
			}
			passenger, ok := data["passenger"].(map[string]interface{})
			if !ok {
				fmt.Printf("Error parsing JSON: no 'passenger' field\n")
			}
			passengerID, ok := passenger["id"].(string)
			if !ok {
				fmt.Printf("Error parsing JSON: no 'passenger.id' field\n")
			}
			moveFrom, ok := passenger["movingFromId"].(string)
			if !ok {
				fmt.Printf("Error parsing JSON: no 'passenger.movingFromId' field\n")
			}
			moveTo, ok := passenger["movingTowardsId"].(string)
			if !ok {
				fmt.Printf("Error parsing JSON: no 'passenger.movingTowardsId' field\n")
			}

			_ = r.service.AddRoute(passengerID, entities.Path{RootRouteID: routeID, MoveFromID: moveFrom, MoveToID: moveTo})
		}
	}()
	<-forever
}

func (r routesEventsListener) Close() {
	_ = r.ch.Close()
	_ = r.conn.Close()
}
