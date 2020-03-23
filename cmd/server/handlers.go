package server

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

const (
	templatesPath = "web/templates"
)

func (server *ServerStruct) handlerIndexPage() http.HandlerFunc {
	var (
		tpl *template.Template
		tplBad *template.Template
		err error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "index.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badSvc.gohtml"),
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
		tpl *template.Template
		tplBad *template.Template
		err error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "admin", "index.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badSvc.gohtml"),
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
		tpl *template.Template
		err error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "admin", "postPage.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		err = tpl.Execute(writer, struct {}{})
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (server *ServerStruct) handlerPostNew() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		postTitle := request.PostFormValue("postTitle")
		if postTitle == "" {
			log.Print("postTitle can't be empty")
			return
		}

		postPoster := request.PostFormValue("postPoster")
		if postPoster == "" {
			log.Print("postPoster can't be empty")
			return
		}

		postCategory := request.PostFormValue("postCategory")
		if postCategory == "" {
			log.Print("postCategory can't be empty")
			return
		}

		postFile := request.PostFormValue("postFile")
		if postFile == "" {
			log.Print("postFile can't be empty")
			return
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		err := server.postClient.NewPost(ctx, postTitle, postPoster, postCategory, postFile)
		if err != nil {
			_, err := fmt.Fprintf(writer, "error while create new post: %v", err)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
	}
}

func (server *ServerStruct) handlerPostFiles() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		postPosterUrl := request.PostFormValue("postPosterUrl")
		if postPosterUrl == "" {
			log.Print("postPosterUrl can't be empty")
			return
		}

		postFilerUrl := request.PostFormValue("postFilerUrl")
		if postFilerUrl == "" {
			log.Print("postFilerUrl can't be empty")
			return
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		posterUrl, fileUrl, err := server.fileClient.UploadFile(ctx, postPosterUrl, postFilerUrl)
		if err != nil {
			_, err := fmt.Fprintf(writer, "error while create new post: %v", err)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
		log.Print(posterUrl, fileUrl)

		http.Redirect(writer, request, Root, http.StatusTemporaryRedirect)
	}
}

func (server *ServerStruct) handlerHealthAuth() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		healthAuth := server.authClient.HealthAuth()
		if healthAuth != true {
			log.Print("Verify service is down")
		}else {
			log.Print("Verify service is life")
		}
	}
}

func (server *ServerStruct) handlerHealthFile() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		healthFile := server.fileClient.HealthFile()
		if healthFile != true {
			log.Print("File service is down")
		}else {
			log.Print("File service is life")
		}
	}
}

func (server *ServerStruct) handlerHealthPost() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		healthPost := server.postClient.HealthPost()
		if healthPost != true {
			log.Print("Post service is down")
		}else {
			log.Print("Post service is life")
		}
	}
}

func (server *ServerStruct) handlerVerifyRedir() http.HandlerFunc {
	var (
		tpl *template.Template
		tplBad *template.Template
		err error
	)
	_, err = template.ParseFiles(
		filepath.Join(templatesPath, "badSvc.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
		)
	if err != nil {
		log.Print(err, tpl)
	}

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badAuthSvc.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		healthAuth := server.authClient.HealthAuth()
		if healthAuth != true {
			err := tplBad.Execute(writer, struct{}{})
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		ctx, _ := context.WithTimeout(request.Context(), time.Second)
		bytes, err := server.authClient.VerifyReDirect(ctx)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, err = writer.Write(bytes)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (server *ServerStruct) handlerPostsMovies() http.HandlerFunc {
	var (
		tpl *template.Template
		tplBad *template.Template
		err error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "movie.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badSvc.gohtml"),
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
		tpl *template.Template
		tplBad *template.Template
		err error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "game.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badSvc.gohtml"),
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
		tpl *template.Template
		tplBad *template.Template
		err error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "soft.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badSvc.gohtml"),
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
		tpl *template.Template
		tplBad *template.Template
		err error
	)

	tpl, err = template.ParseFiles(
		filepath.Join(templatesPath, "music.gohtml"),
		filepath.Join(templatesPath, "base.gohtml"),
	)
	if err != nil {
		panic(err)
	}

	tplBad, err = template.ParseFiles(
		filepath.Join(templatesPath, "badSvc.gohtml"),
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
