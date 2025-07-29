package api

import (
	"encoding/json"
	"library/internal/service"
	"net/http"
	"strconv"
)

type RecommendationHandler struct {
	recommendService *service.RecommendationService
}

func NewRecommendationHandler(recommendService *service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{
		recommendService: recommendService,
	}
}

func (h *RecommendationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/recommendations/personal":
		h.GetPersonalRecommendations(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/recommendations/categories":
		h.GetCategoryRecommendations(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/recommendations/similar":
		h.GetSimilarBooks(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *RecommendationHandler) GetPersonalRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10 // 默认返回10本书
	}

	recommendations, err := h.recommendService.GetPersonalizedRecommendations(userID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(recommendations)
}

func (h *RecommendationHandler) GetCategoryRecommendations(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10 // 每个分类默认返回10本书
	}

	categories, err := h.recommendService.GetCategoryRecommendations(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(categories)
}

func (h *RecommendationHandler) GetSimilarBooks(w http.ResponseWriter, r *http.Request) {
	bookID := r.URL.Query().Get("book_id")
	if bookID == "" {
		http.Error(w, "book_id is required", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	recommendations, err := h.recommendService.GetSimilarBooks(bookID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(recommendations)
}
