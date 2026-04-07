package repository

import (
	"context"
	"fmt"

	"github.com/yhdessa/attends-moi/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BoardRepo struct {
	pool *pgxpool.Pool
}

func NewBoardRepo(pool *pgxpool.Pool) *BoardRepo {
	return &BoardRepo{pool: pool}
}

func (r *BoardRepo) Create(ctx context.Context, board *model.Board) error {
	query := `
		INSERT INTO boards (title, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query, board.Title, board.Description).
		Scan(&board.ID, &board.CreatedAt, &board.UpdatedAt)
}

func (r *BoardRepo) GetAll(ctx context.Context) ([]model.Board, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, title, description, created_at, updated_at FROM boards ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("query boards: %w", err)
	}
	defer rows.Close()

	boards := make([]model.Board, 0)
	for rows.Next() {
		var b model.Board
		if err := rows.Scan(&b.ID, &b.Title, &b.Description, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan board: %w", err)
		}
		boards = append(boards, b)
	}
	return boards, rows.Err()
}

func (r *BoardRepo) GetByID(ctx context.Context, id string) (*model.Board, error) {
	var b model.Board
	query := "SELECT id, title, description, created_at, updated_at FROM boards WHERE id = $1"
	err := r.pool.QueryRow(ctx, query, id).Scan(&b.ID, &b.Title, &b.Description, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get board: %w", err)
	}
	return &b, nil
}
