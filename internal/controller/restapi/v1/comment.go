package v1

import (
	"errors"
	"net/http"

	"github.com/andreyxaxa/Comment-Tree/internal/controller/restapi/v1/request"
	"github.com/andreyxaxa/Comment-Tree/internal/controller/restapi/v1/response"
	"github.com/andreyxaxa/Comment-Tree/internal/controller/restapi/v1/utils"
	"github.com/andreyxaxa/Comment-Tree/internal/dto"
	"github.com/andreyxaxa/Comment-Tree/internal/entity"
	"github.com/andreyxaxa/Comment-Tree/pkg/types/errs"
	"github.com/gofiber/fiber/v2"
)

// @Summary Create new comment
// @Description Creates new comment
// @Tags comments
// @Accept json
// @Produce json
// @Param request body request.CreateCommentRequest true "Comment"
// @Success 201 {object} response.CreateCommentResponse
// @Failure 400 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /v1/comments [post]
func (r *V1) create(ctx *fiber.Ctx) error {
	var body request.CreateCommentRequest

	err := ctx.BodyParser(&body)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	comment, err := r.c.CreateComment(ctx.UserContext(), body.ParentID, body.Content)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "parent not found")
		}
		r.l.Error(err, "restapi - v1 - create")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	resp := response.CreateCommentResponse{
		ID:        comment.ID,
		ParentID:  utils.NullInt64ToPtr(comment.ParentID),
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}

	return ctx.Status(http.StatusCreated).JSON(resp)
}

// @Summary Get comments
// @Description Get comment(s) with all replies, search, sort
// @Tags comments
// @Produce json
// @Param parent_id query string false "Parent ID"
// @Param search query string false "Search text"
// @Param sort_by query string false "Sort option" Enums(created_at, id)
// @Param order query string false "Sort order" Enums(asc, ASC, desc, DESC)
// @Param limit query string false "Limit of comments on one page, default 20"
// @Param offset query string false "Offset for displaying a specific page, default 0"
// @Success 200 {object} response.PaginatedCommentsResponse
// @Failure 400 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /v1/comments [get]
func (r *V1) getComments(ctx *fiber.Ctx) error {
	var req request.GetCommentsReqeust

	err := ctx.QueryParser(&req)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid query parameters")
	}

	req.Validate()

	result, err := r.c.GetComments(ctx.UserContext(), dto.GetCommentsParams{
		ParentID: req.ParentID,
		Search:   req.Search,
		SortBy:   req.SortBy,
		Order:    req.Order,
		Limit:    req.Limit,
		Offset:   req.Offset,
	})
	if err != nil {
		r.l.Error(err, "restapi - v1 - getComments")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	var trees []*response.CommentTreeResponse

	// собираем все комменты для каждого корня
	rootMap := make(map[int64][]entity.Comment)
	for _, c := range result.Comments {
		if c.Depth == 0 {
			rootMap[c.ID] = []entity.Comment{c}
		} else {
			if len(c.Path) > 0 {
				rootID := c.Path[0]
				rootMap[rootID] = append(rootMap[rootID], c)
			}
		}
	}

	// проходимся по всем корням
	// строим дерево для каждого корня
	for _, comments := range rootMap {
		tree := utils.BuildTree(comments)
		if tree != nil {
			trees = append(trees, tree)
		}
	}

	// считаем страницы
	pages := result.Total / result.Limit
	if result.Total%result.Limit != 0 {
		pages++
	}
	// текущая страница
	page := result.Offset/result.Limit + 1

	resp := response.PaginatedCommentsResponse{
		Comments: trees,
		Total:    result.Total,
		Limit:    result.Limit,
		Offset:   result.Offset,
		Page:     page,
		Pages:    pages,
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

// @Summary Delete comments
// @Description Deletes comment by ID
// @Tags comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 204
// @Failure 400 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /v1/comments/{id} [delete]
func (r *V1) deleteCommentTree(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid comment id")
	}

	err = r.c.DeleteCommentWithChildren(ctx.UserContext(), int64(id))
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "comment not found")
		}
		r.l.Error(err, "restapi - v1 - deleteCommentTree")

		return errorResponse(ctx, http.StatusInternalServerError, "storage problems")
	}

	return ctx.SendStatus(http.StatusNoContent)
}
