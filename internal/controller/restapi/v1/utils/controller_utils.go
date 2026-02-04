package utils

import (
	"database/sql"

	"github.com/andreyxaxa/Comment-Tree/internal/controller/restapi/v1/response"
	"github.com/andreyxaxa/Comment-Tree/internal/entity"
)

func NullInt64ToPtr(n sql.NullInt64) *int64 {
	if n.Valid {
		return &n.Int64
	}

	return nil
}

func BuildTree(comments []entity.Comment) *response.CommentTreeResponse {
	nodeMap := make(map[int64]*response.CommentTreeResponse)

	// для быстрого доступа по ID
	for _, c := range comments {
		nodeMap[c.ID] = &response.CommentTreeResponse{
			ID:        c.ID,
			ParentID:  NullInt64ToPtr(c.ParentID),
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
			Depth:     c.Depth,
			Children:  []*response.CommentTreeResponse{},
		}
	}

	var root *response.CommentTreeResponse

	for _, c := range comments {
		node := nodeMap[c.ID]

		// если есть родитель
		if c.ParentID.Valid {
			parent, exists := nodeMap[c.ParentID.Int64]
			if exists {
				parent.Children = append(parent.Children, node)
			} else {
				if root == nil {
					root = node
				}
			}
		} else {
			// если родителя нет - корень
			root = node
		}
	}

	if root == nil && len(comments) > 0 {
		root = nodeMap[comments[0].ID]
	}

	return root
}
