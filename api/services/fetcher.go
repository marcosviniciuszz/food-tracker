package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"food-tracker/repositories"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type FetcherService struct {
	repo       repositories.OrderRepository
	apiURL     string
	postAPIURL string
	interval   time.Duration
	token      string
}

func NewFetcherService(repo repositories.OrderRepository) *FetcherService {
	token := os.Getenv("BEARER_TOKEN_IFOOD")

	return &FetcherService{
		repo:       repo,
		apiURL:     "https://merchant-api.ifood.com.br/events/v1.0/events:polling",
		postAPIURL: "https://merchant-api.ifood.com.br/events/v1.0/events/acknowledgment",
		interval:   30 * time.Second,
		token:      token,
	}
}

func (f *FetcherService) Start() {
	go func() {
		for {
			if err := f.fetchAndStore(); err != nil {
				log.Printf("Error fetching and storing data: %v", err)
			}
			time.Sleep(f.interval)
		}
	}()
}

// fetchAndStore handles the fetching of data and storing it in MongoDB.
func (f *FetcherService) fetchAndStore() error {
	data, err := f.fetchData()
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	ids, err := f.processData(data)
	if err != nil {
		return err
	}

	if len(ids) > 0 {
		return f.sendPostRequest(ids)
	}

	return nil
}

// fetchData performs the HTTP GET request and returns the response data.
func (f *FetcherService) fetchData() ([]bson.M, error) {
	req, err := http.NewRequest("GET", f.apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+f.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		log.Println("No content in response")
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		log.Println("No data received in response")
		return nil, nil
	}

	var data []bson.M
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// processData processes the fetched data and inserts valid entries into MongoDB.
func (f *FetcherService) processData(data []bson.M) ([]map[string]string, error) {
	var ids []map[string]string

	for _, item := range data {

		if fullCode, ok := item["fullCode"].(string); ok && fullCode == "PLACED" {
			if err := f.repo.Insert(context.TODO(), item); err != nil {
				return nil, err
			}
		}

		if id, ok := item["id"].(string); ok {
			ids = append(ids, map[string]string{"id": id})
		}
	}

	return ids, nil
}

// sendPostRequest sends a POST request with the collected IDs.
func (f *FetcherService) sendPostRequest(ids []map[string]string) error {
	payload, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", f.postAPIURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+f.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("POST request failed: %s", body)
	}

	return nil
}
