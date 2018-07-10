package server

import (
	"github.com/gorilla/mux"
)

func registerHandlers(r *mux.Router) {
	registerHttpHandlers(r)
	registerUserHandlers()
}

func registerHttpHandlers(r *mux.Router) {
	r.HandleFunc("/v1/userInfo", )//在逗号后面绑定handler
	r.HandleFunc("/v1/userCharge",)//在逗号后面绑定handler
	r.HandleFunc("/v1/cardChgLog", )//在逗号后面绑定handler
}

func registerUserHandlers() {
	GetMsgRegistry().RegisterMsg(int32(MessageId_LOGIN),)//在逗号后面绑定handler
	GetMsgRegistry().RegisterMsg(int32(MessageId_CREATE_ROOM), )//在逗号后面绑定handler
	GetMsgRegistry().RegisterMsg(int32(MessageId_ENTER_ROOM), )//在逗号后面绑定handler
	GetMsgRegistry().RegisterMsg(int32(MessageId_BACK_ROOM),)//在逗号后面绑定handler
}
