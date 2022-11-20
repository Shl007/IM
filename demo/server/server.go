package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	once    sync.Once
	id      string
	address string
	sync.Mutex
	// a list of sessions
	users map[string]net.Conn
}

type ServerOptions struct {
	writewait time.Duration //写超时时间
	readwait  time.Duration //读超时时间
}

func NewServer(id, address string) *Server {
	return newServer(id, address)
}

func newServer(id, address string) *Server {
	return &Server{
		id:      id,
		address: address,
		users:   make(map[string]net.Conn, 100),
	}
}

func (s *Server) Start() error {
	// Multiplexing 多路复用
	mux := http.NewServeMux()
	log := logrus.WithFields(logrus.Fields{
		"module": "Server",
		"listen": s.address,
		"id":     s.id,
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// step1 升级长连接
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			conn.Close()
			return
		}
		// step2 拿到user信息
		user := r.URL.Query().Get("user")
		if user == "" {
			conn.Close()
			return
		}
		// step3 添加到list
		old, ok := s.addUser(user, conn)
		if ok {
			old.Close()
		}
		log.Infof("user %s in", user)

		go func(user string, conn net.Conn) {
			// step4 read
			err := s.readloop(user, conn)
			if err != nil {
				log.Error(err)
			}
			conn.Close()
			// step5
			s.delUser(user)
			log.Infof("connection of %s closed", user)
		}(user, conn)
	})
	log.Infof("started")
	return http.ListenAndServe(s.address, mux)
}

func (s *Server) addUser(user string, conn net.Conn) (net.Conn, bool) {
	s.Lock()
	defer s.Unlock()
	old, ok := s.users[user]
	s.users[user] = conn
	return old, ok
}

func (s *Server) delUser(user string) {
	s.Lock()
	defer s.Unlock()
	delete(s.users, user)
}

func (s *Server) Shutdown() {
	s.once.Do(func() {
		s.Lock()
		defer s.Unlock()
		for _, conn := range s.users {
			conn.Close()
		}
	})
}

func (s *Server) readloop(user string, conn net.Conn) interface{} {
	for {
		// 从TCP缓冲区读取一帧的消息
		frame, err := ws.ReadFrame(conn)
		if err != nil {
			return err
		}
		if frame.Header.OpCode == ws.OpClose {
			return errors.New("remote side close the conn")
		}
		if frame.Header.Masked {
			/**
			websocket协议规定客户端发送数据时，必须使用一个随机的mask值对消息体做一次编码，
			因此在服务端就需要解码，否则内容是乱的。
			*/
			ws.Cipher(frame.Payload, frame.Header.Mask, 0)
		}
		if frame.Header.OpCode == ws.OpText {
			go s.handel(user, string(frame.Payload))
		}
	}
}

func (s *Server) handel(user string, message string) {
	logrus.Infof("recv message %s from %s", message, user)
	s.Lock()
	// map不安全，需要加锁
	defer s.Unlock()
	broadcast := fmt.Sprintf("%s -- from %s", message, user)
	for u, conn := range s.users {
		if u == user {
			continue
		}
		logrus.Infof("send to % s : %s", u, broadcast)
		err := s.writeText(conn, broadcast)
		if err != nil {
			logrus.Errorf("write to %s failed, error: %v", user, err)
		}
	}
}

func (s *Server) writeText(conn net.Conn, message string) error {
	// 创建文本帧数据
	f := ws.NewTextFrame([]byte(message))
	return ws.WriteFrame(conn, f)
}

const (
	CommandPing = 100
	CommandPong = 101
)

func (s *Server) handleBinary(user string, message []byte) {
	logrus.Infof("recv message %v from %s", message, user)
	s.Lock()
	defer s.Unlock()
	// handle ping request
	i := 0
	command := binary.BigEndian.Uint16(message[i : i+2])
	i += 2
	payloadLen := binary.BigEndian.Uint32(message[i : i+4])
	logrus.Infof("command: %v payloadLen: %v", command, payloadLen)
	if command == CommandPing {
		u := s.users[user]
		// return pong
		err := wsutil.WriteServerBinary(u, []byte{0, CommandPong, 0, 0, 0, 0})
		if err != nil {
			logrus.Errorf("write to %s failed, error: %v", user, err)
		}
	}
}
