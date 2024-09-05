package services

import (
	"fmt"
	"net/http"
	"os"
)

func UpdateStatus(id string, operation string) error {
	token := os.Getenv("BEARER_TOKEN_IFOOD")
	apiURL := "https://merchant-api.ifood.com.br/order/v1.0"
	var endpoint string

	switch operation {
	case "ConfirmOrder":
		endpoint = fmt.Sprintf("%s/orders/%s/confirm", apiURL, id)
	case "StartPreparation":
		endpoint = fmt.Sprintf("%s/orders/%s/startPreparation", apiURL, id)
	case "Dispatch":
		endpoint = fmt.Sprintf("%s/orders/%s/dispatch", apiURL, id)
	default:
		return fmt.Errorf("invalid operation")
	}

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error on update status ifood: %s | endpoint: %s", resp.Status, endpoint)
	}

	return nil
}
