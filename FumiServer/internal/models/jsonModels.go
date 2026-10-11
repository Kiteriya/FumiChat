package models

type AuthJson struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	JwtTOKEN string `json:"token"`
}
