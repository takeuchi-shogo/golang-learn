package model

import (
	"github.com/google/uuid"
)

type AuthType int

const (
	AuthTypeNone AuthType = iota
	AuthTypeCognito
)

func (a AuthType) String() string {
	switch a {
	case AuthTypeNone:
		return "none"
	case AuthTypeCognito:
		return "cognito"
	}
	return "unknown"
}

type (
	User struct {
		ID          uuid.UUID `json:"id"`
		UserIDToken string    `json:"user_id_token"`
		Name        Name      `json:"name"`
		Email       Email     `json:"email"`
	}

	Name  string
	Email string
)

func newUser(id uuid.UUID, name Name, email Email, userIDToken string) *User {
	return &User{
		ID:          id,
		Name:        name,
		Email:       email,
		UserIDToken: userIDToken,
	}
}

func NewUser(name Name, email Email, userIDToken string) *User {
	id := uuid.Must(uuid.NewV7())
	return newUser(id, name, email, userIDToken)
}

func (u *User) GetID() uuid.UUID {
	return u.ID
}

func (u *User) GetName() Name {
	return u.Name
}

func (u *User) GetEmail() Email {
	return u.Email
}

// GetUserIDToken はユーザーIDトークンを取得する
// 引数: なし
// 戻り値: ユーザーIDトークン(Cognito Username)
// 実装: UserIDTokenフィールドを返す
// 注意事項: Cognitoのユーザー識別子として使用される
func (u *User) GetUserIDToken() string {
	return u.UserIDToken
}

// UpdateName はユーザーの名前を更新する
// 引数:
//   - name: 新しいユーザー名
// 戻り値: なし
// 実装: ユーザーの名前フィールドを更新する
func (u *User) UpdateName(name Name) {
	u.Name = name
}

// UpdateEmail はユーザーのメールアドレスを更新する
// 引数:
//   - email: 新しいメールアドレス
// 戻り値: なし
// 実装: ユーザーのメールアドレスフィールドを更新する
func (u *User) UpdateEmail(email Email) {
	u.Email = email
}

// ReconstructUser は既存のユーザーを再構築する
// 引数:
//   - id: ユーザーID
//   - name: ユーザー名
//   - email: メールアドレス
//   - userIDToken: ユーザーIDトークン
// 戻り値: 再構築されたユーザーオブジェクト
// 実装: DBから取得したデータを使ってユーザーを再構築する際に使用
// 注意事項: 新規ユーザー作成には使用せず、既存データの復元にのみ使用する
func ReconstructUser(id uuid.UUID, name Name, email Email, userIDToken string) *User {
	return newUser(id, name, email, userIDToken)
}
