package main

import (
	"flag"
	"fmt"
	"github.com/burhon94/alifMux/pkg/mux"
	"github.com/burhon94/bdi/pkg/di"
	"github.com/burhon94/frontend/cmd/server"
	"github.com/burhon94/frontend/core/authSvc"
	"github.com/burhon94/frontend/core/fileSvc"
	"github.com/burhon94/frontend/core/postSvc"
	"log"
	"net"
	"net/http"
)
// -host 0.0.0.0 -port 4444 -authSvcUrl http://localhost:10000 -fileSvcUrl http://localhost:20000 -postSvcUrl http://localhost:7777
var (
	hostFlag   = flag.String("host", "", "Server host")
	portFlag   = flag.String("port", "", "Server host")
	authSvcUrl = flag.String("authSvcUrl", "", "Auth Service URL")
	fileSvcUrl = flag.String("fileSvcUrl", "", "File Service URL")
	postSvcUrl = flag.String("postSvcUrl", "", "File Service URL")
)

func main() {
	flag.Parse()
	addr := net.JoinHostPort(*hostFlag, *portFlag)
	start(addr, authSvc.Url(*authSvcUrl), fileSvc.Url(*fileSvcUrl), postSvc.Url(*postSvcUrl))
}

func start(addr string, authSvcUrl authSvc.Url, fileSvcUrl fileSvc.Url, postSvcUrl postSvc.Url) {
	container := di.NewContainer()
	err := container.Provide(
		server.NewServerStruct,
		mux.NewExactMux,
		func() authSvc.Url { return authSvcUrl },
		func() fileSvc.Url { return fileSvcUrl },
		func() postSvc.Url { return postSvcUrl },
		authSvc.NewAuthClient,
		fileSvc.NewFileClient,
		postSvc.NewPostClient,
		)
	if err != nil {
		panic(fmt.Sprintf("can't set provide: %v", err))
	}

	container.Start()
	var appServer *server.ServerStruct
	container.Component(&appServer)

	log.Printf("fronEnd start on: %s", addr)
	panic(http.ListenAndServe(addr, appServer))
}
