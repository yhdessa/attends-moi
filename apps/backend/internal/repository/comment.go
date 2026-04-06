package repository

import (
	"context"
	"fmt"

	"github.com/yourusername/attends-moi/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepo struct {
	pool *pgxpool.Pool
}

func NewCommentRepo(pool *pgxpool.Pool) *CommentRepo {
	return &CommentRepo{pool: pool}
}

func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) error {
	query := `
		INSERT INTO comments (card_id, author, body)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query, comment.CardID, comment.Author, comment.Body).
		Scan(&comment.ID, &comment.CreatedAt)
}

func (r *CommentRepo) GetByCardID(ctx context.Context, cardID string) ([]model.Comment, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, card_id, author, body, created_at FROM comments WHERE card_id = $1 ORDER BY created_at ASC",
		cardID,
	)
	if err != nil {
		return nil, fmt.Errorf("query comments: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.CardID, &c.Author, &c.Body, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}
