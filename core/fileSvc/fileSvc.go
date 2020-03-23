package fileSvc

import (
	"context"
	"github.com/burhon94/frontend/core/getSvcHealth"
	"log"
)

type Url string

type FileClient struct {
	url Url
}

func NewFileClient(url Url) *FileClient {
	return &FileClient{url: url}
}

func (f *FileClient) HealthFile() bool {
	status, err := getSvcHealth.GetHealth(getSvcHealth.Url(f.url))
	if err != nil {
		log.Print(err)
		return false
	}
	return status
}

func (f *FileClient) UploadFile(ctx context.Context, posterUrl, FileUrl string) (postPosterUrl, postFileUrl string, err error) {


	return "", "", nil
}