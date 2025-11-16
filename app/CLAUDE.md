---
project_name: "Golang Learn Application"
version: "1.0.0"
license: "MIT"
created_date: "2025-11-16"
updated_date: "2025-11-16"
---

# Overview

Golang の高度な使い方を学ぶためのアプリケーションです。

## Local Setting

- 振る舞いなどの個人設定は `CLAUDE.local.md` に記載しています。

## Context

<!-- プロジェクトの背景を3行で要約 -->
- **What**: Golang の高度な使い方を学ぶため
- **Why**: 自身の実務で活かせるように学びたいため
- **Who**: 自分自身

## Project Structure

```
golang-learn/
├── backend/
|   ├── base
|   |   ├── idtype/
|   |   └── server/
|   ├── cmd/
|   |   ├── server/
|   |   └── main.go
|   ├── internal/
|   |   ├── controllers/
|   |   |   ├── authcontroller/
|   |   |   └── usercontroller/
|   |   ├── domain/
|   |   |   ├── model/
|   |   |   └── service/
|   |   ├── application/
|   |   |   ├── authapplication/
|   |   |   └── userapplication/
|   |   ├── infra/
|   |   |   ├── database/
|   |   |   └── routes/
|   |   ├── middleware/
|   |   ├── repository/
|   |   ├── service/
|   |   |   ├── command/
|   |   |   └── query/
|   |   └── factory/
|   |      ├── authregistory/
|   |      └── userregistory/
|   ├── pkg/
|   |   ├── aws/
|   |   ├── http/
|   |   ├── jwt/
|   |   └── validator/
|   ├── dbconfig.yml
|   ├── go.mod
|   ├── go.sum
|   ├── .air.toml
|   └── .gitignore
├── local/
├── server/
├── compose.yml
├── README.md
└── Taskfile.yml
```

## Tech Stack

- Go 1.25.1
- MySQL 8.0
- Docker
- Docker Compose
- Task
- Air
- sql-migrate

## Dependency

- go.mod, go.sum を参照してください。

## Coding Guidelines

- 基本的に Golang のコーディング規約に従ってください。
- また、迷う場合はこちらを参照してください。
  - [Google Style Guide](https://google.github.io/styleguide/go/)
  - [Effective Go](https://go.dev/doc/effective_go)
  - [Uber Go Style Guide](https://github.com/uber-go/guide)

## Branch / Commit

- ブランチは master ブランチから派生させてください。
- コミットメッセージは以下の形式に従ってください。
  - feat: 新機能の追加
  - fix: バグの修正
  - chore: ビルド関連の変更
  - docs: ドキュメントの更新
  - style: コードのスタイル修正
  - refactor: コードのリファクタリング
  - test: テストの追加や修正
  - perf: パフォーマンスの改善
  - ci: CI関連の変更

## Pull Request

- プルリクエストは master ブランチに対して行ってください。
