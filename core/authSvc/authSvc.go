package authSvc

import (
	"github.com/burhon94/frontend/core/getSvcHealth"
	"log"
)

type Url string

type AuthClient struct {
	url Url
}

func NewAuthClient(url Url) *AuthClient {
	return &AuthClient{url: url}
}

func (a *AuthClient) HealthAuthSvc() bool {
	status, err := getSvcHealth.GetHealth("http://localhost:10000")
	if err != nil {
		log.Print(err)
		return false
	}
	return status
}

func (a *AuthClient) HealthAuthCore() bool {
	status, err := getSvcHealth.GetHealth("http://localhost:9999")
	if err != nil {
		log.Print(err)
		return false
	}
	return status
}
/*
func (a *AuthClient) VerifyReDirect(ctx context.Context) (fileWithByte[]byte, err error) {

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/verify", a.url),
		bytes.NewBuffer([]byte("VerifyMe")),
	)
	if err != nil {
		err = errors.New(fmt.Sprintf("can't create request: %v", err))
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		err = errors.New(fmt.Sprintf("can't sent request: %v", err))
		return nil, err
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			err = errors.New(fmt.Sprintf("can't success request: %v", err))
			return
		}
	}()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New("some Error while readBuf")
		return nil, err
	}

	switch response.StatusCode {
	case 200:
		return responseBody, nil
	case 400:
		return nil, err
	default:
		return nil, err
	}

}
*/