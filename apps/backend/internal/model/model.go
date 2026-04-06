package model

import "time"

type Board struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CardStatus string

const (
	StatusBacklog    CardStatus = "backlog"
	StatusTodo       CardStatus = "todo"
	StatusInProgress CardStatus = "in_progress"
	StatusReview     CardStatus = "review"
	StatusDone       CardStatus = "done"
)

type CardPriority string

const (
	PriorityLow      CardPriority = "low"
	PriorityMedium   CardPriority = "medium"
	PriorityHigh     CardPriority = "high"
	PriorityCritical CardPriority = "critical"
)

type Card struct {
	ID          string     `json:"id"`
	BoardID     string     `json:"board_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      CardStatus `json:"status"`
	Priority    CardPriority `json:"priority"`
	Labels      []string   `json:"labels"`
	Assignee    string     `json:"assignee"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Comment struct {
	ID        string    `json:"id"`
	CardID    string    `json:"card_id"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}
