package response

import "time"

type CreateCommentResponse struct {
	ID        int64     `json:"id" example:"12"`
	ParentID  *int64    `json:"parent_id" example:"1"`
	Content   string    `json:"content" example:"nice picture!!!"`
	CreatedAt time.Time `json:"created_at" example:"2026-02-02T14:31:00Z"`
}
