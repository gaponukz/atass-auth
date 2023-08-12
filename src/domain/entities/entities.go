package entities

type Path struct {
	RootRouteID string `json:"rootRouteId"`
	MoveFromID  string `json:"movingFromId"`
	MoveToID    string `json:"movingTowardsId"`
}

type User struct {
	ID                  string `json:"id"`
	Gmail               string `json:"gmail"`
	Password            string `json:"password"`
	Phone               string `json:"phone"`
	FullName            string `json:"fullName"`
	AllowsAdvertisement bool   `json:"allowsAdvertisement"`
	PurchasedRouteIds   []Path `json:"purchasedRouteIds"`
}
