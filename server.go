package im

import (
	"context"
	"net"
	"time"
)

// OpCode OpCode
type OpCode byte

const (
	DefaultReadWait         = time.Minute * 3
	DefaultWriteWait        = time.Second * 10
	DefaultLoginWait        = time.Second * 10
	DefaultHeartbeat        = time.Second * 55
	OpContinuation   OpCode = 0x0
	OpText           OpCode = 0x1
	OpBinary         OpCode = 0x2
	OpClose          OpCode = 0x8
	OpPing           OpCode = 0x9
	OpPong           OpCode = 0xa
)

type Server interface {
	SetAcceptor(Acceptor)
	SetMessageListener(MessageListener)
	SetStateListener(StateListener)
	SetReadWait(time.Duration)
	SetChannelMap(ChannelMap)

	Start() error
	Push(string, []byte) error
	Shutdown(context.Context) error
}

// Conn Connection
type Conn interface {
	net.Conn
	ReadFrame() (Frame, error)
	WriteFrame(OpCode, []byte) error
	Flush() error
}

// Frame 通过抽象一个Frame接口来解决底层封包与拆包问题。
type Frame interface {
	SetOpCode(OpCode)
	GetOpCode() OpCode
	SetPayload([]byte)
	GetPayload() []byte
}

type Acceptor interface {
	// Accept 返回一个握手完成的Channel对象或者一个error。
	Accept(Conn, time.Duration) (string, error)
}

type StateListener interface {
	Disconnect(string) error
}

// Agent is interface of client side
type Agent interface {
	ID() string
	Push([]byte) error
}

// MessageListener 监听消息
type MessageListener interface {
	// 收到消息回调
	Receive(Agent, []byte)
}

type DialerContext struct {
	Id      string
	Name    string
	Address string
	Timeout time.Duration
}

// Dialer Dialer
type Dialer interface {
	DialAndHandshake(DialerContext) (net.Conn, error)
}

// Client is interface of client side
type Client interface {
	ID() string
	Name() string
	// connect to server
	Connect(string) error
	// SetDialer 设置拨号处理器
	SetDialer(Dialer)
	Send([]byte) error
	Read() (Frame, error)
	// Close 关闭
	Close()
}
