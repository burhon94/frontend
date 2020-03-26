package server

import (
	"context"
	"github.com/burhon94/frontend/core/postSvc"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

const (
	templatesPath     = "web/templates"
	multipartMaxBytes = 10 * 1024 * 1024 * 1024
)

func (server *ServerStruct) handlerIndexPage() http.HandlerFunc {
	var (
		tpl    *template.Template
		tplBad *template.Template
		err    error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "index.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badPostSvc.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		post := server.postClient.HealthPost()
		if post != true {
			err := tplBad.Execute(writer, struct{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Hour)
		posts, err := server.postClient.GetAllPosts(ctx)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = tpl.Execute(writer, posts)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (server *ServerStruct) handlerIndexPageAdmin() http.HandlerFunc {
	var (
		tpl    *template.Template
		tplBad *template.Template
		err    error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "admin", "index.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badPostSvc.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		post := server.postClient.HealthPost()
		if post != true {
			err := tplBad.Execute(writer, struct{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Hour)
		posts, err := server.postClient.GetAllPosts(ctx)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		//n := rand.Int() % len(posts)
		//fmt.Print(posts[n])

		err = tpl.Execute(writer, posts)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (server *ServerStruct) handlerPostNewPage() http.HandlerFunc {
	var (
		tpl           *template.Template
		tplBadPostSvc *template.Template
		tplBadFileSvc *template.Template
		err           error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "admin", "postPage.gohtml"),
		filepath.Join(templatesPath, "admin", "basePostPage.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBadPostSvc, err = template.ParseFiles(
		filepath.Join(templatesPath, "badPostSvc.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBadFileSvc, err = template.ParseFiles(
		filepath.Join(templatesPath, "badFileSvc.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		post := server.postClient.HealthPost()
		if post != true {
			err := tplBadPostSvc.Execute(writer, struct{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		status := server.fileClient.HealthFile()
		if status != true {
			err := tplBadFileSvc.Execute(writer, struct{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		err = tpl.Execute(writer, struct{}{})
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (server *ServerStruct) handlerPostNew() http.HandlerFunc {
	var (
		tplBadFile *template.Template
		tplBadPost *template.Template
		err        error
	)

	tplBadFile, err = template.ParseFiles(
		filepath.Join(templatesPath, "badFileSvc.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}
	tplBadPost, err = template.ParseFiles(
		filepath.Join(templatesPath, "badPostSvc.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		status := server.fileClient.HealthFile()
		if status != true {
			err := tplBadFile.Execute(writer, struct{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
		status = server.postClient.HealthPost()
		if status != true {
			err := tplBadPost.Execute(writer, struct{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
		err := request.ParseMultipartForm(multipartMaxBytes)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		postNew := postSvc.PostNew{}

		postNew.Title = request.PostFormValue("postTitle")
		if postNew.Title == "" {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		postNew.Category = request.PostFormValue("postCategory")
		if postNew.Category == "" {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		multipartForm:= request.MultipartForm
		file := multipartForm.File["files"]
		if len(file) != 2 {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Print(len(file))
		filesUrls, err := server.fileClient.UploadFiles(request)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		postNew.Poster = filesUrls[0]
		postNew.FileUrl = filesUrls[1]

		ctx, _ := context.WithTimeout(request.Context(), time.Hour)
		err = server.postClient.NewPost(ctx, postNew.Title, postNew.Poster, postNew.Category, postNew.FileUrl)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
	}
}

func (server *ServerStruct) handlerHealthAuth() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		healthAuth := server.authClient.HealthAuthCore()
		if healthAuth != true {
			log.Print("Verify Core is down")
		} else {
			log.Print("Verify Core is life")
		}
		healthAuth = server.authClient.HealthAuthSvc()
		if healthAuth != true {
			log.Print("Verify service is down")
		} else {
			log.Print("Verify service is life")
		}
	}
}

func (server *ServerStruct) handlerHealthFile() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		healthFile := server.fileClient.HealthFile()
		if healthFile != true {
			log.Print("File service is down")
		} else {
			log.Print("File service is life")
		}
	}
}

func (server *ServerStruct) handlerHealthPost() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		healthPost := server.postClient.HealthPost()
		if healthPost != true {
			log.Print("Post service is down")
		} else {
			log.Print("Post service is life")
		}
	}
}

func (server *ServerStruct) handlerVerifyRedir() http.HandlerFunc {
	var (
		tplBad *template.Template
		err    error
	)

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badAuthSvc.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		healthAuth := server.authClient.HealthAuthCore()
		if healthAuth != true {
			err := tplBad.Execute(writer, struct{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		healthAuth = server.authClient.HealthAuthSvc()
		if healthAuth != true {
			err := tplBad.Execute(writer, struct{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
		http.Redirect(writer, request, "http://localhost:10000/login", http.StatusTemporaryRedirect)
	}
}

func (server *ServerStruct) handlerPostsMovies() http.HandlerFunc {
	var (
		tpl    *template.Template
		tplBad *template.Template
		err    error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "movie.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badPostSvc.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		post := server.postClient.HealthPost()
		if post != true {
			err := tplBad.Execute(writer, struct{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		movies, err := server.postClient.PostsMovies(ctx)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = tpl.Execute(writer, movies)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (server *ServerStruct) handlerPostsGames() http.HandlerFunc {
	var (
		tpl    *template.Template
		tplBad *template.Template
		err    error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "game.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badPostSvc.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		post := server.postClient.HealthPost()
		if post != true {
			err := tplBad.Execute(writer, struct{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		games, err := server.postClient.PostsGames(ctx)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = tpl.Execute(writer, games)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (server *ServerStruct) handlerPostsSofts() http.HandlerFunc {
	var (
		tpl    *template.Template
		tplBad *template.Template
		err    error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "soft.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badPostSvc.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		post := server.postClient.HealthPost()
		if post != true {
			err := tplBad.Execute(writer, struct{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		soft, err := server.postClient.PostsSofts(ctx)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = tpl.Execute(writer, soft)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (server *ServerStruct) handlerPostsMusics() http.HandlerFunc {
	var (
		tpl    *template.Template
		tplBad *template.Template
		err    error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "music.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badPostSvc.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		post := server.postClient.HealthPost()
		if post != true {
			err := tplBad.Execute(writer, struct{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		musics, err := server.postClient.PostsMusics(ctx)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = tpl.Execute(writer, musics)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
