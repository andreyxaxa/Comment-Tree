package restapi

import (
	v1 "github.com/andreyxaxa/Comment-Tree/internal/controller/restapi/v1"
	"github.com/andreyxaxa/Comment-Tree/internal/usecase"
	"github.com/andreyxaxa/Comment-Tree/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, c usecase.CommentUseCase, l logger.Interface) {
	apiV1Group := app.Group("/v1")
	{
		v1.NewCommentRoutes(apiV1Group, c, l)
	}
}
