package entities

type User struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
}
