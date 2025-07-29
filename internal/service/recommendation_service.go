package service

import (
	"library/internal/domain/model"
	"library/internal/gorse"
	"library/internal/repository"
)

// RecommendationService 推荐服务
type RecommendationService struct {
	gorseClient  *gorse.Client
	bookRepo     repository.BookRepository
	behaviorRepo repository.BehaviorRepository
}

// NewRecommendationService 创建推荐服务实例
func NewRecommendationService(
	gorseClient *gorse.Client,
	bookRepo repository.BookRepository,
	behaviorRepo repository.BehaviorRepository,
) *RecommendationService {
	return &RecommendationService{
		gorseClient:  gorseClient,
		bookRepo:     bookRepo,
		behaviorRepo: behaviorRepo,
	}
}

// GetPersonalizedRecommendations 获取个性化推荐
func (s *RecommendationService) GetPersonalizedRecommendations(userID string, limit int) (*model.RecommendationResponse, error) {
	// 1. 从 Gorse 获取推荐的图书 ID
	itemIDs, err := s.gorseClient.GetRecommend(userID, "", limit)
	if err != nil {
		return nil, err
	}

	// 2. 获取图书详细信息
	books, err := s.bookRepo.FindByIDs(itemIDs)
	if err != nil {
		return nil, err
	}

	// 3. 构建响应
	return &model.RecommendationResponse{
		Books:    books,
		Category: "personalized",
		Reason:   "根据您的阅读历史推荐",
	}, nil
}

// GetCategoryRecommendations 获取分类推荐
func (s *RecommendationService) GetCategoryRecommendations(limit int) ([]model.RecommendationCategory, error) {
	categories := []struct {
		id    string
		name  string
		desc  string
		fetch func(string, int) ([]string, error)
	}{
		{
			id:    "popular",
			name:  "热门推荐",
			desc:  "最受欢迎的图书",
			fetch: s.gorseClient.GetPopular,
		},
		{
			id:    "latest",
			name:  "最新上架",
			desc:  "新到馆的图书",
			fetch: s.gorseClient.GetLatest,
		},
	}

	var results []model.RecommendationCategory
	for _, cat := range categories {
		// 获取推荐的图书 ID
		itemIDs, err := cat.fetch("", limit)
		if err != nil {
			continue // 跳过失败的类别
		}

		// 获取图书详细信息
		books, err := s.bookRepo.FindByIDs(itemIDs)
		if err != nil {
			continue
		}

		results = append(results, model.RecommendationCategory{
			ID:          cat.id,
			Name:        cat.name,
			Description: cat.desc,
			Books:       books,
			Total:       len(books),
		})
	}

	return results, nil
}

// GetSimilarBooks 获取相似图书
func (s *RecommendationService) GetSimilarBooks(bookID string, limit int) (*model.RecommendationResponse, error) {
	// 1. 从 Gorse 获取相似图书 ID
	itemIDs, err := s.gorseClient.GetItemNeighbors(bookID, "", limit)
	if err != nil {
		return nil, err
	}

	// 2. 获取图书详细信息
	books, err := s.bookRepo.FindByIDs(itemIDs)
	if err != nil {
		return nil, err
	}

	// 3. 构建响应
	return &model.RecommendationResponse{
		Books:    books,
		Category: "similar",
		Reason:   "与当前图书相似的推荐",
	}, nil
}
