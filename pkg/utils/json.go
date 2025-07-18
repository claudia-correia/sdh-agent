package utils

import (
	"encoding/json"
	"fmt"
)

// MarshalJSON converts a request body to JSON
func MarshalJSON(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	return jsonData, nil
}

// UnmarshalJSON parses JSON data into a struct
func UnmarshalJSON(data []byte, jsonData interface{}) error {
	if err := json.Unmarshal(data, jsonData); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %w. Data: %s", err, string(data))
	}
	return nil
}
