package server
import (
	"github.com/gorilla/mux"
	"net/http"
	"code.google.com/p/go.net/websocket"
	"github.com/golang/glog"
	"fmt"
	"os"
	"sync"
	"time"
	"encoding/json"
	"sync/atomic"
)
type MsgDispatcher interface {
	RegisterHandlers(r *mux.Router)
	DispatchMsg(msg *ClientMsg, sess *Session) []byte
}
type GameServer struct {
	dispatcher    MsgDispatcher
	sigChan       chan os.Signal
	waitGroup     *sync.WaitGroup
	stopChan      chan bool
	stopOnce      sync.Once
	refuseService int32
}
var bindHost string
var s *GameServer
func (s *GameServer) IsRefuseService() bool {
	return atomic.AddInt32(&s.refuseService, 0) > 0
}
func GetServerInstance() *GameServer {
	return s
}
func (s *GameServer) StartServer(dispatcher MsgDispatcher) {
	s.dispatcher = dispatcher
	r := mux.NewRouter()
	http.Handle("/", r)
	http.Handle("/ws/", websocket.Server{Handler: s.handleClient, Handshake: nil})
	s.dispatcher.RegisterHandlers(r)
	glog.Fatal(http.ListenAndServe(fmt.Sprintf("%v", bindHost), nil))
}

func (s *GameServer) handleClient(conn *websocket.Conn) {
	if s.IsRefuseService() {
		conn.Close()
		return
	}
	sess := newSess(conn)
	go sess.Run(s.dispatcher)
	defer sess.cleanSess()
	for {
		var data []byte
		conn.SetReadDeadline(time.Now().Add(time.Minute * 10))
		err := websocket.Message.Receive(conn, &data)
		if err != nil {
			break
		}
		conn.SetReadDeadline(time.Time{})
		msg := &ClientMsg{}
		errMsg := json.Unmarshal(data, msg)
		if errMsg != nil {
			glog.Error("unmarshal client msg failed!")
			break
		}
		sess.mq <- msg
	}
}
