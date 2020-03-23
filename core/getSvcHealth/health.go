package getSvcHealth

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Url string

func GetHealth(urlSvc Url) (status bool, err error) {
	reqHealth, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/health", urlSvc), bytes.NewBuffer([]byte("Get Health")))
	if err != nil {
		return false, fmt.Errorf("can't create request: %w", err)
	}
	responseHealth, err := http.DefaultClient.Do(reqHealth)
	if err != nil {
		return false, fmt.Errorf("can't send request: %w", err)
	}
	defer func() {
		err = responseHealth.Body.Close()
		if err != nil {
			return
		}
	}()

	responseBodyHealth, err := ioutil.ReadAll(responseHealth.Body)
	if err != nil {
		return false, fmt.Errorf("can't parse response: %w", err)
	}

	if string(responseBodyHealth) != "Health ok" {
		return false, fmt.Errorf("service is down: %w", err)
	}

	return true, nil
}