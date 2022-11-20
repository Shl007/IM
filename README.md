# IM
IM system written in GO

# demo

有锁无chan的websocket demo

会群发消息给除自己其它在线用户

## 1.1环境准备

golang

## 1.2启动服务

`$ go run main.go chat`

## 1.3测试

随便打开一个websocket测试工具，例如（[WebSocket在线测试工具 (wstool.js.org)](http://wstool.js.org/)）

```
ws://localhost:8000?user=userA
ws://localhost:8000?user=userB
```

