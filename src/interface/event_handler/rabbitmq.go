package event_handler

import (
	"auth/src/domain/entities"
	"auth/src/domain/events"
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
		"events",
		"",
		"events_exchange",
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}
	msgs, err := ch.Consume(
		"events", // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
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
			var event events.BookingEvent

			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				fmt.Printf("Error parsing JSON: %v\n", err)
			}

			_ = r.service.AddRoute(
				event.PassengerID,
				entities.Path{
					RootRouteID: event.RouteID,
					MoveFromID:  event.MoveFromID,
					MoveToID:    event.MoveToID,
				},
			)
		}
	}()
	<-forever
}

func (r routesEventsListener) Close() {
	_ = r.ch.Close()
	_ = r.conn.Close()
}
