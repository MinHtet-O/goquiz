package scraper

import (
	"net/http"
	"time"
)

func validateImageURL(url string) error {
	client := clientHTTP()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

// TODO: set one http client instance for each scraper
func clientHTTP() *http.Client {
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
