package param

type PostDiscussionRequest struct {
	Title   string    `json:"title" validate:"required"`
	Content string    `json:"content" validate:"required"`
}

type PostCommentRequest struct {
	Content    string    `json:"content" validate:"required"`
}