package controller

type userInfoDTO struct {
	ID string `json:"id"`
}

type createTokenDTO struct {
	ID          string `json:"id"`
	RememberHim bool   `json:"rememberHim"`
}

type signInDTO struct {
	Gmail       string `json:"gmail"`
	Password    string `json:"password"`
	RememberHim bool   `json:"rememberHim"`
}

type signUpDTO struct {
	Gmail               string `json:"gmail"`
	Password            string `json:"password"`
	FullName            string `json:"fullName"`
	Phone               string `json:"phone"`
	Key                 string `json:"key"`
	AllowsAdvertisement bool   `json:"allowsAdvertisement"`
}

type passwordResetDTO struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
	Key      string `json:"key"`
}
