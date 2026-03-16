# Sub2API

<div align="center">

[![Go](https://img.shields.io/badge/Go-1.25.7-00ADD8.svg)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D.svg)](https://vuejs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791.svg)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D.svg)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)

**AI API 网关平台 - 订阅配额分发管理**

中文主说明 | [English](README_EN.md)

</div>

---

## 项目概述

Sub2API 是一个 AI API 网关平台，用于分发和管理 AI 产品订阅（如 Claude Code $200/月）的 API 配额。用户通过平台生成的 API Key 调用上游 AI 服务，平台负责鉴权、计费、负载均衡和请求转发。

Sub2API 将 OpenAI 兼容能力作为平台内建功能统一提供，包括：

- `/v1/responses`
- `/v1/chat/completions`
- `/v1/embeddings`

这些接口共享同一套 API Key、会员/订阅校验、账号调度、用量记录与计费链路，而不是依赖独立旁路服务。

## 核心功能

- **多账号管理**：支持多种上游账号类型（OAuth、API Key）
- **API Key 分发**：为用户生成和管理 API Key
- **统一 OpenAI v1 兼容层**：内建支持 Responses、Chat Completions、Embeddings，并与同一套网关和计费体系打通
- **精确计费**：Token 级别的用量追踪和成本计算
- **智能调度**：智能账号选择，支持粘性会话
- **并发控制**：用户级和账号级并发限制
- **速率限制**：可配置的请求和 Token 速率限制
- **管理后台**：Web 界面进行监控和管理

## 文档入口

- 完整中文说明：[README_CN.md](README_CN.md)
- 英文说明：[README_EN.md](README_EN.md)

## 在线体验

体验地址：**https://v2.pincc.ai/**

演示账号（共享演示环境；自建部署不会自动创建该账号）：

| 邮箱 | 密码 |
|------|------|
| admin@sub2api.com | admin123 |

## 快速开始

### 脚本安装

```bash
curl -sSL https://raw.githubusercontent.com/Wei-Shaw/sub2api/main/deploy/install.sh | sudo bash
```

### Docker Compose

```bash
mkdir -p sub2api-deploy && cd sub2api-deploy
curl -sSL https://raw.githubusercontent.com/Wei-Shaw/sub2api/main/deploy/docker-deploy.sh | bash
docker-compose up -d
```

## 说明

- GitHub 仓库首页默认以本文件作为中文主说明入口。
- 更完整的中文部署、配置和使用说明请阅读 [README_CN.md](README_CN.md)。
