package dto

type GetCommentsParams struct {
	ParentID *int64
	Search   string
	SortBy   string
	Order    string
	Limit    int
	Offset   int
}
