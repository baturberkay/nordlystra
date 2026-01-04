package bridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Authenticate(bridgeIP string) (string, error) {
	url := fmt.Sprintf("http://%s/api", bridgeIP)
	payload := []byte(`{"devicetype":"nordlystra#golang"}`)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("failed to authenticate: %w", err)
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result) == 0 || result[0]["error"] != nil {
		if errData, ok := result[0]["error"].(map[string]interface{}); ok {
			if errData["description"] == "link button not pressed" {
				fmt.Println("Link button not pressed. Please press the button on the Hue Bridge.")
				fmt.Println("Press Enter after pressing the button to continue...")
				fmt.Scanln()
				return Authenticate(bridgeIP)
			}
		}
		return "", fmt.Errorf("authentication failed: %v", result)
	}

	username := result[0]["success"].(map[string]interface{})["username"].(string)
	return username, nil
}
