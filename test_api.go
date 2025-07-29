package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	baseURL := "http://localhost:8080"

	// 测试用户行为追踪
	fmt.Println("=== 测试用户行为追踪 ===")

	// 测试点击行为
	testBehaviorTracking(baseURL, map[string]interface{}{
		"user_id":       "test_user_001",
		"book_title":    "深入理解计算机系统",
		"behavior_type": "click",
	})

	// 测试浏览行为
	testBehaviorTracking(baseURL, map[string]interface{}{
		"user_id":       "test_user_001",
		"book_title":    "深入理解计算机系统",
		"behavior_type": "view",
	})

	// 测试停留时间行为
	testBehaviorTracking(baseURL, map[string]interface{}{
		"user_id":           "test_user_001",
		"book_title":        "深入理解计算机系统",
		"behavior_type":     "stay_time",
		"stay_time_seconds": 45,
	})

	// 测试阅读行为
	testBehaviorTracking(baseURL, map[string]interface{}{
		"user_id":           "test_user_001",
		"book_title":        "深入理解计算机系统",
		"behavior_type":     "read",
		"read_time_minutes": 10,
	})

	// 等待一下让数据处理
	time.Sleep(2 * time.Second)

	// 测试推荐获取
	fmt.Println("\n=== 测试推荐获取 ===")

	// 测试个性化推荐
	testGetRecommendations(baseURL, "test_user_001")

	// 测试热门图书
	testGetPopularBooks(baseURL)

	// 测试相似图书
	testGetSimilarBooks(baseURL, "深入理解计算机系统")
}

func testBehaviorTracking(baseURL string, data map[string]interface{}) {
	url := baseURL + "/behavior/track"

	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ 行为追踪失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Printf("✅ 行为追踪成功 (%s): %s\n", data["behavior_type"], string(body))
	} else {
		fmt.Printf("❌ 行为追踪失败 (%s): %s\n", data["behavior_type"], string(body))
	}
}

func testGetRecommendations(baseURL, userID string) {
	url := fmt.Sprintf("%s/behavior/recommendations?user_id=%s&limit=5", baseURL, userID)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ 获取推荐失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		fmt.Printf("✅ 个性化推荐成功: 获取到 %v 本图书\n", result["count"])
		if recommendations, ok := result["recommendations"].([]interface{}); ok && len(recommendations) > 0 {
			if book, ok := recommendations[0].(map[string]interface{}); ok {
				fmt.Printf("   第一本推荐图书: %s\n", book["title"])
			}
		}
	} else {
		fmt.Printf("❌ 获取推荐失败: %s\n", string(body))
	}
}

func testGetPopularBooks(baseURL string) {
	url := fmt.Sprintf("%s/behavior/popular?limit=5", baseURL)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ 获取热门图书失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		fmt.Printf("✅ 热门图书获取成功: 获取到 %v 本图书\n", result["count"])
	} else {
		fmt.Printf("❌ 获取热门图书失败: %s\n", string(body))
	}
}

func testGetSimilarBooks(baseURL, title string) {
	url := fmt.Sprintf("%s/behavior/similar?title=%s&limit=5", baseURL, title)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ 获取相似图书失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		fmt.Printf("✅ 相似图书获取成功: 获取到 %v 本图书\n", result["count"])
	} else {
		fmt.Printf("❌ 获取相似图书失败: %s\n", string(body))
	}
}
