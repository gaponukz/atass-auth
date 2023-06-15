package entities

type User struct {
	Gmail               string   `json:"gmail"`
	Password            string   `json:"password"`
	Phone               string   `json:"phone"`
	FullName            string   `json:"fullName"`
	RememberHim         bool     `json:"rememberHim"`
	AllowsAdvertisement bool     `json:"allowsAdvertisement"`
	PurchasedRouteIds   []string `json:"purchasedRouteIds"`
}

type GmailWithKeyPair struct {
	Gmail string `json:"gmail"`
	Key   string `json:"id"`
}
