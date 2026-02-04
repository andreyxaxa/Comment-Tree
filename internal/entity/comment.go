package entity

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int64         `json:"id"`
	ParentID  sql.NullInt64 `json:"parent_id"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"created_at"`

	Depth int     `json:"depth"`
	Path  []int64 `json:"path,omitempty"`
}
