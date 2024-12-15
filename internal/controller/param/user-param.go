package param

type LoginUserRequest struct {
	Username    string  `json:"username" validate:"required"`
    Password    string  `json:"password" validate:"required"`
    Code        string  `json:"vercode" validate:"required"`
}

type UserRequest struct {
    Username    string  `json:"username"`
    Password    string  `json:"password"`
    Email       string  `json:"email"`
}