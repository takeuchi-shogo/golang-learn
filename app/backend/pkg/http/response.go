package http

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse はエラーレスポンスの構造体
// 統一されたエラー形式を提供する
type ErrorResponse struct {
	Error   string `json:"error"`             // エラーメッセージ
	Code    string `json:"code"`              // HTTPステータスコードのテキスト表現
	Details string `json:"details,omitempty"` // 追加の詳細情報(オプショナル)
}

// WriteJSON はJSONレスポンスを返す
// 引数:
//   - w: HTTPレスポンスライター
//   - data: レスポンスデータ
//   - status: HTTPステータスコード
// 実装:
//   1. Content-Typeヘッダーを設定
//   2. HTTPステータスコードを設定
//   3. JSONエンコードしてレスポンスを返す
// 注意事項: エンコードエラーは無視される
func WriteJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteError はエラーレスポンスを返す
// 引数:
//   - w: HTTPレスポンスライター
//   - message: エラーメッセージ
//   - status: HTTPステータスコード
// 実装: WriteJSONを使用してエラーレスポンスを返す
// 注意事項: ErrorResponse構造体を使用して統一されたエラー形式を提供
func WriteError(w http.ResponseWriter, message string, status int) {
	WriteJSON(w, ErrorResponse{
		Error: message,
		Code:  http.StatusText(status),
	}, status)
}
