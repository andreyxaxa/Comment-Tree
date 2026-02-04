package response

import "time"

type CommentTreeResponse struct {
	ID        int64                  `json:"id"`
	ParentID  *int64                 `json:"parent_id"`
	Content   string                 `json:"content"`
	CreatedAt time.Time              `json:"created_at"`
	Depth     int                    `json:"depth"`
	Children  []*CommentTreeResponse `json:"children,omitempty"`
}
