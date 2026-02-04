package response

type PaginatedCommentsResponse struct {
	Comments []*CommentTreeResponse `json:"comments"`
	Total    int                    `json:"total"`
	Limit    int                    `json:"limit"`
	Offset   int                    `json:"offset"`
	Page     int                    `json:"page"`
	Pages    int                    `json:"pages"`
}
