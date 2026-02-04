package persistent

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/andreyxaxa/Comment-Tree/internal/entity"
	"github.com/andreyxaxa/Comment-Tree/pkg/postgres"
	"github.com/andreyxaxa/Comment-Tree/pkg/types/errs"
	"github.com/jackc/pgx/v5"
)

const (
	// Table
	commentsTable = "comments"

	// Columns
	idColumn        = "id"
	parentIDColumn  = "parent_id"
	contentColumn   = "content"
	createdAtColumn = "created_at"
)

type CommentRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *CommentRepo {
	return &CommentRepo{pg}
}

func (r *CommentRepo) CreateComment(ctx context.Context, parentID *int64, content string) (entity.Comment, error) {
	sqlq, args, err := r.Builder.
		Insert(commentsTable).
		Columns(parentIDColumn, contentColumn).
		Values(parentID, content).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return entity.Comment{}, fmt.Errorf("CommentRepo - CreateComment - r.Builder.ToSql(): %w", err)
	}

	c := entity.Comment{
		Content: content,
	}

	if parentID != nil {
		c.ParentID = sql.NullInt64{Int64: *parentID, Valid: true}
	} else {
		c.ParentID = sql.NullInt64{Valid: false}
	}

	err = r.Pool.QueryRow(ctx, sqlq, args...).Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		return entity.Comment{}, fmt.Errorf("CommentRepo - CreateComment - r.Pool.QueryRow.Scan: %w", err)
	}

	return c, nil
}

// returns:
// comment exists - nil;
// squirrel build error - err;
// pool.QueryRow.Scan error - err;
// comment not exists - errs.ErrRecordNotFound.
func (r *CommentRepo) CommentExists(ctx context.Context, id int64) error {
	sql, args, err := r.Builder.
		Select("1").
		From(commentsTable).
		Where(squirrel.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("CommentRepo - CommentExists - r.Builder.ToSql: %w", err)
	}

	var exists int
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("CommentRepo - CommentExists: %w", errs.ErrRecordNotFound)
		}
		return fmt.Errorf("CommentRepo - CommentExists - r.Pool.QueryRow.Scan: %w", err)
	}

	return nil
}

func (r *CommentRepo) GetCommentWithChildren(ctx context.Context, id int64) ([]entity.Comment, error) {
	sql := `
	WITH RECURSIVE comment_tree AS (
		SELECT 
			id, 
			parent_id, 
			content, 
			created_at, 
			0 AS depth, 
			ARRAY[id] AS path
		FROM comments
		WHERE id = $1

		UNION ALL

		SELECT
			c.id,
			c.parent_id,
			c.content,
			c.created_at,
			ct.depth + 1,
			ct.path || c.id
		FROM comments c
		INNER JOIN comment_tree ct ON c.parent_id = ct.id
	)
	SELECT id, parent_id, content, created_at, depth, path 
	FROM comment_tree
	ORDER BY path;
	`

	rows, err := r.Pool.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("CommentRepo - GetComments - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var comments []entity.Comment

	for rows.Next() {
		var c entity.Comment

		err = rows.Scan(
			&c.ID,
			&c.ParentID,
			&c.Content,
			&c.CreatedAt,
			&c.Depth,
			&c.Path,
		)
		if err != nil {
			return nil, fmt.Errorf("CommentRepo - GetComments - rows.Scan: %w", err)
		}

		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("CommentRepo - GetComments - rows.Err: %w", err)
	}

	if len(comments) == 0 {
		return nil, fmt.Errorf("CommentRepo - GetComments: %w", errs.ErrRecordNotFound)
	}

	return comments, nil
}

func (r *CommentRepo) DeleteCommentWithChildren(ctx context.Context, id int64) error {
	sql, args, err := r.Builder.
		Delete(commentsTable).
		Where(squirrel.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("CommentRepo - DeleteCommentWithChildren - r.Builder.ToSql: %w", err)
	}

	tag, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("CommentRepo - DeleteCommentWithChildren - r.Pool.Exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("CommentRepo - DeleteCommentWithChildren: %w", errs.ErrRecordNotFound)
	}

	return nil
}

func (r *CommentRepo) SearchComments(ctx context.Context, search string, sortBy, order string, limit, offset int) ([]entity.Comment, int, error) {
	sql := fmt.Sprintf(`
		SELECT id, parent_id, content, created_at, COUNT(*) OVER() as total
		FROM comments
		WHERE content_tsv @@ plainto_tsquery('english', $1)
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := r.Pool.Query(ctx, sql, search, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("CommentRepo - SearchComments - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var comments []entity.Comment
	var total int

	for rows.Next() {
		var c entity.Comment
		err = rows.Scan(&c.ID, &c.ParentID, &c.Content, &c.CreatedAt, &total)
		if err != nil {
			return nil, 0, fmt.Errorf("CommentRepo - SearchComments - rows.Scan: %w", err)
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("CommentRepo - SearchComments - rows.Err: %w", err)
	}

	return comments, total, nil
}

func (r *CommentRepo) GetRootComments(ctx context.Context, sortBy, order string, limit, offset int) ([]entity.Comment, int, error) {
	sql := fmt.Sprintf(`
		SELECT id, parent_id, content, created_at, COUNT(*) OVER()
		FROM comments
		WHERE parent_id IS NULL
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, sortBy, order)

	rows, err := r.Pool.Query(ctx, sql, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("CommentRepo - GetRootComments - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var comments []entity.Comment
	var total int

	for rows.Next() {
		var c entity.Comment
		err = rows.Scan(&c.ID, &c.ParentID, &c.Content, &c.CreatedAt, &total)
		if err != nil {
			return nil, 0, fmt.Errorf("CommentRepo - GetRootComments - rows.Scan: %w", err)
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("CommentRepo - GetRootComments - rows.Err: %w", err)
	}

	return comments, total, nil
}

func (r *CommentRepo) GetTreesForRoots(ctx context.Context, rootIDs []int64) ([]entity.Comment, error) {
	if len(rootIDs) == 0 {
		return []entity.Comment{}, nil
	}

	anchorQuery, anchorArgs, err := r.Builder.
		Select(idColumn, parentIDColumn, contentColumn, createdAtColumn, "0 AS depth", "ARRAY[id] AS path").
		From(commentsTable).
		Where(squirrel.Eq{idColumn: rootIDs}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CommentRepo - GetTreesForRoots - r.Builder.ToSql: %w", err)
	}

	sql := fmt.Sprintf(`
	WITH RECURSIVE comment_tree AS (
		%s

		UNION ALL

		SELECT
			c.id, 
			c.parent_id, 
			c.content,
			c.created_at,
			ct.depth + 1,
			ct.path || c.id
		FROM comments c
		INNER JOIN comment_tree ct ON c.parent_id = ct.id
	)
	SELECT id, parent_id, content, created_at, depth, path
	FROM comment_tree
	ORDER BY path;
	`, anchorQuery)

	rows, err := r.Pool.Query(ctx, sql, anchorArgs...)
	if err != nil {
		return nil, fmt.Errorf("CommentRepo - GetTreesForRoots - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var comments []entity.Comment

	for rows.Next() {
		var c entity.Comment
		err = rows.Scan(&c.ID, &c.ParentID, &c.Content, &c.CreatedAt, &c.Depth, &c.Path)
		if err != nil {
			return nil, fmt.Errorf("CommentRepo - GetTreesForRoots - rows.Scan: %w", err)
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("CommentRepo - GetTreesForRoots - rows.Err: %w", err)
	}

	return comments, nil
}
