package responses

type Token struct {
	Token string `json:"token"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
