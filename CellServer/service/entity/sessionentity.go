package entity

func CreateSessionJoinSystem(session interface{}) *SessionJoinEntity {
	return &SessionJoinEntity{
		SessionData: session,
	}
}

func CreateSessionLeaveSystem(session interface{}) *SessionLeaveEntity {
	return &SessionLeaveEntity{
		SessionData: session,
	}
}

//进入房间
type SessionJoinEntity struct {
	SessionData interface{}
}

//离开房间
type SessionLeaveEntity struct {
	SessionData interface{}
}
