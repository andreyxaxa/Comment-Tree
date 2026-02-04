package dto

import "github.com/andreyxaxa/Comment-Tree/internal/entity"

type PaginatedComments struct {
	Comments []entity.Comment
	Total    int
	Limit    int
	Offset   int
}
