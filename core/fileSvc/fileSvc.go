package fileSvc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/burhon94/frontend/core/getSvcHealth"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type Url string

type FileClient struct {
	url Url
}

type fileStruct struct {
	FileName string `json:"fileName"`
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

func (f *FileClient) UploadFiles(request *http.Request) (filesUrls []string, err error) {
	status := f.HealthFile()
	if status != true {
		err = errors.New("fileSvc is down")
		return nil, err
	}

	const multipartMaxBytes = 10 * 1024 * 1024 * 1024 * 1024
	err = request.ParseMultipartForm(multipartMaxBytes)
	if err != nil {
		return nil, err
	}

	multiPartForm := request.MultipartForm
	files := multiPartForm.File

	for _, file := range files["files"] {
		openFile, err := file.Open()
		if err != nil {
			return nil, err
		}

		bodyReq := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyReq)
		part, err := writer.CreateFormFile("files", filepath.Base(file.Filename))
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(part, openFile)
		err = writer.Close()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/api/files/", f.url),
			bodyReq,
		)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		fileUrl, err := readResponse(response)
		if err != nil {
			return nil, err
		}

		var fileResp *fileStruct
		err = json.Unmarshal(fileUrl, &fileResp)
		if err != nil {
			return nil, err
		}

		filesUrls = append(filesUrls, fileResp.FileName)

	}

	return filesUrls, nil
}

func readResponse(response *http.Response) (responseData []byte, err error) {
	defer func() {
		err := response.Body.Close()
		if err != nil {
			return
		}
	}()

	responseData, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case 400:
		err = errors.New("bad request")
		return nil, err
	case 200:
		return responseData, nil
	}

	err = errors.New("unknown error")
	return nil, err
}