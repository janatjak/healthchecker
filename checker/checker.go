package checker

import (
	"net/http"
	"time"
)

type Checker struct {
	client *http.Client
	url    string
}

func (c *Checker) Check() (bool, error) {
	response, err := c.client.Get(c.url)
	if err != nil {
		return false, err
	}

	return response.StatusCode == 200, nil
}

func (c *Checker) Url() string {
	return c.url
}

func New(url string, timeout time.Duration) *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: timeout,
		},
		url: url,
	}
}
