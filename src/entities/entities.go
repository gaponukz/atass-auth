package entities

type User struct {
	Gmail             string   `json:"gmail"`
	Password          string   `json:"password"`
	Phone             string   `json:"phone"`
	FullName          string   `json:"fullName"`
	RememberHim       bool     `json:"rememberHim"`
	PurchasedRouteIds []string `json:"purchasedRouteIds"`
}

type FutureUser struct {
	UniqueKey string `json:"uniqueKey"`
	User
}
