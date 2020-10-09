package libray

////////// to client /////////////
// login
type LoginReq struct {
	Id   string `json:"id"`
	Cmd  string `json:"cmd"` // login
	Uuid uint64 `json:"uuid"`
}

// heartbeat
type HeartBeatReq struct {
	Id  string `json:"id"`
	Cmd string `json:"cmd"` // heartbeat
}

// push
type PushClientReq struct {
	Id  string `json:"id"`
	Cmd string `json:"cmd"`
	Msg string `json:"msg"`
}

type Response struct {
	Id   string `json:"id"`
	Cmd  string `json:"cmd"`
	Data struct {
		Code   uint32 `json:"code"`
		Result string `json:"result"`
	}
}

// -------- route -------
type RouteReq struct {
	Id   string
	Uuid uint64
	Msg  string
}

type RouteRsp struct {
	Seq string
}

// - - - - server - - - -
// push msg
type PushReq struct {
	Id   string
	Uuid uint64
	Msg  string
}

type PushRsp struct {
	Seq string
}

// query isOnline
type QueryReq struct {
	Uuid uint64
}

type QueryRsp struct {
	IsOnline bool
}

// push transfer
type TransferData struct {
	Seq  string `json:"seq"`
	Uuid uint64 `json:"uuid"`
	Msg  string `json:"msg"`
}
