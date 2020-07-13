package net

import (
	"sync/atomic"
	"time"
)

const (
	WaitTimeOut       = 10 * time.Second          //客户端网络连接超时时间
	MaxTimeOutCount   = 3                         //允许超时次数
	MaxMessageSize    = 256                       //客户端数据包最大字节数
	MaxTokenCount     = 15                        //最大每秒通讯令牌数
	RoomEmptyKeepTime = int64(time.Second * 1800) //房间空置关闭时间
)

var Connections = newConnectionService()

type connectionService struct {
	count int64
	sid   int64
}

func newConnectionService() *connectionService {
	return &connectionService{sid: 0}
}

func (c *connectionService) Increment() {
	atomic.AddInt64(&c.count, 1)
}

func (c *connectionService) Decrement() {
	atomic.AddInt64(&c.count, -1)
}

func (c *connectionService) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

func (c *connectionService) Reset() {
	atomic.StoreInt64(&c.count, 0)
	atomic.StoreInt64(&c.sid, 0)
}

func (c *connectionService) SessionID() int64 {
	return atomic.AddInt64(&c.sid, 1)
}
