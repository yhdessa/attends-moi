package handler

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/yhdessa/attends-moi/internal/model"
    "github.com/yhdessa/attends-moi/internal/repository"
)

type CardHandler struct {
    repo *repository.CardRepo
}

func NewCardHandler(repo *repository.CardRepo) *CardHandler {
    return &CardHandler{repo: repo}
}

func (h *CardHandler) Create(w http.ResponseWriter, r *http.Request) {
    boardID := chi.URLParam(r, "boardID")

    var input struct {
        Title       string             `json:"title"`
        Description string             `json:"description"`
        Status      model.CardStatus   `json:"status"`
        Priority    model.CardPriority `json:"priority"`
        Labels      []string           `json:"labels"`
        Assignee    string             `json:"assignee"`
        DueDate     *time.Time         `json:"due_date"`
        Position    *float64           `json:"position"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    position := float64(time.Now().UnixMilli())
    if input.Position != nil {
        position = *input.Position
    }

    card := &model.Card{
        BoardID:     boardID,
        Title:       input.Title,
        Description: input.Description,
        Status:      input.Status,
        Priority:    input.Priority,
        Labels:      input.Labels,
        Assignee:    input.Assignee,
        DueDate:     input.DueDate,
        Position:    position,
    }
    if card.Status == "" {
        card.Status = model.StatusBacklog
    }
    if card.Priority == "" {
        card.Priority = model.PriorityMedium
    }

    if err := h.repo.Create(r.Context(), card); err != nil {
        http.Error(w, "failed to create card", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) GetByBoardID(w http.ResponseWriter, r *http.Request) {
    boardID := chi.URLParam(r, "boardID")
    cards, err := h.repo.GetByBoardID(r.Context(), boardID)
    if err != nil {
        http.Error(w, "failed to fetch cards", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cards)
}

func (h *CardHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    card, err := h.repo.GetByID(r.Context(), id)
    if err != nil {
        http.Error(w, "failed to fetch card", http.StatusInternalServerError)
        return
    }
    if card == nil {
        http.Error(w, "card not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) Update(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    card, err := h.repo.GetByID(r.Context(), id)
    if err != nil {
        http.Error(w, "failed to fetch card", http.StatusInternalServerError)
        return
    }
    if card == nil {
        http.Error(w, "card not found", http.StatusNotFound)
        return
    }

    var input struct {
        Title       *string             `json:"title"`
        Description *string             `json:"description"`
        Status      *model.CardStatus   `json:"status"`
        Priority    *model.CardPriority `json:"priority"`
        Labels      *[]string           `json:"labels"`
        Assignee    *string             `json:"assignee"`
        DueDate     *time.Time          `json:"due_date"`
        Position    *float64            `json:"position"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if input.Title != nil {
        card.Title = *input.Title
    }
    if input.Description != nil {
        card.Description = *input.Description
    }
    if input.Status != nil {
        card.Status = *input.Status
    }
    if input.Priority != nil {
        card.Priority = *input.Priority
    }
    if input.Labels != nil {
        card.Labels = *input.Labels
    }
    if input.Assignee != nil {
        card.Assignee = *input.Assignee
    }
    if input.DueDate != nil {
        card.DueDate = input.DueDate
    }
    if input.Position != nil {
        card.Position = *input.Position
    }

    if err := h.repo.Update(r.Context(), card); err != nil {
        http.Error(w, "failed to update card", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if err := h.repo.Delete(r.Context(), id); err != nil {
        http.Error(w, "failed to delete card", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
