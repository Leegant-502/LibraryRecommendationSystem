package api

import (
	"encoding/json"
	"library/internal/domain/model"
	"library/internal/service"
	"net/http"
	"time"
)

type BehaviorHandler struct {
	behaviorService *service.BehaviorService
}

func NewBehaviorHandler(behaviorService *service.BehaviorService) *BehaviorHandler {
	return &BehaviorHandler{
		behaviorService: behaviorService,
	}
}

func (h *BehaviorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/behaviors":
		h.TrackBehavior(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/behaviors/stats":
		h.GetBehaviorStats(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *BehaviorHandler) TrackBehavior(w http.ResponseWriter, r *http.Request) {
	var behavior model.UserBehavior
	if err := json.NewDecoder(r.Body).Decode(&behavior); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	behavior.Timestamp = time.Now()
	if err := h.behaviorService.TrackBehavior(&behavior); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *BehaviorHandler) GetBehaviorStats(w http.ResponseWriter, r *http.Request) {
	bookID := r.URL.Query().Get("book_id")
	if bookID == "" {
		http.Error(w, "book_id is required", http.StatusBadRequest)
		return
	}

	// 默认获取最近24小时的数据
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	stats, err := h.behaviorService.GetBookBehaviors(bookID, startTime, endTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}
