# Hackathon24spring02 Realtime Server(WebSocket)
WebSocketを用いたオンラインゲームのリアルタイム通信サーバーです。
イベント駆動型アーキテクチャに書き直しました。

## 概要
```
 .
├──  domain
│  └──  player.go 受送信するデータの型の定義
├──  go.mod
├──  go.sum
├──  main.go IPv4とIPv6の両方でリッスン(/ws)
└──  server
   ├──  client.go 接続しているクライアントの管理
   ├──  handlers.go データを受信するロジック
   ├──  item-generator.go LLMを用いたアイテム生成
   ├──  message-broadcaster.go データを送信するロジック
   ├──  read-write-loops.go データの受送信を実際に行う処理
   └──  state.go ゲームの状態を操作する処理
```

`read-write-loops.go`(readloop) → `handlers.go` → `message-broadcaster.go` → `read-write-loops.go`(writeloop)

