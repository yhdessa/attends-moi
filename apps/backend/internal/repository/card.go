package repository

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/yhdessa/attends-moi/internal/model"
)

type CardRepo struct {
    pool *pgxpool.Pool
}

func NewCardRepo(pool *pgxpool.Pool) *CardRepo {
    return &CardRepo{pool: pool}
}

func (r *CardRepo) Create(ctx context.Context, card *model.Card) error {
    query := `
        INSERT INTO cards (board_id, title, description, status, priority, labels, assignee, due_date, position)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at, updated_at
    `
    return r.pool.QueryRow(ctx, query,
        card.BoardID, card.Title, card.Description,
        card.Status, card.Priority, card.Labels, card.Assignee, card.DueDate, card.Position,
    ).Scan(&card.ID, &card.CreatedAt, &card.UpdatedAt)
}

func (r *CardRepo) GetByBoardID(ctx context.Context, boardID string) ([]model.Card, error) {
    rows, err := r.pool.Query(ctx,
        "SELECT id, board_id, title, description, status, priority, labels, assignee, due_date, position, created_at, updated_at FROM cards WHERE board_id = $1 ORDER BY position ASC, created_at ASC",
        boardID,
    )
    if err != nil {
        return nil, fmt.Errorf("query cards: %w", err)
    }
    defer rows.Close()

    cards := make([]model.Card, 0)
    for rows.Next() {
        var c model.Card
        if err := rows.Scan(&c.ID, &c.BoardID, &c.Title, &c.Description, &c.Status, &c.Priority, &c.Labels, &c.Assignee, &c.DueDate, &c.Position, &c.CreatedAt, &c.UpdatedAt); err != nil {
            return nil, fmt.Errorf("scan card: %w", err)
        }
        cards = append(cards, c)
    }
    return cards, rows.Err()
}

func (r *CardRepo) GetByID(ctx context.Context, id string) (*model.Card, error) {
    var c model.Card
    query := "SELECT id, board_id, title, description, status, priority, labels, assignee, due_date, position, created_at, updated_at FROM cards WHERE id = $1"
    err := r.pool.QueryRow(ctx, query, id).Scan(&c.ID, &c.BoardID, &c.Title, &c.Description, &c.Status, &c.Priority, &c.Labels, &c.Assignee, &c.DueDate, &c.Position, &c.CreatedAt, &c.UpdatedAt)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("get card: %w", err)
    }
    return &c, nil
}

func (r *CardRepo) Update(ctx context.Context, card *model.Card) error {
    query := `
        UPDATE cards SET title = $2, description = $3, status = $4, priority = $5, labels = $6, assignee = $7, due_date = $8, position = $9, updated_at = NOW()
        WHERE id = $1
        RETURNING updated_at
    `
    return r.pool.QueryRow(ctx, query,
        card.ID, card.Title, card.Description, card.Status, card.Priority, card.Labels, card.Assignee, card.DueDate, card.Position,
    ).Scan(&card.UpdatedAt)
}

func (r *CardRepo) Delete(ctx context.Context, id string) error {
    _, err := r.pool.Exec(ctx, "DELETE FROM cards WHERE id = $1", id)
    return err
}
