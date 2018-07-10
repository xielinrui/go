package server
import (
	"game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"sync"
	"time"
	"fmt"
)

type MsgRegistry struct {
	registry           map[int32]func(msg *server.ClientMsg, sess *server.Session) []byte
	unLoginMsgRegistry map[int32]bool
	mu                 sync.RWMutex
}

func (registry *MsgRegistry) RegisterMsg(msgId int32, f func(msg *server.ClientMsg, sess *server.Session) []byte) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.registry[msgId] = f
}

func (registry *MsgRegistry) isUnLoginMsg(msgId int32) bool {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	_, ok := registry.unLoginMsgRegistry[msgId]
	return ok
}

func (registry *MsgRegistry) getHandler(msgId int32) func(msg *server.ClientMsg, sess *server.Session) []byte {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	return registry.registry[msgId]
}

func (registry *MsgRegistry) DispatchMsg(msg *server.ClientMsg, sess *server.Session) []byte {
	start := time.Now()
	fmt.Println(msg.MsgId)
	defer func() {
		elapseTime := time.Since(start)
		if elapseTime.Seconds() > 0.1 {
			p := user.GetPlayer(sess.Data)
			if p != nil {
				user.SaveSlowMsg(p.User.UserId, msg.MsgId, start, elapseTime.String())
			}
		}
	}()

	if !registry.isUnLoginMsg(msg.MsgId) {
		if user.GetPlayer(sess.Data) == nil {
			return nil
		}
	}
	if msg.MsgId != 2{

		pp := user.GetPlayer(sess.Data)
		if pp != nil {

		}
	}
	f := registry.getHandler(msg.MsgId)
	if f == nil {
		glog.Error("msgId:  has no handler")
		return nil
	}
	return f(msg, sess)
}

func (registry *MsgRegistry) RegisterHandlers(r *mux.Router) {
	registerHandlers(r)
}

var registry *MsgRegistry

func init() {
	registry = &MsgRegistry{
		registry:           make(map[int32]func(msg *server.ClientMsg, sess *server.Session) []byte),
		unLoginMsgRegistry: make(map[int32]bool),
		mu:                 sync.RWMutex{},
	}
}

func GetMsgRegistry() *MsgRegistry {
	return registry
}
