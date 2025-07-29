# 图书推荐系统迁移总结

## 迁移概述

本次迁移将图书推荐系统从基于借阅记录的推荐方式完全转换为基于用户在图书推荐页面行为数据的推荐系统。

## 已完成的工作

### 1. 删除借阅记录相关内容

#### 删除的文件：
- `internal/domain/model/borrow.go` - 借阅记录模型
- `book_borrow_info.sql` - 借阅记录数据库表结构

#### 修改的接口：
- 从 `internal/domain/service/recommendation_service.go` 中删除了 `RecordBookBorrow` 方法

### 2. 增强用户行为追踪功能

#### 新增的行为类型：
- **浏览行为** (`view`): 用户浏览图书详情页
- **点击行为** (`click`): 用户点击图书链接
- **阅读行为** (`read`): 用户深度阅读图书内容
- **停留时间** (`stay_time`): 用户在图书页面的停留时间

#### 新增的API接口：
- `POST /books/stay-time` - 记录图书页面停留时间
- `POST /behavior/track` - 统一的用户行为追踪接口

#### 智能行为分析：
- 停留时间超过30秒自动视为深度阅读行为
- 不同行为类型对应不同的推荐权重

### 3. 完善推荐系统

#### 推荐算法优化：
- **个性化推荐**: 基于用户行为的协同过滤
- **热门推荐**: 基于所有用户行为统计的热门度排序
- **相似推荐**: 基于用户行为模式的物品协同过滤

#### 返回数据增强：
- 推荐结果现在返回完整的图书信息，而不仅仅是标题
- 包含推荐原因和算法说明
- 提供推荐分数和类别信息

#### 新增的推荐API：
- `GET /behavior/recommendations` - 基于用户行为的个性化推荐
- `GET /behavior/popular` - 基于用户行为的热门图书
- `GET /behavior/similar` - 基于用户行为的相似图书

### 4. Gorse配置优化

#### 更新的配置项：
```toml
# 正向反馈类型（表示用户喜欢）
positive_feedback_types = ["click", "read"]

# 阅读反馈类型（表示用户查看）
read_feedback_types = ["read", "view"]
```

### 5. 系统架构改进

#### 新增的组件：
- `api/behavior_tracking_handler.go` - 统一的行为追踪处理器
- `internal/domain/valueobject/book_status.go` - 图书状态值对象

#### 服务层增强：
- `BookService` 新增 `RecordBookStayTime` 方法
- 所有推荐方法现在返回完整的 `BookInfo` 对象
- 添加 `getBooksByTitles` 方法用于批量获取图书详细信息

## 系统特点

### 1. 完全基于用户行为
- 不再依赖借阅记录
- 实时收集用户在推荐页面的行为数据
- 支持多种行为类型的综合分析

### 2. 智能推荐算法
- 基于Gorse推荐引擎的协同过滤算法
- 支持冷启动用户的默认推荐策略
- 实时更新推荐结果

### 3. 丰富的API接口
- 保持向后兼容性
- 提供新的统一行为追踪接口
- 返回完整的图书信息和推荐元数据

### 4. 可扩展的架构
- 清晰的分层架构
- 易于添加新的行为类型
- 支持多种推荐策略

## 使用指南

### 前端集成建议

1. **行为数据收集**：
   ```javascript
   // 记录用户点击
   fetch('/behavior/track', {
     method: 'POST',
     headers: { 'Content-Type': 'application/json' },
     body: JSON.stringify({
       user_id: 'user123',
       book_title: '深入理解计算机系统',
       behavior_type: 'click'
     })
   });

   // 记录停留时间
   fetch('/behavior/track', {
     method: 'POST',
     headers: { 'Content-Type': 'application/json' },
     body: JSON.stringify({
       user_id: 'user123',
       book_title: '深入理解计算机系统',
       behavior_type: 'stay_time',
       stay_time_seconds: 45
     })
   });
   ```

2. **获取推荐**：
   ```javascript
   // 获取个性化推荐
   fetch('/behavior/recommendations?user_id=user123&limit=10')
     .then(response => response.json())
     .then(data => {
       console.log('推荐图书:', data.recommendations);
       console.log('推荐算法:', data.algorithm);
     });
   ```

### 部署注意事项

1. **Gorse服务**: 确保Gorse推荐引擎正常运行
2. **数据库**: 确保PostgreSQL数据库包含图书信息表
3. **配置**: 更新Gorse配置文件中的反馈类型设置

## 兼容性

- 保留了所有原有的API接口，确保现有前端代码无需修改
- 新增的接口提供更丰富的功能和更好的用户体验
- 推荐逐步迁移到新的API接口以获得更好的推荐效果

## 下一步建议

1. **测试验证**: 全面测试所有API接口的功能
2. **性能优化**: 监控推荐系统的响应时间和准确性
3. **数据分析**: 分析用户行为数据，优化推荐算法参数
4. **前端迁移**: 逐步将前端代码迁移到新的API接口
