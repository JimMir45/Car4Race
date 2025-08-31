# B站视频管理系统

一个用于管理B站个人视频的Web应用，支持视频数据获取、统计分析和可视化展示。

## 功能特点

- 自动获取B站个人视频数据
- 视频列表展示（支持分页）
- 数据统计分析（播放量、点赞、硬币、收藏等）
- 分类统计图表
- 搜索、筛选、排序功能
- 实时数据刷新

## 安装步骤

1. 进入项目目录：
```bash
cd bilibili-video-manager
```

2. 安装依赖：
```bash
npm install
```

3. 配置Cookie（重要）：
   - 编辑 `src/config.js` 文件
   - 将 `COOKIE` 字段替换为你的B站Cookie

## 使用方法

1. 首次获取视频数据：
```bash
npm run fetch
```

2. 启动Web服务器：
```bash
npm run dev
```

3. 打开浏览器访问：
```
http://localhost:3000
```

## 功能说明

### 数据获取
- 自动获取所有视频信息
- 支持批量获取（每页50条）
- 数据保存在 `data/` 目录

### Web界面
- **统计概览**：显示总视频数、总播放量等关键指标
- **分类图表**：可视化展示视频分类分布
- **视频列表**：卡片式展示视频信息
- **筛选功能**：
  - 按标题搜索
  - 按分类筛选
  - 按状态筛选（已发布/待发布）
  - 多种排序方式

### API接口
- `GET /api/videos` - 获取所有视频
- `GET /api/statistics` - 获取统计数据
- `GET /api/video/:bvid` - 获取单个视频详情
- `POST /api/refresh` - 刷新数据
- `GET /api/statistics/categories` - 分类统计
- `GET /api/statistics/timeline` - 时间趋势

## 注意事项

1. **Cookie获取方法**：
   - 登录B站创作中心
   - 打开浏览器开发者工具（F12）
   - 切换到Network标签
   - 刷新页面，找到任意请求
   - 复制请求头中的Cookie值

2. **Cookie有效期**：
   - B站Cookie有有效期限制
   - 如果数据获取失败，请更新Cookie

3. **请求频率**：
   - 避免频繁刷新数据
   - 建议间隔至少5分钟

## 数据文件说明

- `data/videos.json` - 视频详细信息
- `data/statistics.json` - 统计汇总数据

## 技术栈

- Node.js + Express（后端）
- 原生JavaScript + HTML + CSS（前端）
- Fetch API（数据请求）

## License

MIT