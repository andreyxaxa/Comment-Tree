package request

import "strings"

type GetCommentsReqeust struct {
	ParentID *int64 `query:"parent_id"`
	Search   string `query:"search"`
	SortBy   string `query:"sort_by"`
	Order    string `query:"order"`
	Limit    int    `query:"limit"`
	Offset   int    `query:"offset"`
}

func (r *GetCommentsReqeust) Validate() {
	if r.Limit <= 0 || r.Limit > 100 {
		r.Limit = 20
	}

	if r.Offset < 0 {
		r.Offset = 0
	}

	if r.SortBy != "created_at" && r.SortBy != "id" {
		r.SortBy = "created_at"
	}

	r.Order = strings.ToUpper(r.Order)

	if r.Order != "DESC" && r.Order != "ASC" {
		r.Order = "DESC"
	}
}
