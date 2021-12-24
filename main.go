package main

import (
	"encoding/json"
	"fmt"
	"github.com/janatjak/healthchecker/checker"
	"github.com/janatjak/healthchecker/update_status_client"
	"os"
	"sync"
	"time"
)

type Component struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type Config struct {
	APIKey          string      `json:"apiKey"`
	PageID          string      `json:"pageId"`
	MainComponentID string      `json:"mainComponentId"`
	Components      []Component `json:"components"`
}

const INTERVAL = 60

var componentStatuses []bool

func loop(index int, config *Config) {
	component := config.Components[index]

	updateStatusClient := update_status_client.New(config.APIKey, config.PageID, component.Id)
	checkerClient := checker.New(component.Url, time.Second*10)

	for {
		println("CHECK " + checkerClient.Url())
		status := update_status_client.Operational
		ok, _ := checkerClient.Check()
		if !ok {
			status = update_status_client.MajorOutage
			componentStatuses[index] = false
		} else {
			componentStatuses[index] = true
		}

		updateStatusClient.UpdateStatus(status)

		time.Sleep(time.Second * INTERVAL)
	}
}

func main() {
	configString := os.Getenv("CONFIG")
	configBytes := []byte(configString)

	var config = &Config{}
	err := json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Errorf("%v", err)
		return
	}

	fmt.Printf("%+v\n", config)

	if config.PageID == "" || config.APIKey == "" || config.MainComponentID == "" {
		fmt.Println("Missing required parameters")
		os.Exit(1)
	}

	componentStatuses = make([]bool, len(config.Components))
	for i := range componentStatuses {
		componentStatuses[i] = true
	}

	var wg sync.WaitGroup
	wg.Add(len(config.Components) + 1)

	// components
	for index := range config.Components {
		go loop(index, config)
	}

	// main component
	go func() {
		updateStatusClient := update_status_client.New(config.APIKey, config.PageID, config.MainComponentID)

		for {
			globalStatus := update_status_client.Operational
			for _, ok := range componentStatuses {
				if !ok {
					globalStatus = update_status_client.PartialOutage
					break
				}
			}

			updateStatusClient.UpdateStatus(globalStatus)

			time.Sleep(time.Second * INTERVAL / 2)
		}
	}()

	wg.Wait()
}
