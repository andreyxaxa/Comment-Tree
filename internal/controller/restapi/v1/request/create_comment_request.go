package request

type CreateCommentRequest struct {
	ParentID *int64 `json:"parent_id"`
	Content  string `json:"content"`
}
