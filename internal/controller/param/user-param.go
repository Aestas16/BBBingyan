package param

type LoginUserRequest struct {
	Username    string  `json:"username"`
    Password    string  `json:"password"`
    Code        string  `json:"vercode"`
}