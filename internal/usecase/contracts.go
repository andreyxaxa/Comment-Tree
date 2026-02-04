package usecase

import (
	"context"

	"github.com/andreyxaxa/Comment-Tree/internal/dto"
	"github.com/andreyxaxa/Comment-Tree/internal/entity"
)

type (
	CommentUseCase interface {
		CreateComment(ctx context.Context, parentID *int64, content string) (entity.Comment, error)
		DeleteCommentWithChildren(ctx context.Context, id int64) error
		GetComments(ctx context.Context, params dto.GetCommentsParams) (dto.PaginatedComments, error)
	}
)
