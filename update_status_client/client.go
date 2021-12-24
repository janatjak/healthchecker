package update_status_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UpdateStatusClient struct {
	client     *http.Client
	url        string
	apiKey     string
	lastStatus Status
}

type Status string

const (
	Operational         Status = "operational"
	UnderMaintenance    Status = "under_maintenance"
	DegradedPerformance Status = "degraded_performance"
	PartialOutage       Status = "partial_outage"
	MajorOutage         Status = "major_outage"
)

func (c *UpdateStatusClient) UpdateStatus(status Status) {
	if c.lastStatus != status {
		println("NEW status: " + status)
		c.call(status)
	}

	c.lastStatus = status
}

func (c *UpdateStatusClient) call(status Status) {
	jsonValue, _ := json.Marshal(map[string]interface{}{
		"component": map[string]Status{
			"status": status,
		},
	})

	request, _ := http.NewRequest(
		"PATCH",
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

func New(apiKey string, pageId string, componentId string) *UpdateStatusClient {
	return &UpdateStatusClient{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		url:    "https://api.statuspage.io/v1/pages/" + pageId + "/components/" + componentId,
		apiKey: apiKey,
	}
}
