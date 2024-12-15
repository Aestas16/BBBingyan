package param

type SendVerCodeRequest struct {
	Username    string    `json:"username" validate:"required"`
}