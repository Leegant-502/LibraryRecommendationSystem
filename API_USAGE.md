# 图书推荐系统 API 使用指南

## 概述

本系统基于 Gorse 推荐引擎，通过收集用户在图书推荐页面的行为数据（点击、浏览、停留时间等），为用户提供个性化的图书推荐服务。

## 用户行为追踪 API

### 1. 统一行为追踪接口（推荐使用）

**POST** `/behavior/track`

用于记录用户在图书推荐页面的各种行为。

#### 请求参数

```json
{
  "user_id": "用户ID",
  "book_title": "图书标题",
  "behavior_type": "行为类型", // "view", "click", "read", "stay_time"
  "stay_time_seconds": 30,    // 停留时间（秒），仅当 behavior_type 为 "stay_time" 时需要
  "read_time_minutes": 5,     // 阅读时间（分钟），仅当 behavior_type 为 "read" 时需要
  "extra": {}                 // 额外信息（可选）
}
```

#### 行为类型说明

- `view`: 浏览图书详情页
- `click`: 点击图书链接
- `read`: 深度阅读图书内容
- `stay_time`: 在图书页面的停留时间

#### 响应示例

```json
{
  "message": "用户行为记录成功",
  "behavior_type": "click",
  "user_id": "user123",
  "book_title": "数据结构与算法"
}
```

### 2. 分别的行为记录接口（兼容旧版本）

#### 记录图书浏览
**POST** `/books/view`
```json
{
  "user_id": "用户ID",
  "title": "图书标题"
}
```

#### 记录图书点击
**POST** `/books/click`
```json
{
  "user_id": "用户ID",
  "title": "图书标题"
}
```

#### 记录图书阅读
**POST** `/books/read`
```json
{
  "user_id": "用户ID",
  "title": "图书标题",
  "read_time_minutes": 5
}
```

#### 记录停留时间
**POST** `/books/stay-time`
```json
{
  "user_id": "用户ID",
  "title": "图书标题",
  "stay_time_seconds": 30
}
```

## 推荐获取 API

### 1. 个性化推荐（推荐使用）

**GET** `/behavior/recommendations?user_id={用户ID}&limit={数量}`

基于用户行为数据的个性化推荐。

#### 响应示例

```json
{
  "success": true,
  "recommendations": [
    {
      "id": "123",
      "title": "深入理解计算机系统",
      "primary_author": "Randal E. Bryant",
      "publisher": "机械工业出版社",
      "publication_date": "2016-11-01T00:00:00Z",
      // ... 其他图书信息
    }
  ],
  "count": 10,
  "user_id": "user123",
  "algorithm": "基于用户行为的协同过滤推荐",
  "message": "推荐结果基于您的浏览、点击和停留时间等行为数据生成"
}
```

### 2. 热门图书推荐

**GET** `/behavior/popular?limit={数量}`

基于所有用户行为统计的热门图书。

#### 响应示例

```json
{
  "success": true,
  "popular_books": [
    {
      "id": "456",
      "title": "算法导论",
      "primary_author": "Thomas H. Cormen",
      // ... 其他图书信息
    }
  ],
  "count": 10,
  "algorithm": "基于用户行为统计的热门度排序",
  "message": "热门图书基于所有用户的点击、浏览和停留时间等行为数据统计生成"
}
```

### 3. 相似图书推荐

**GET** `/behavior/similar?title={图书标题}&limit={数量}`

基于用户行为模式的相似图书推荐。

#### 响应示例

```json
{
  "success": true,
  "similar_books": [
    {
      "id": "789",
      "title": "数据结构与算法分析",
      "primary_author": "Mark Allen Weiss",
      // ... 其他图书信息
    }
  ],
  "count": 10,
  "base_title": "数据结构与算法",
  "algorithm": "基于用户行为的物品协同过滤",
  "message": "相似图书基于用户对图书的行为模式相似性推荐"
}
```

## 兼容性接口

为了保持向后兼容，系统仍然保留了原有的推荐接口：

- **GET** `/recommendations?user_id={用户ID}&limit={数量}` - 个性化推荐
- **GET** `/books/popular?limit={数量}` - 热门图书
- **GET** `/books/similar?title={图书标题}&limit={数量}` - 相似图书

## 使用建议

1. **前端集成**：建议使用新的 `/behavior/track` 接口统一记录用户行为
2. **行为数据收集**：
   - 用户点击图书链接时记录 `click` 行为
   - 用户浏览图书详情页时记录 `view` 行为
   - 用户在页面停留超过一定时间时记录 `stay_time` 行为
   - 用户深度阅读时记录 `read` 行为

3. **推荐获取**：建议使用新的 `/behavior/` 系列接口获取推荐结果

## 系统特点

- **完全基于用户行为**：不再依赖借阅记录，完全基于用户在推荐页面的行为数据
- **实时推荐**：用户行为数据实时传递给 Gorse 推荐引擎
- **多种推荐算法**：支持个性化推荐、热门推荐、相似推荐
- **完整图书信息**：推荐结果包含图书的完整信息，而不仅仅是标题
- **智能行为分析**：根据停留时间等行为特征智能判断用户兴趣程度
