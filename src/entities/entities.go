package entities

type User struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
}

type FutureUser struct {
	UniqueKey string `json:"uniqueKey"`
	User
}
