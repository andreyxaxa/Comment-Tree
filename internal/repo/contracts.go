package repo

import (
	"context"

	"github.com/andreyxaxa/1/internal/entity"
)

type (
	CommentRepo interface {
		CreateComment(ctx context.Context, parentID *int64, content string) (entity.Comment, error)
		CommentExists(ctx context.Context, id int64) error
		GetCommentWithChildren(ctx context.Context, id int64) ([]entity.Comment, error)
		DeleteCommentWithChildren(ctx context.Context, id int64) error
		SearchComments(ctx context.Context, search string, sortBy, order string, limit, offset int) ([]entity.Comment, int, error)
		GetRootComments(ctx context.Context, sortBy, order string, limit, offset int) ([]entity.Comment, int, error)
		GetTreesForRoots(ctx context.Context, rootIDs []int64) ([]entity.Comment, error)
	}
)
