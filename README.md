# HPA - Home Private Academy

> 轻量级私域学习网站 MVP 文档模板

## 项目简介

HPA（Home Private Academy）是一个面向个人知识分享者的私域网站项目文档模板，包含从立项到运营的完整文档体系。

**适用场景**：
- 个人博客 + 付费课程
- 知识付费平台
- 会员制内容网站

**技术栈**：
- 后端：Go + Gin + SQLite
- 前端：Vue 3 + Vite
- 部署：Docker 单容器

## 文档结构

```
docs/
├── README.md                 # 项目导航
├── 01-立项/
│   ├── 项目概述.md           # 项目背景、目标、范围
│   └── 需求分析.md           # 功能需求、用户故事
├── 02-设计/
│   ├── 技术架构.md           # 系统架构、技术选型
│   ├── 数据库设计.md         # 数据模型、表结构
│   └── 接口设计.md           # API 规范、安全设计
├── 03-测试/
│   └── 测试用例.md           # 功能/安全/性能测试
├── 04-部署/
│   └── 部署指南.md           # Docker 部署、Nginx 配置
└── 05-运营/
    └── 使用手册.md           # 内容管理、日常运维
```

## 核心功能

| 功能 | 说明 |
|-----|------|
| 笔记展示 | Markdown 笔记渲染，游客可浏览 |
| 用户系统 | 手机号+验证码注册，用户名自动生成 |
| 课程销售 | 第三方支付 + 邀请码兑换 |
| 会员体系 | 年消费满额解锁全站下载权限 |
| 安全防护 | 频率限制、游客排队、注入防护 |

## 快速开始

### 1. 克隆仓库

```bash
git clone https://github.com/your-username/hpa-docs.git
cd hpa-docs
```

### 2. 浏览文档

推荐使用支持 Markdown 的编辑器（VS Code、Typora 等）查看 `docs/` 目录下的文档。

### 3. 基于模板开发

1. Fork 本仓库
2. 根据实际需求修改文档
3. 按照文档进行开发

## 文档亮点

- **完整的项目生命周期**：从立项到运营全覆盖
- **详细的安全设计**：注入防护、CSRF、频率限制等，带优先级标注
- **可落地的技术方案**：包含代码示例和配置模板
- **清晰的业务规则**：用户角色、支付流程、会员体系

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 许可证

本项目采用 [MIT License](LICENSE) 开源。

## 联系方式

如有问题或建议，请提交 [Issue](https://github.com/your-username/hpa-docs/issues)。

---

**Star** 本项目如果对你有帮助！
