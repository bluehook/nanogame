package core

import (
	"fmt"
	"sync"
	"time"
)

var instanceWord *Word
var onceWord sync.Once

func GetWord() *Word {
	onceWord.Do(func() {
		instanceWord = &Word{
			EntityMgr:     CreateEntityManager(),
			SystemMgr:     CreateSystemManager(),
			HeartTime:     (int64)(time.Millisecond * 50),
			commandStream: make(chan Command, 1024),
			updateWait:    &sync.WaitGroup{},
		}
	})
	return instanceWord
}

//服务框架
type Word struct {
	//所有实体
	EntityMgr *EntityManager
	//所有系统
	SystemMgr *SystemManager
	//当前时间
	CurTime int64
	//每帧实际消耗时间
	DeltaTime int64
	//系统更新间隔时间
	HeartTime int64
	//更新并发等待
	updateWait *sync.WaitGroup
	//全局命令管道
	commandStream chan Command
}

//开始服务
func (word *Word) Start() {
	var elapse int64
	for {
		word.CurTime = time.Now().UnixNano()

		//心跳到来时进行系统更新
		if elapse > word.HeartTime {
			//word.DebugStatus()
			word.update()
			elapse = 0
		}

		//其他时间处理命令
		word.processCommands()
		time.Sleep(time.Nanosecond)

		word.DeltaTime = time.Now().UnixNano() - word.CurTime
		elapse += word.DeltaTime
	}
}

func (word *Word) Commit(command Command) {
	word.commandStream <- command
}

//处理系统更新
func (word *Word) update() {
	word.EntityMgr.Map(func(t int, c *Chunck) {
		word.updateWait.Add(1)
		go word.SystemMgr.Update(t, c, word.updateWait)
	})
	word.updateWait.Wait()
}

//处理延迟命令
func (word *Word) processCommands() {
	for {
		select {
		case c := <-word.commandStream:
			c.Execute()
		default:
			return
		}
	}
}

//调试统计信息
func (word *Word) DebugStatus() {
	fmt.Println("实体类型个数:", word.EntityMgr.Len())
	word.EntityMgr.Map(func(t int, c *Chunck) {
		fmt.Println("--实体类型:", t)
		fmt.Println("----实体数:", c.Len())
	})

	fmt.Println("系统类型个数:", word.SystemMgr.Len())
	fmt.Println("----------------------------------")
}
