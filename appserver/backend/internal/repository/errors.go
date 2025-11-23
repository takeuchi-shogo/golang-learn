package repository

import "errors"

// リポジトリ層で発生するエラーを定義
// 使用例:
//   - errors.Is(err, repository.ErrUserNotFound) でエラーの種類を判定
//   - HTTPステータスコードのマッピングに利用
var (
	// ErrUserNotFound はユーザーが見つからない場合のエラー
	ErrUserNotFound = errors.New("user not found")

	// ErrDuplicateEmail はメールアドレスが重複している場合のエラー
	ErrDuplicateEmail = errors.New("email already exists")
)
