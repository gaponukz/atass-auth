package events

type BookingEvent struct {
	EventType   string `json:"type"`
	RouteID     string `json:"route_id"`
	MoveFromID  string `json:"move_from_id"`
	MoveToID    string `json:"move_to_id"`
	PassengerID string `json:"passenger_id"`
}
