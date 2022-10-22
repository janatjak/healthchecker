package metric_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MetricClient struct {
	client   *http.Client
	url      string
	apiKey   string
	metricId string
}

type metric struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

func (c *MetricClient) AddMetric(timestamp int64, value float64) {
	fmt.Printf("NEW metric: %f\n", value)

	jsonValue, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			c.metricId: []metric{
				{
					Timestamp: timestamp,
					Value:     value,
				},
			},
		},
	})

	request, _ := http.NewRequest(
		"POST",
		c.url,
		bytes.NewBuffer(jsonValue),
	)
	request.Header.Add("Authorization", c.apiKey)
	request.Header.Add("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	println(response.StatusCode)
}

func New(apiKey string, pageId string, metricId string) *MetricClient {
	return &MetricClient{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		url:      "https://api.statuspage.io/v1/pages/" + pageId + "/metrics/data",
		apiKey:   apiKey,
		metricId: metricId,
	}
}
