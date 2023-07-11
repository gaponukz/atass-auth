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
