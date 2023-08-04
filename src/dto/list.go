package dto

type GmailWithKeyPairDTO struct {
	Gmail string `json:"gmail"`
	Key   string `json:"id"`
}

type SignInDTO struct {
	Gmail       string `json:"gmail"`
	Password    string `json:"password"`
	RememberHim bool   `json:"rememberHim"`
}

type SignUpDTO struct {
	Gmail               string `json:"gmail"`
	Password            string `json:"password"`
	FullName            string `json:"fullName"`
	Phone               string `json:"phone"`
	Key                 string `json:"key"`
	AllowsAdvertisement bool   `json:"allowsAdvertisement"`
}

type PasswordResetDTO struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
	Key      string `json:"key"`
}

type UpdateUserDTO struct {
	FullName            string `json:"fullName"`
	Phone               string `json:"phone"`
	AllowsAdvertisement bool   `json:"allowsAdvertisement"`
}

type UserInfoDTO struct {
	ID                  string `json:"id"`
	Gmail               string `json:"gmail"`
	FullName            string `json:"fullName"`
	Phone               string `json:"phone"`
	AllowsAdvertisement bool   `json:"allowsAdvertisement"`
}

type CreateTokenDTO struct {
	RememberHim bool `json:"rememberHim"`
	UserInfoDTO
}

type UpdateTokenDTO struct {
	FullName            string `json:"fullName"`
	Phone               string `json:"phone"`
	AllowsAdvertisement bool   `json:"allowsAdvertisement"`
}
