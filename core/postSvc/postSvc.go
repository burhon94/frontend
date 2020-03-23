package postSvc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/burhon94/frontend/core/getSvcHealth"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Url string

type PostClient struct {
	url Url
}

func NewPostClient(url Url) *PostClient {
	return &PostClient{url: url}
}

type postReq struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	DataAndTime time.Time `json:"dataandtime"`
	Poster      string    `json:"poster"`
	Category    string    `json:"category"`
	FileUrl     string    `json:"file_url"`
}

type PostNew struct {
	Title    string `json:"title"`
	Poster   string `json:"poster"`
	Category string `json:"category"`
	FileUrl  string `json:"file_url"`
}

var ErrUnknown = errors.New("unknown error")

func (p *PostClient) HealthPost() bool {
	status, err := getSvcHealth.GetHealth(getSvcHealth.Url(p.url))
	if err != nil {
		log.Print(err)
		return false
	}
	return status
}

func getPosts(ctx context.Context, svcUrl, pattern string) (post []postReq, err error) {
	post = make([]postReq, 0)
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s%s", svcUrl, pattern),
		bytes.NewBuffer([]byte("")),
	)
	if err != nil {
		return post, fmt.Errorf("can't create request: %w", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return post, fmt.Errorf("can't send request: %w", err)
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			return
		}
	}()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return post, fmt.Errorf("can't parse response: %w", err)
	}


	switch response.StatusCode {
	case 500:
		return post, errors.New("can't decode response")
	}
	err = json.Unmarshal(responseBody, &post)
	if err != nil {
		return post, fmt.Errorf("can't decode response: %w", err)
	}
	return post, nil
}

func (p *PostClient) GetAllPosts(ctx context.Context) (posts []postReq, err error) {
	posts, err = getPosts(ctx, string(p.url), "/api/posts")
	if err != nil {
		return posts, err
	}

	return posts, nil
}

func (p *PostClient) PostsMovies(ctx context.Context) (postMovies []postReq, err error) {
	postMovies, err = getPosts(ctx, string(p.url), "/api/posts/movies")
	if err != nil {
		return postMovies, err
	}

	return postMovies, nil
}

func (p *PostClient) PostsGames(ctx context.Context) (postGames []postReq, err error) {
	postGames, err = getPosts(ctx, string(p.url), "/api/posts/games")
	if err != nil {
		return postGames, err
	}

	return postGames, nil
}

func (p *PostClient) PostsSofts(ctx context.Context) (postSofts []postReq, err error) {
	postSofts, err = getPosts(ctx, string(p.url), "/api/posts/programmes")
	if err != nil {
		return postSofts, err
	}

	return postSofts, nil
}

func (p *PostClient) PostsMusics(ctx context.Context) (postMusics []postReq, err error) {
	postMusics, err = getPosts(ctx, string(p.url), "/api/posts/musics")
	if err != nil {
		return postMusics, err
	}

	return postMusics, nil
}

func (p *PostClient) NewPost(ctx context.Context,  postTitle, postPoster, postCategory, postFile string) (err error) {
	requestData := PostNew{
		Title:    postTitle,
		Poster:   postPoster,
		Category: postCategory,
		FileUrl:  postFile,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("can't encode requestBody %v: %w", requestData, err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/post/0", p.url),
		bytes.NewBuffer(requestBody),
	)

	if err != nil {
		return fmt.Errorf("can't create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("can't send request: %w", err)
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			return
		}
	}()

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("can't parse response: %w", err)
	}

	switch response.StatusCode {
	case 200:
		return nil
	case 400:
		return err
	default:
		return ErrUnknown
	}
}