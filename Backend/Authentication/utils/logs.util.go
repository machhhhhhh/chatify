package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HeadersToString(header http.Header) string {
	var result string

	if len(header) == 0 {
		return result
	}

	for key, value := range header {
		for i := range value {
			result += fmt.Sprintf("%s: %s\n", key, value[i])
		}
	}

	return result
}

func BodyToString(body map[string]map[string]any) (string, error) {
	if body == nil {
		return "", nil
	}

	json_byte, err := json.Marshal(body)

	return string(json_byte), err
}
