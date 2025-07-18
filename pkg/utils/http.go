package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CreateDefaultHTTPClient creates an HTTP client with a default timeout
func CreateDefaultHTTPClient() *http.Client {
	return CreateHTTPClient(
		time.Second * 120, // Generous timeout for AI generation
	)
}

// CreateHTTPClient creates an HTTP client with a custom timeout
func CreateHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

// CreateRequest creates an HTTP request with JSON body and headers
func CreateRequest(method, url string, jsonData []byte, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default content type header
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

// ReadResponse reads and returns the response body
func ReadResponse(resp *http.Response) ([]byte, error) {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return respBody, nil
}

// SendJSONRequest sends a JSON request to an API and unmarshals the response
func SendJSONRequest(client *http.Client, method, url string, requestBody interface{}, responseBody interface{}, headers map[string]string) error {
	// Marshal request body to JSON
	jsonData, err := MarshalJSON(requestBody)
	if err != nil {
		return err
	}

	// Create request with headers
	req, err := CreateRequest(method, url, jsonData, headers)
	if err != nil {
		return err
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ReadResponse(resp)
	if err != nil {
		return err
	}

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API returned non-success status code %d: %s", resp.StatusCode, string(respBody))
	}

	// Unmarshal response if a response struct is provided
	if responseBody != nil {
		if err := UnmarshalJSON(respBody, responseBody); err != nil {
			return err
		}
	}

	return nil
}
