package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yhdessa/attends-moi/internal/model"
	"github.com/yhdessa/attends-moi/internal/repository"
)

type BoardHandler struct {
	repo *repository.BoardRepo
}

func NewBoardHandler(repo *repository.BoardRepo) *BoardHandler {
	return &BoardHandler{repo: repo}
}

func (h *BoardHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	board := &model.Board{Title: input.Title, Description: input.Description}
	if err := h.repo.Create(r.Context(), board); err != nil {
		http.Error(w, "failed to create board", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(board)
}

func (h *BoardHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	boards, err := h.repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch boards", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boards)
}

func (h *BoardHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	board, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to fetch board", http.StatusInternalServerError)
		return
	}
	if board == nil {
		http.Error(w, "board not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}
