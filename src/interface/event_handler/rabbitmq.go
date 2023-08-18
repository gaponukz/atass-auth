package event_handler

import (
	"auth/src/domain/entities"
	"auth/src/domain/events"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type addRouteService interface {
	AddRoute(userID string, path entities.Path) error
}

type deleteRouteService interface {
	DeleteRoute(string, entities.Path) error
}

type routesEventsListener struct {
	addRouteService    addRouteService
	deleteRouteService deleteRouteService
	conn               *amqp.Connection
	ch                 *amqp.Channel
	msgs               <-chan amqp.Delivery
}

func NewRoutesEventsListener(as addRouteService, ds deleteRouteService, rabbitUrl string) (*routesEventsListener, error) {
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

	return &routesEventsListener{addRouteService: as, deleteRouteService: ds, conn: conn, ch: ch, msgs: msgs}, nil
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

			switch event.EventType {
			case "booked":
				_ = r.addRouteService.AddRoute(
					event.PassengerID,
					entities.Path{
						RootRouteID: event.RouteID,
						MoveFromID:  event.MoveFromID,
						MoveToID:    event.MoveToID,
					},
				)

			case "removed":
				_ = r.deleteRouteService.DeleteRoute(
					event.PassengerID,
					entities.Path{
						RootRouteID: event.RouteID,
						MoveFromID:  event.MoveFromID,
						MoveToID:    event.MoveToID,
					},
				)
			}
		}
	}()
	<-forever
}

func (r routesEventsListener) Close() {
	_ = r.ch.Close()
	_ = r.conn.Close()
}
