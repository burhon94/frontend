package server

import (
	"github.com/burhon94/alifMux/pkg/mux"
	"github.com/burhon94/frontend/core/authSvc"
	"github.com/burhon94/frontend/core/fileSvc"
	"github.com/burhon94/frontend/core/postSvc"
	"net/http"
)

type ServerStruct struct {
	router *mux.ExactMux
	authClient *authSvc.AuthClient
	fileClient *fileSvc.FileClient
	postClient *postSvc.PostClient
}

func NewServerStruct(router *mux.ExactMux, authClient *authSvc.AuthClient, fileClient *fileSvc.FileClient, postClient *postSvc.PostClient) *ServerStruct {
	return &ServerStruct{router: router, authClient: authClient, fileClient: fileClient, postClient: postClient}
}


func (server *ServerStruct) Start()  {
	server.InitRoutes()
}

func (server *ServerStruct) ServeHTTP(writer http.ResponseWriter, request *http.Request)  {
	server.router.ServeHTTP(writer, request)
}
