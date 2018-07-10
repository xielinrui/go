package server

import (
	"code.google.com/p/go.net/websocket"
	"sync"
	"github.com/golang/glog"
	"util"
	"strings"
	"fmt"
	"encoding/json"
	"time"
)

type ClientMsg struct {
	MsgId   int32       `json:"msgId"`
	MsgBody interface{} `json:"msgBody"`
}
type Session struct {
	conn      *websocket.Conn
	IP        string
	mq        chan *ClientMsg
	Data      interface{}
	LoggedIn  bool
	exitChan  chan bool
	cleanOnce sync.Once
	kickOnce  sync.Once
	OnLogout  func()
}
func newSess(conn *websocket.Conn) *Session {
	sess := &Session{}
	sess.conn = conn
	return sess
}
func (s *Session) logout() {
	if s.LoggedIn && s.OnLogout != nil {
		s.OnLogout()
	}
	s.LoggedIn = false
}
func (s *Session) Run(dispatcher MsgDispatcher) {
	defer util.PrintPanicStack()
	s.IP = s.conn.Request().Header.Get("X-Real-Ip")
	if len(s.IP) == 0 {
		s.IP = strings.Split(s.conn.Request().RemoteAddr, ":")[0]
	}

	s.mq = make(chan *ClientMsg, 100)
	s.exitChan = make(chan bool)
	GetServerInstance().waitGroup.Add(1)
	defer func() {
		util.PrintPanicStack()
		s.cleanSess()
		s.logout()
		GetServerInstance().waitGroup.Done()
	}()

	for {
		select {
		case <-GetServerInstance().stopChan:
			return
		case msg, ok := <-s.mq:
			if !ok {
				return
			}

			res := dispatcher.DispatchMsg(msg, s)
			if res != nil {
				if msg.MsgId != 2{
				}
				s.SendToClient(res)
			}
		case <-s.exitChan:
			return
		}
	}
}
func (s *Session) SendToClient(msg []byte) {
	ret := true

	defer func() {
		if r := recover(); r != nil {
			ret = false
			clientMsg := &ClientMsg{}
			err := json.Unmarshal(msg, clientMsg)
			if err != nil {
				glog.Error("SendToClient json.Unmarshal err:", err)
			}
		}
	}()

	if msg != nil {
		if glog.V(2) {
			clientMsg := &ClientMsg{}
			err := json.Unmarshal(msg, clientMsg)
			if err != nil {
				glog.Error("====>SendToClient unmarshal msg failed:", err)
			}
		}
		fmt.Println()
		fmt.Println(s.IP)
		s.conn.SetWriteDeadline(time.Now().Add(time.Second))
		sendMsg := string(msg)
		err := websocket.Message.Send(s.conn, sendMsg)
		if err != nil {
			s.cleanSess()
			return
		}
		s.conn.SetWriteDeadline(time.Time{})
	}
}
func (s *Session) cleanSess() {
	s.cleanOnce.Do(func() {
		s.conn.Close()
		if s.mq != nil {
			close(s.mq)
		}
	})
}