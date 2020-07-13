package net

import (
	"CellServer/core"
	"CellServer/service"
	"CellServer/service/command"
	"CellServer/service/entity"
	"log"
)

func CreateGate() *Gate {
	return &Gate{
		broadcastStream: make(chan []byte),
		registerStrem:   make(chan *Session),
		unregisterStrem: make(chan *Session),
		sessions:        make(map[*Session]bool),
	}
}

type Gate struct {
	sessions        map[*Session]bool //登录用户
	broadcastStream chan []byte       //广播消息管道
	registerStrem   chan *Session     //登录请求管道
	unregisterStrem chan *Session     //退出请求管道
}

func (g *Gate) Run() {
	for {
		select {
		case session := <-g.registerStrem:
			g.sessions[session] = true
			//进入系统处理
			cmd := &command.AddEntityCommand{
				Type: service.EntityTypeSessionJoin,
				Data: entity.CreateSessionJoinSystem(session),
			}
			core.GetWord().Commit(cmd)
			log.Println("Client Join: [", session.id, "]")
		case session := <-g.unregisterStrem:
			if _, ok := g.sessions[session]; ok {
				//离开系统处理
				cmd := &command.AddEntityCommand{
					Type: service.EntityTypeSessionLeave,
					Data: entity.CreateSessionLeaveSystem(session),
				}
				core.GetWord().Commit(cmd)
				delete(g.sessions, session)
				close(session.send)
				log.Println("Client Quit: [", session.id, "]")
			}
		case message := <-g.broadcastStream:
			for session := range g.sessions {
				select {
				case session.send <- message:
				default:
					close(session.send)
					delete(g.sessions, session)
				}
			}
		}
	}
}
