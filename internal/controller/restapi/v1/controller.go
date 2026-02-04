package v1

import (
	"github.com/andreyxaxa/Comment-Tree/internal/usecase"
	"github.com/andreyxaxa/Comment-Tree/pkg/logger"
)

type V1 struct {
	c usecase.CommentUseCase
	l logger.Interface
}
