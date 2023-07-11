package entities

type User struct {
	Gmail               string   `json:"gmail"`
	Password            string   `json:"password"`
	Phone               string   `json:"phone"`
	FullName            string   `json:"fullName"`
	AllowsAdvertisement bool     `json:"allowsAdvertisement"`
	PurchasedRouteIds   []string `json:"purchasedRouteIds"`
}

type UserEntity struct {
	ID string
	User
}
