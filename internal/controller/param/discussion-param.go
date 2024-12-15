package param

type PostDiscussionRequest struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
}

type PostCommentRequest struct {
	Content    string    `json:"content"`
}