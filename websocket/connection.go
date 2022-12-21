package websocket

import (
	"github.com/gobwas/ws"
	"im"
	"net"
)

type Frame struct {
	raw ws.Frame
}

func (f *Frame) SetOpCode(code im.OpCode) {
	f.raw.Header.OpCode = ws.OpCode(code)
}

func (f *Frame) GetOpCode() im.OpCode {
	return im.OpCode(f.raw.Header.OpCode)
}

func (f *Frame) SetPayload(payload []byte) {
	f.raw.Payload = payload
}

func (f *Frame) GetPayload() []byte {
	// 表示是否要对数据载荷进行掩码操作。从客户端向服务端发送数据时，需要对数据进行掩码操作；从服务端向客户端发送数据时，不需要对数据进行掩码操作
	if f.raw.Header.Masked {
		ws.Cipher(f.raw.Payload, f.raw.Header.Mask, 0)
	}
	f.raw.Header.Masked = false
	return f.raw.Payload
}

type WsConn struct {
	net.Conn
}

func NewConn(conn net.Conn) *WsConn {
	return &WsConn{
		Conn: conn,
	}
}

func (c *WsConn) ReadFrame() (im.Frame, error) {
	f, err := ws.ReadFrame(c.Conn)
	if err != nil {
		return nil, err
	}
	return &Frame{raw: f}, nil
}

func (c *WsConn) WriteFrame(code im.OpCode, payload []byte) error {
	f := ws.NewFrame(ws.OpCode(code), true, payload)
	return ws.WriteFrame(c.Conn, f)
}

func (c *WsConn) Flush() error {
	return nil
}
