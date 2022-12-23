package tcp

import (
	"im"
	"im/wire/endian"
	"io"
	"net"
)

type Frame struct {
	OpCode  im.OpCode
	Payload []byte
}

// SetOpCode SetOpCode
func (f *Frame) SetOpCode(code im.OpCode) {
	f.OpCode = code
}

// GetOpCode GetOpCode
func (f *Frame) GetOpCode() im.OpCode {
	return f.OpCode
}

// SetPayload SetPayload
func (f *Frame) SetPayload(payload []byte) {
	f.Payload = payload
}

// GetPayload GetPayload
func (f *Frame) GetPayload() []byte {
	return f.Payload
}

type TcpConn struct {
	net.Conn
}

// NewConn NewConn
func NewConn(conn net.Conn) *TcpConn {
	return &TcpConn{
		Conn: conn,
	}
}

func (c *TcpConn) ReadFrame() (im.Frame, error) {
	// 从 reader 中读取一个 uint8
	opcode, err := endian.ReadUint8(c.Conn)
	if err != nil {
		return nil, err
	}
	// 从 reader 中读取一个 []byte
	payload, err := endian.ReadBytes(c.Conn)
	if err != nil {
		return nil, err
	}
	return &Frame{
		OpCode:  im.OpCode(opcode),
		Payload: payload,
	}, nil
}

func (c *TcpConn) WriteFrame(code im.OpCode, payload []byte) error {
	return WriteFrame(c.Conn, code, payload)
}

// Flush Flush
func (c *TcpConn) Flush() error {
	return nil
}

// WriteFrame write a frame to w
func WriteFrame(w io.Writer, code im.OpCode, payload []byte) error {
	if err := endian.WriteUint8(w, uint8(code)); err != nil {
		return err
	}
	if err := endian.WriteBytes(w, payload); err != nil {
		return err
	}
	return nil
}
