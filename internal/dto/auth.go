package dto

type UserRequestAdd struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=2,max=72"`
}

type UserAdd struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

type User struct {
	Id    uint   `json:"id"`
	Email string `json:"email"`
}

type StatusResponse struct {
	Status string `json:"status" example:"OK"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}
