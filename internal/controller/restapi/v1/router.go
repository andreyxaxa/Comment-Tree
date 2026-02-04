package v1

import (
	"github.com/andreyxaxa/Comment-Tree/internal/usecase"
	"github.com/andreyxaxa/Comment-Tree/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func NewCommentRoutes(apiV1Group fiber.Router, c usecase.CommentUseCase, l logger.Interface) {
	r := &V1{c: c, l: l}

	commentsGroup := apiV1Group.Group("/comments")

	{
		// API
		commentsGroup.Post("/", r.create)
		commentsGroup.Get("/", r.getComments)
		commentsGroup.Delete("/:id", r.deleteCommentTree)

		// UI
		apiV1Group.Get("/ui", r.showUI)
	}
}
