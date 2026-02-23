# 31. http.Client Reuse

## 核心ポイント
**リクエストごとに http.Client を new するな！接続プールが無駄になる。**

## なぜバグなのか

```go
// 悪い例: リクエストごとに新しい Client
func fetch(url string) {
    client := &http.Client{}  // ← 毎回新しい Transport = 新しい接続プール
    resp, _ := client.Get(url)
    // ...
}
```

新しい Client = 新しい Transport = **接続プールが共有されない**。
毎リクエストで TCP + TLS ハンドシェイクが発生する。

## http.DefaultClient の罠

```go
http.Get(url)  // http.DefaultClient を使う
```

`DefaultClient` は**タイムアウトなし**。
プロダクションでは必ずタイムアウト付きの Client を使う。

## 正しいパターン

```go
// グローバルまたは DI で再利用
var httpClient = &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConnsPerHost: 10,
    },
}
```

## 実行方法

```bash
go run ./31_http_client_reuse/
```
