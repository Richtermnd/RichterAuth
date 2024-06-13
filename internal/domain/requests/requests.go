package requests

type Register struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeatPassword"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPassword struct {
	OldPassword    string `json:"oldPassword"`
	NewPassword    string `json:"newPassword"`
	RepeatPassword string `json:"repeatPassword"`
}
