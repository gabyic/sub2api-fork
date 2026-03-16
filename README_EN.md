# Sub2API

<div align="center">

[![Go](https://img.shields.io/badge/Go-1.25.7-00ADD8.svg)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D.svg)](https://vuejs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791.svg)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D.svg)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)

**AI API Gateway Platform for Subscription Quota Distribution**

[中文主说明](README.md) | English

</div>

---

## Demo

Try Sub2API online: **https://demo.sub2api.org/**

Demo credentials (shared demo environment; **not** created automatically for self-hosted installs):

| Email | Password |
|-------|----------|
| admin@sub2api.com | admin123 |

## Overview

Sub2API is an AI API gateway platform designed to distribute and manage API quotas from AI product subscriptions (like Claude Code $200/month). Users can access upstream AI services through platform-generated API Keys, while the platform handles authentication, billing, load balancing, and request forwarding.

Sub2API provides OpenAI-compatible gateway capabilities as a unified part of the platform, including `/v1/responses`, `/v1/chat/completions`, and `/v1/embeddings`. These endpoints share the same API key system, subscription and quota checks, account scheduling, usage logging, and billing pipeline instead of relying on a separate sidecar service.

## Features

- **Multi-Account Management** - Support multiple upstream account types (OAuth, API Key)
- **API Key Distribution** - Generate and manage API Keys for users
- **Unified OpenAI v1 Compatibility** - Built-in support for Responses, Chat Completions, and Embeddings under the same gateway and billing model
- **Precise Billing** - Token-level usage tracking and cost calculation
- **Smart Scheduling** - Intelligent account selection with sticky sessions
- **Concurrency Control** - Per-user and per-account concurrency limits
- **Rate Limiting** - Configurable request and token rate limits
- **Admin Dashboard** - Web interface for monitoring and management
- **External System Integration** - Embed external systems (e.g. payment, ticketing) via iframe to extend the admin dashboard

## Ecosystem

Community projects that extend or integrate with Sub2API:

| Project | Description | Features |
|---------|-------------|----------|
| [Sub2ApiPay](https://github.com/touwaeriol/sub2apipay) | Self-service payment system | Self-service top-up and subscription purchase; supports YiPay protocol, WeChat Pay, Alipay, Stripe; embeddable via iframe |
| [sub2api-mobile](https://github.com/ckken/sub2api-mobile) | Mobile admin console | Cross-platform app (iOS/Android/Web) for user management, account management, monitoring dashboard, and multi-backend switching; built with Expo + React Native |

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.25.7, Gin, Ent |
| Frontend | Vue 3.4+, Vite 5+, TailwindCSS |
| Database | PostgreSQL 15+ |
| Cache/Queue | Redis 7+ |
