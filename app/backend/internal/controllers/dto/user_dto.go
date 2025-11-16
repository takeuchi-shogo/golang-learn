package dto

import (
	"errors"
	"regexp"
)

// メールアドレスの正規表現パターン
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// UpdateUserRequest はユーザー更新リクエストの構造体
// ポインタ型を使用することで、部分更新とnull値の区別が可能
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// Validate はリクエストのバリデーションを行う
// 戻り値:
//   - error: バリデーションエラー
// 実装:
//   1. 少なくとも1つのフィールドが提供されているか確認
//   2. nameが提供されている場合、空文字列でないか、長さが100文字以下か確認
//   3. emailが提供されている場合、正しいメールアドレス形式か確認
// 注意事項:
//   - nameとemailの両方がnilの場合はエラー
//   - nameが空文字列の場合はエラー
//   - emailの形式が不正な場合はエラー
func (r *UpdateUserRequest) Validate() error {
	// 少なくとも1つのフィールドが提供されているか確認
	if r.Name == nil && r.Email == nil {
		return errors.New("at least one field must be provided")
	}

	// nameのバリデーション
	if r.Name != nil {
		if len(*r.Name) == 0 {
			return errors.New("name cannot be empty")
		}
		if len(*r.Name) > 100 {
			return errors.New("name must be 100 characters or less")
		}
	}

	// emailのバリデーション
	if r.Email != nil {
		if !emailRegex.MatchString(*r.Email) {
			return errors.New("invalid email format")
		}
	}

	return nil
}
