package controller

type credentials struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
}

type userInfoDTO struct {
	Gmail       string `json:"gmail"`
	RememberHim bool   `json:"rememberHim"`
}

type signInDTO struct {
	Gmail       string `json:"gmail"`
	Password    string `json:"password"`
	RememberHim bool   `json:"rememberHim"`
}

type singUpDTO struct {
	Gmail               string `json:"gmail"`
	Password            string `json:"password"`
	FullName            string `json:"fullName"`
	Phone               string `json:"phone"`
	Key                 string `json:"key"`
	AllowsAdvertisement bool   `json:"allowsAdvertisement"`
	RememberHim         bool   `json:"rememberHim"`
}

type passwordResetDTO struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
	Key      string `json:"key"`
}
