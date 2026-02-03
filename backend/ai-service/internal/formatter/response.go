package formatter

import (
	"encoding/json"
	"errors"
	"strings"
)

func ExtractJSON(raw string) ([]byte, error) {
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")

	if start == -1 || end == -1 || end <= start {
		return nil, errors.New("no valid JSON found")
	}

	return []byte(raw[start : end+1]), nil
}

func ValidateJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}