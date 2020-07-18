package main

import (
	"CellServer/core"
	"CellServer/net"
	"CellServer/service"
	"CellServer/service/command"
	"CellServer/service/entity"
	"CellServer/service/system"
	"flag"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"time"
	//_ "net/http/pprof"
)

func initialize() *chi.Mux {

	//创建网关
	gate = net.CreateGate()

	//注册系统
	core.GetWord().SystemMgr.AddSystem(
		core.CreateSystem(service.EntityTypeNetMessage, system.NetMessageSystem))
	core.GetWord().SystemMgr.AddSystem(
		core.CreateSystem(service.EntityTypeRoom, system.RoomSystem))
	core.GetWord().SystemMgr.AddSystem(
		core.CreateSystem(service.EntityTypeSessionJoin, system.SessionJoinSystem))
	core.GetWord().SystemMgr.AddSystem(
		core.CreateSystem(service.EntityTypeSessionLeave, system.SessionLeaveSystem))
	core.GetWord().SystemMgr.AddSystem(
		core.CreateSystem(service.EntityTypeRoomPacket, system.RoomPacketSystem))

	//注册路由
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome."))
	})

	r.Get("/create", func(w http.ResponseWriter, r *http.Request) {
		//创建房间
		room := entity.CreateRoomEntity()
		cmd := &command.AddEntityCommand{
			Type: service.EntityTypeRoom,
			Data: room,
		}
		core.GetWord().Commit(cmd)
		w.Write([]byte("RoomID: " + room.ID))
		log.Println("Create Room: '", room.ID, "'")
	})

	r.Get("/{rid}", func(w http.ResponseWriter, r *http.Request) {
		rid := chi.URLParam(r, "rid")
		net.PeerHandler(w, r, gate, rid)
	})

	return r
}

var addr = flag.String("addr", ":36000", "http service address")
var gate *net.Gate

func main() {

	//解析命令行
	flag.Parse()
	log.Println(*addr, "Server running...")

	//初始化
	r := initialize()

	//启动服务
	go core.GetWord().Start()

	//启动网关
	go gate.Run()

	//性能测试
	/*
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	*/

	//启动网络
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
