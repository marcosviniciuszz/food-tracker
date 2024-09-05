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

func (f *FetcherService) fetchAndStore() error {
	// Create request
	req, err := http.NewRequest("GET", f.apiURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+f.token)

	// Send the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response status is 'No Content'
	if resp.StatusCode == http.StatusNoContent {
		log.Println("No content in response")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, body)
	}

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Check if the response body is empty
	if len(body) == 0 {
		log.Println("No data received in response")
		return nil
	}

	var data []bson.M
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if len(data) == 0 {
		log.Println("No valid data found in response")
		return nil
	}

	// Insert data into the MongoDB
	var ids []map[string]string
	for _, item := range data {
		if err := f.repo.Insert(context.TODO(), item); err != nil {
			return err
		}

		// Collect IDs for the POST acknowledge
		if id, ok := item["id"].(string); ok {
			ids = append(ids, map[string]string{"id": id})
		}
	}

	// Send IDs via POST request if any IDs are present
	if len(ids) > 0 {
		if err := f.sendPostRequest(ids); err != nil {
			return err
		}
	}

	return nil
}

func (f *FetcherService) sendPostRequest(ids []map[string]string) error {
	// Create request
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

	// Send request
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
