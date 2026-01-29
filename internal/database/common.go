package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func (db *Database) RechargeStatusCheck(
	partnerRequestID string,
) (string, error) {
	apiUrl := `https://v2a.rechargkit.biz/recharge/statusCheck`
	reqBody, err := json.Marshal(map[string]any{
		"partner_request_id": partnerRequestID,
	})
	if err != nil {
		return "", err
	}
	apiRequest, err := http.NewRequest(
		http.MethodPost,
		apiUrl,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return "", err
	}

	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set("Authorization", "Bearer "+os.Getenv("RKIT_API_TOKEN"))

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(apiRequest)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var apiResponse struct {
		Status  int    `json:"status"`
		Message string `json:"msg"`
	}

	if err := json.Unmarshal(respBytes, &apiResponse); err != nil {
		return "", err
	}

	if apiResponse.Status == 1 {
		return "SUCCESS", nil
	}

	if apiResponse.Status == 2 {
		return "PENDING", nil
	}

	if apiResponse.Status == 3 {
		return "FAILURE", nil
	}

	return "", fmt.Errorf("invalid status code in response: %s", apiResponse.Message)
}
