/**
 * Copyright 2018 hydra Author. All rights reserved.
 * Author: Chuck1024
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"hydra/common"
	"net/url"
	"strconv"
	"sync/atomic"
	"time"
)

var addr = flag.String("addr", "127.0.0.1:10240", "http service address")
var uuid = flag.Uint64("uuid", 10240, "uuid")

func main() {
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/hydra"}
	dialer := &websocket.Dialer{}

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go heartbeat(conn)
	go handle(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		rsp := &common.Response{}

		json.Unmarshal(message, rsp)
		if rsp.Cmd == common.PushCmd {
			resp := &common.PushClientReq{}
			json.Unmarshal(message, resp)
			fmt.Println("received: ", *resp)

			pr := &common.Response{
				Id:  resp.Id,
				Cmd: common.PushCmd,
			}
			pr.Data.Code = 200
			pr.Data.Result = "ok"

			prb, _ := json.Marshal(pr)
			conn.WriteMessage(websocket.BinaryMessage, prb)
		} else {
			fmt.Println("received: ", *rsp)
		}
	}
}

var globalSeq uint32

func nextSeq() string {
	return strconv.FormatUint(uint64(atomic.AddUint32(&globalSeq, 1)), 10)
}

func heartbeat(conn *websocket.Conn) {
	t := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-t.C:
			h := &common.HeartBeatReq{
				Id:  nextSeq(),
				Cmd: common.HeartbeatCmd,
			}

			hh, _ := json.Marshal(h)
			conn.WriteMessage(websocket.BinaryMessage, hh)
		}
	}
}

func handle(conn *websocket.Conn) {
	for {
		var input string
		fmt.Scan(&input)
		if len(input) > 0 {
			switch input {
			case common.LoginCmd:
				l := &common.LoginReq{
					Id:   nextSeq(),
					Cmd:  common.LoginCmd,
					Uuid: *uuid,
				}
				ll, _ := json.Marshal(l)
				conn.WriteMessage(websocket.BinaryMessage, ll)

			case common.HeartbeatCmd:
				h := &common.HeartBeatReq{
					Id:  nextSeq(),
					Cmd: common.HeartbeatCmd,
				}

				hh, _ := json.Marshal(h)
				conn.WriteMessage(websocket.BinaryMessage, hh)
			}
		}
	}
}
