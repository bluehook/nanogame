package net

import (
	"CellServer/core"
	"CellServer/message"
	"CellServer/service"
	"CellServer/service/command"
	"CellServer/service/entity"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  MaxMessageSize,
	WriteBufferSize: MaxMessageSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func PeerHandler(w http.ResponseWriter, r *http.Request, gate *Gate, rid string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	session := CreateSession(gate, conn)
	session.Set("RoomID", rid)
	session.gate.registerStrem <- session

	go session.writePump()
	go session.readPump()
}

func CreateSession(gate *Gate, conn *websocket.Conn) *Session {
	s := &Session{
		gate: gate,
		conn: conn,
		id:   Connections.SessionID(),
		data: make(map[string]interface{}),
		send: make(chan []byte, 16),
	}
	s.Set("timeout", 0)
	s.Set("lastTime", time.Now().Unix())
	return s
}

type Session struct {
	sync.RWMutex
	gate *Gate
	id   int64
	conn *websocket.Conn
	data map[string]interface{}
	send chan []byte
}

func (s *Session) SendStream(msg []byte) {
	defer func() {
		if recover() != nil {
		}
	}()
	s.send <- msg
}

func (s *Session) KickSelf() {
	defer func() {
		s.gate.unregisterStrem <- s
		s.conn.Close()
	}()
	return
}

func (s *Session) readPump() {
	defer func() {
		s.gate.unregisterStrem <- s
		s.conn.Close()
	}()
	s.conn.SetReadLimit(MaxMessageSize)
	tokenCount := MaxTokenCount
	limitTime := int64(time.Second / MaxTokenCount)
	var curt, elapse int64
	for {
		curt = time.Now().UnixNano()
		_, msg, err := s.conn.ReadMessage()
		if err != nil {
			return
		}
		//超时次数清零
		s.Set("timeout", 0)

		//流量控制
		if tokenCount > 0 {
			//快速检测是否是心跳包
			if int8(msg[0]) != message.MsgNull {
				cmd := &command.AddEntityCommand{
					Type: service.EntityTypeNetMessage,
					Data: entity.CreateNetMessageEntity(msg, s),
				}
				core.GetWord().Commit(cmd)
				tokenCount -= 1
			}
		}
		elapse = time.Now().UnixNano() - curt
		//计算令牌
		tokenCount += int(elapse / limitTime)
		if tokenCount > MaxTokenCount {
			tokenCount = MaxTokenCount
		}
	}
}

func (s *Session) writePump() {
	//超时触发器
	ticker := time.NewTicker(WaitTimeOut)
	defer func() {
		ticker.Stop()
		s.conn.Close()
	}()
	for {
		select {
		case message, ok := <-s.send:
			if !ok {
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := s.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			count := s.Int("timeout")
			if count < MaxTimeOutCount {
				count += 1
				s.Set("timeout", count)
			} else {
				//离线处理
				log.Println("Session Ping Timeout Leave: [", s.ID(), "]")
				return
			}
		}
	}
}

func (s *Session) Conn() *websocket.Conn {
	return s.conn
}

func (s *Session) Gate() *Gate {
	return s.gate
}

func (s *Session) ID() int64 {
	return s.id
}

func (s *Session) Remove(key string) {
	s.Lock()
	defer s.Unlock()

	delete(s.data, key)
}

func (s *Session) Set(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()

	s.data[key] = value
}

func (s *Session) HasKey(key string) bool {
	s.RLock()
	defer s.RUnlock()

	_, has := s.data[key]
	return has
}

func (s *Session) Int(key string) int {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Int8(key string) int8 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int8)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Int16(key string) int16 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int16)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Int32(key string) int32 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int32)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Int64(key string) int64 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int64)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint(key string) uint {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint8(key string) uint8 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint8)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint16(key string) uint16 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint16)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint32(key string) uint32 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint32)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint64(key string) uint64 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint64)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Float32(key string) float32 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(float32)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Float64(key string) float64 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(float64)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) String(key string) string {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return ""
	}

	value, ok := v.(string)
	if !ok {
		return ""
	}
	return value
}

func (s *Session) Value(key string) interface{} {
	s.RLock()
	defer s.RUnlock()

	return s.data[key]
}

func (s *Session) State() map[string]interface{} {
	s.RLock()
	defer s.RUnlock()

	return s.data
}

func (s *Session) Restore(data map[string]interface{}) {
	s.Lock()
	defer s.Unlock()

	s.data = data
}

func (s *Session) Clear() {
	s.Lock()
	defer s.Unlock()

	s.data = map[string]interface{}{}
}
