package controller

type userInfoDTO struct {
	ID                  string `json:"id"`
	Gmail               string `json:"gmail"`
	FullName            string `json:"fullName"`
	Phone               string `json:"phone"`
	AllowsAdvertisement bool   `json:"allowsAdvertisement"`
}

type createTokenDTO struct {
	RememberHim bool `json:"rememberHim"`
	userInfoDTO
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

type updateUserDTO struct {
	FullName            string `json:"fullName"`
	Phone               string `json:"phone"`
	AllowsAdvertisement bool   `json:"allowsAdvertisement"`
}
