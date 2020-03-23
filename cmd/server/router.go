package server

import (
	"github.com/burhon94/frontend/core/middleware/authenticated"
	"github.com/burhon94/frontend/core/middleware/jwt"
	"github.com/burhon94/frontend/core/middleware/logger"
	"github.com/burhon94/frontend/core/middleware/unauthenticated"
	jwt2 "github.com/burhon94/jsonwebtoken/pkg/cmd"
	"reflect"
)

var (
	Root = "/"
	Admin = "/admin"
)

func (server *ServerStruct) InitRoutes()  {
	jwtMW := jwt.JWT(jwt.SourceCookie, true, Root, reflect.TypeOf((*Payload)(nil)).Elem(), jwt2.Secret("alifkey"))
	authMW := authenticated.Authenticated(jwt.IsContextNonEmpty, true, Root)
	unAuthMW := unauthenticated.Unauthenticated(jwt.IsContextNonEmpty, true, Admin)
	//index page
	server.router.GET("/", server.handlerIndexPage(),unAuthMW, jwtMW, logger.Logger("HTTP"))
	server.router.POST("/", server.handlerIndexPage(), unAuthMW, jwtMW, logger.Logger("HTTP"))
	//admin page
	server.router.GET("/admin", server.handlerIndexPageAdmin(), authMW, jwtMW, logger.Logger("HTTP"))
	server.router.POST("/admin", server.handlerIndexPageAdmin(), authMW, jwtMW, logger.Logger("HTTP"))
	server.router.GET("/post/0", server.handlerPostNewPage(), authMW, jwtMW, logger.Logger("HTTP"))
	server.router.POST("/postNew", server.handlerPostNew(), authMW, jwtMW, logger.Logger("HTTP"))
	server.router.POST("/postFiles", server.handlerPostFiles(), authMW, jwtMW, logger.Logger("HTTP"))
	//get services status
	server.router.GET("/authSvc", server.handlerHealthAuth(), logger.Logger("HTTP"))
	server.router.GET("/fileSvc", server.handlerHealthFile(), logger.Logger("HTTP"))
	server.router.GET("/postSvc", server.handlerHealthPost(), logger.Logger("HTTP"))
	//goto verify page
	server.router.GET("/api/verify", server.handlerVerifyRedir(), logger.Logger("VerifyRedirect"))
	//get post
	server.router.GET("/movies", server.handlerPostsMovies(), logger.Logger("PostsMoies"))
	server.router.GET("/games", server.handlerPostsGames(), logger.Logger("PostsGames"))
	server.router.GET("/softs", server.handlerPostsSofts(), logger.Logger("PostsSofts"))
	server.router.GET("/musics", server.handlerPostsMusics(), logger.Logger("PostsMusics"))
}
