package comment

import (
	"context"
	"fmt"

	"github.com/andreyxaxa/Comment-Tree/internal/dto"
	"github.com/andreyxaxa/Comment-Tree/internal/entity"
	"github.com/andreyxaxa/Comment-Tree/internal/repo"
)

type CommentUseCase struct {
	repo repo.CommentRepo
}

func New(r repo.CommentRepo) *CommentUseCase {
	return &CommentUseCase{
		repo: r,
	}
}

func (uc *CommentUseCase) CreateComment(ctx context.Context, parentID *int64, content string) (entity.Comment, error) {
	if parentID != nil {
		err := uc.repo.CommentExists(ctx, *parentID)
		if err != nil {
			return entity.Comment{}, fmt.Errorf("CommentUseCase - CreateComment - uc.repo.CommentExists: %w", err)
		}
	}

	c, err := uc.repo.CreateComment(ctx, parentID, content)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("CommentUseCase - CreateComment - uc.repo.CreateComment: %w", err)
	}

	return c, nil
}

func (uc *CommentUseCase) DeleteCommentWithChildren(ctx context.Context, id int64) error {
	err := uc.repo.DeleteCommentWithChildren(ctx, id)
	if err != nil {
		return fmt.Errorf("CommentUseCase - DeleteCommentWithChildren - uc.repo.DeleteCommentWithChildren: %w", err)
	}

	return nil
}

func (uc *CommentUseCase) GetComments(ctx context.Context, params dto.GetCommentsParams) (dto.PaginatedComments, error) {
	var comments []entity.Comment
	var total int
	var err error

	// 1. если указан поисковой запрос - полнотекстовый поиск
	if params.Search != "" {
		comments, total, err = uc.repo.SearchComments(ctx, params.Search, params.SortBy, params.Order, params.Limit, params.Offset)
		if err != nil {
			return dto.PaginatedComments{}, fmt.Errorf("CommentUseCase - GetComments - uc.repo.SearchComments: %w", err)
		}

		return dto.PaginatedComments{
			Comments: comments,
			Total:    total,
			Limit:    params.Limit,
			Offset:   params.Offset,
		}, nil
	}

	// 2. если указан конкретный родитель - получаем его дерево
	if params.ParentID != nil {
		comments, err = uc.repo.GetCommentWithChildren(ctx, *params.ParentID)
		if err != nil {
			return dto.PaginatedComments{}, fmt.Errorf("CommentUseCase - GetComments - uc.repo.GetCommentWithChildren: %w", err)
		}

		return dto.PaginatedComments{
			Comments: comments,
			Total:    len(comments),
			Limit:    params.Limit,
			Offset:   params.Offset,
		}, nil
	}

	// 3. иначе - получаем корневые комменты
	roots, total, err := uc.repo.GetRootComments(ctx, params.SortBy, params.Order, params.Limit, params.Offset)
	if err != nil {
		return dto.PaginatedComments{}, fmt.Errorf("CommentUseCase - GetComments - uc.repo.GetRootComments: %w", err)
	}

	// 3.1 получаем id корневых комментов
	rootIDs := make([]int64, len(roots))
	for i, r := range roots {
		rootIDs[i] = r.ID
	}

	// 3.2 получаем их деревья
	comments, err = uc.repo.GetTreesForRoots(ctx, rootIDs)
	if err != nil {
		return dto.PaginatedComments{}, fmt.Errorf("CommentUseCase - GetComments - uc.repo.GetTreesForRoots: %w", err)
	}

	return dto.PaginatedComments{
		Comments: comments,
		Total:    total,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}, nil
}
