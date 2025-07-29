# 图书推荐系统项目结构

## 项目概述

本项目是一个基于用户行为数据的图书推荐系统，使用 Gorse 推荐引擎，完全摒弃了借阅记录，专注于用户在图书推荐页面的交互行为。

## 目录结构

```
LibraryRecommendationSystem/
├── api/                                    # API处理层
│   ├── behavior_tracking_handler.go       # 用户行为追踪处理器
│   └── book_handler.go                     # 图书相关处理器
├── config/                                 # 配置管理
│   └── config.go                          # 配置文件
├── gorse-deploy/                          # Gorse部署配置
│   ├── config.toml                        # Gorse配置文件
│   └── docker-compose.yml                 # Docker部署配置
├── internal/                              # 内部业务逻辑
│   ├── bookFetch/                         # 图书数据获取
│   │   └── bookFetch.go                   # 从外部API获取图书数据
│   ├── domain/                            # 领域模型
│   │   ├── model/                         # 数据模型
│   │   │   ├── batch.go                   # 批量操作模型
│   │   │   ├── book.go                    # 图书信息模型
│   │   │   ├── recommendation.go          # 推荐响应模型
│   │   │   ├── user.go                    # 用户模型
│   │   │   └── user_behavior.go           # 用户行为模型
│   │   ├── repository/                    # 仓储接口
│   │   │   └── book_repository.go         # 图书仓储接口
│   │   ├── service/                       # 领域服务接口
│   │   │   └── recommendation_service.go  # 推荐服务接口（已删除）
│   │   └── valueobject/                   # 值对象
│   │       └── book_status.go             # 图书状态枚举
│   ├── gorse/                             # Gorse客户端
│   │   └── client.go                      # Gorse API客户端
│   ├── repository/                        # 仓储实现
│   │   └── book_repository.go             # 图书仓储实现
│   └── service/                           # 业务服务
│       └── book_service.go                # 图书业务服务
├── routes/                                # 路由配置
│   └── routes.go                          # API路由设置
├── book_info.sql                          # 图书信息表结构
├── go.mod                                 # Go模块文件
├── go.sum                                 # Go依赖校验文件
├── main.go                                # 应用入口
├── test_api.go                            # API测试脚本
├── API_USAGE.md                           # API使用指南
├── MIGRATION_SUMMARY.md                   # 迁移总结
└── PROJECT_STRUCTURE.md                   # 项目结构说明（本文件）
```

## 核心组件说明

### 1. API层 (`api/`)

- **behavior_tracking_handler.go**: 统一的用户行为追踪处理器
  - 支持点击、浏览、阅读、停留时间等行为类型
  - 提供统一的 `/behavior/track` 接口
  - 提供基于行为的推荐API

- **book_handler.go**: 图书相关的API处理器
  - 保留原有的分离式行为记录接口（向后兼容）
  - 提供图书推荐、热门图书、相似图书等接口

### 2. 领域模型 (`internal/domain/`)

- **model/**: 核心数据模型
  - `book.go`: 图书信息模型，包含完整的图书属性
  - `user_behavior.go`: 用户行为模型，记录用户交互数据
  - `recommendation.go`: 推荐响应模型，定义推荐结果格式

- **valueobject/**: 值对象
  - `book_status.go`: 图书状态枚举（可借阅、已借出等）

### 3. 业务服务 (`internal/service/`)

- **book_service.go**: 核心业务服务
  - 用户行为记录：`RecordBookView`, `RecordBookClick`, `RecordBookRead`, `RecordBookStayTime`
  - 推荐获取：`GetRecommendations`, `GetPopularBooks`, `GetSimilarBooks`
  - 与Gorse推荐引擎的集成

### 4. Gorse集成 (`internal/gorse/`)

- **client.go**: Gorse推荐引擎客户端
  - 用户反馈数据提交
  - 推荐结果获取
  - 热门和相似物品获取

### 5. 数据访问 (`internal/repository/`)

- **book_repository.go**: 图书数据访问层
  - 图书信息的CRUD操作
  - 批量查询和标题查询功能

## 数据流程

### 用户行为收集流程
1. 前端用户交互 → API接口 → BookService → Gorse客户端 → Gorse推荐引擎
2. 同时保存行为数据到本地数据库（UserBehavior表）

### 推荐获取流程
1. 前端请求推荐 → API接口 → BookService → Gorse客户端 → Gorse推荐引擎
2. 获取推荐的图书标题 → BookRepository → 查询完整图书信息 → 返回给前端

## 关键特性

### 1. 行为驱动推荐
- 完全基于用户在推荐页面的行为数据
- 支持多种行为类型：点击、浏览、阅读、停留时间
- 智能行为分析：停留时间超过30秒视为深度阅读

### 2. 实时推荐
- 用户行为数据实时传递给Gorse
- 推荐结果基于最新的用户行为模式
- 支持冷启动用户的默认推荐策略

### 3. 完整数据返回
- 推荐结果包含完整的图书信息
- 提供推荐原因和算法说明
- 支持推荐结果的元数据

### 4. 向后兼容
- 保留原有的API接口
- 新增更强大的统一行为追踪接口
- 渐进式迁移支持

## 配置要求

### 1. 数据库配置
- PostgreSQL数据库
- 自动创建 `book_information` 和 `user_behaviors` 表

### 2. Gorse配置
- Gorse推荐引擎服务
- 配置文件：`gorse-deploy/config.toml`
- 反馈类型：`positive_feedback_types = ["click", "read"]`

### 3. 应用配置
- 服务器端口配置
- 数据库连接信息
- Gorse服务端点和API密钥

## 部署说明

1. **启动Gorse服务**：
   ```bash
   cd gorse-deploy
   docker-compose up -d
   ```

2. **配置数据库**：
   - 确保PostgreSQL运行
   - 配置数据库连接信息

3. **启动应用**：
   ```bash
   go build -o library.exe .
   ./library.exe
   ```

4. **测试API**：
   ```bash
   go run test_api.go
   ```

## 扩展建议

1. **行为类型扩展**：可以添加更多行为类型（如收藏、分享等）
2. **推荐算法优化**：调整Gorse配置参数以优化推荐效果
3. **性能监控**：添加推荐系统的性能监控和分析
4. **A/B测试**：支持不同推荐策略的A/B测试
