---
name: create-pull-request
description: Create a pull request for the current branch
---

# Create Pull Request

## Context

- 現在の git ステータス: !‎`git status`

- 現在の git diff（ステージ済み／未ステージの変更）: !‎`git diff HEAD`

- 現在のブランチ: !‎`git branch --show-current`

## Your Task

1. `master` ブランチにいる場合は、新しいブランチを作成する
2. 適切なメッセージを入力して、単一のコミットを作成する
3. コミットメッセージは以下の形式に従ってください。

- feat: 新機能の追加
- fix: バグの修正
- chore: ビルド関連の変更
- docs: ドキュメントの更新
- style: コードのスタイル修正
- refactor: コードのリファクタリング
- test: テストの追加や修正
- perf: パフォーマンスの改善
- ci: CI関連の変更

4. ブランチを origin にプッシュする
5. プルリクエストを作成する
6. プルリクエストのタイトルは以下の形式に従ってください。

- feat: 新機能の追加
- fix: バグの修正
- chore: ビルド関連の変更
- docs: ドキュメントの更新
- style: コードのスタイル修正
- refactor: コードのリファクタリング
- test: テストの追加や修正
- perf: パフォーマンスの改善
- ci: CI関連の変更

7. プルリクエストの説明は以下の形式に従ってください。

- 変更内容の概要
- 変更内容の詳細
- 変更内容の影響範囲
- 変更内容のテスト方法
- 変更内容のテスト結果
- 変更内容のテストコード
- 変更内容のテストコードの実装
- 変更内容のテストコードの実装
