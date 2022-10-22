package checker

import (
	"net/http"
	"time"
)

type Checker struct {
	client *http.Client
	url    string
}

func (c *Checker) Check() (bool, float64, error) {
	start := time.Now()
	response, err := c.client.Get(c.url)
	if err != nil {
		return false, 0, err
	}

	return response.StatusCode == 200, time.Since(start).Seconds(), nil
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
