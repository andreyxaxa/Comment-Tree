package restapi

import (
	_ "github.com/andreyxaxa/Comment-Tree/docs" // Swagger docs.
	v1 "github.com/andreyxaxa/Comment-Tree/internal/controller/restapi/v1"
	"github.com/andreyxaxa/Comment-Tree/internal/usecase"
	"github.com/andreyxaxa/Comment-Tree/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

// @title Comment tree
// @version 1.0.0
// @host localhost:8080
// @BasePath /v1
func NewRouter(app *fiber.App, c usecase.CommentUseCase, l logger.Interface) {
	apiV1Group := app.Group("/v1")
	{
		v1.NewCommentRoutes(apiV1Group, c, l)
	}
}
