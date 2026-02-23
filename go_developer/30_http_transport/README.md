# 30. net/http Transport

## 核心ポイント
**Body を読み切らないと接続が再利用されない。**

## Transport の内部

`http.Transport` は接続プールを管理する:

```
Transport {
    MaxIdleConns:        100    // 全体のアイドル接続数
    MaxIdleConnsPerHost: 2      // ホストごとのアイドル接続数（デフォルト2!）
    IdleConnTimeout:     90s    // アイドル接続のタイムアウト
}
```

## 接続再利用の条件

レスポンスの Body を**完全に読み切って Close**しないと、接続がプールに戻らない。

```go
// 悪い例: Body を読み切らない
resp, _ := client.Get(url)
resp.Body.Close()  // Body にデータが残っている → 接続が再利用されない

// 良い例: Body を完全に読み切る
resp, _ := client.Get(url)
io.Copy(io.Discard, resp.Body)  // 全て読み切る
resp.Body.Close()                // → 接続がプールに戻る
```

## タイムアウト設定

```
Client.Timeout          ← リクエスト全体のタイムアウト
├─ Transport.DialTimeout     ← TCP 接続
├─ Transport.TLSHandshakeTimeout ← TLS ハンドシェイク
├─ Transport.ResponseHeaderTimeout ← レスポンスヘッダー受信
└─ Body の読み取り          ← Client.Timeout に含まれる
```

## 実行方法

```bash
go run ./30_http_transport/
```
