package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yhdessa/attends-moi/internal/model"
	"github.com/yhdessa/attends-moi/internal/repository"
)

type CommentHandler struct {
	repo *repository.CommentRepo
}

func NewCommentHandler(repo *repository.CommentRepo) *CommentHandler {
	return &CommentHandler{repo: repo}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	cardID := chi.URLParam(r, "cardID")

	var input struct {
		Author string `json:"author"`
		Body   string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	comment := &model.Comment{CardID: cardID, Author: input.Author, Body: input.Body}
	if err := h.repo.Create(r.Context(), comment); err != nil {
		http.Error(w, "failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) GetByCardID(w http.ResponseWriter, r *http.Request) {
	cardID := chi.URLParam(r, "cardID")
	comments, err := h.repo.GetByCardID(r.Context(), cardID)
	if err != nil {
		http.Error(w, "failed to fetch comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
