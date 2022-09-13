package scraper

import (
	"net/http"
	"time"
)

func validateURL(url string) error {
	client := newHTTPClient()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// TODO: set one http client instance for each scraper
func newHTTPClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: t,
	}
	return client
}
