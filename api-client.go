package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type APIClient struct {
	Authorization string
	UserAgent     string

	httpClient *http.Client
}

func (x *APIClient) Init() {
	x.httpClient = &http.Client{
		Timeout: 40 * time.Second,
	}
}

func (x *APIClient) GetAPI(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("User-Agent", x.UserAgent)
	req.Header.Add("Referer", "https://fansly.com/")
	req.Header.Add("Authorization", x.Authorization)

	res, err := x.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("invalid response status code: %d", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(target)
	if err != nil {
		return err
	}

	return nil
}
