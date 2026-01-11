package light

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type State struct {
	On  bool `json:"on"`
	Hue int  `json:"hue"`
	Sat int  `json:"sat"`
	Bri int  `json:"bri"`
}

type Light struct {
	Name  string `json:"name"`
	State State  `json:"state"`
}

func GetState(baseURL, lightID string) (State, error) {
	url := fmt.Sprintf("%s/lights/%s", baseURL, lightID)

	resp, err := http.Get(url)
	if err != nil {
		return State{}, fmt.Errorf("failed to fetch light state: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return State{}, fmt.Errorf("failed to fetch light state: HTTP %d", resp.StatusCode)
	}

	var lightResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&lightResponse); err != nil {
		return State{}, fmt.Errorf("failed to parse response: %w", err)
	}

	stateData, ok := lightResponse["state"].(map[string]interface{})
	if !ok {
		return State{}, fmt.Errorf("invalid response format: missing 'state' field")
	}

	state := State{
		On:  stateData["on"].(bool),
		Hue: int(stateData["hue"].(float64)),
		Sat: int(stateData["sat"].(float64)),
		Bri: int(stateData["bri"].(float64)),
	}

	return state, nil
}

func GetAll(baseURL string) (map[string]Light, error) {
	url := fmt.Sprintf("%s/lights", baseURL)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch lights: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var lights map[string]Light
	if err := json.Unmarshal(body, &lights); err != nil {
		return nil, fmt.Errorf("failed to parse lights: %w", err)
	}

	return lights, nil
}

func ChangeState(baseURL, lightID string, state map[string]interface{}) error {
	url := fmt.Sprintf("%s/lights/%s/state", baseURL, lightID)

	jsonData, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var response interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if respArray, ok := response.([]interface{}); ok {
		for _, item := range respArray {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if errData, exists := itemMap["error"]; exists {
					return fmt.Errorf("error from API: %v", errData)
				}
			}
		}
	} else if respMap, ok := response.(map[string]interface{}); ok {
		if errData, exists := respMap["error"]; exists {
			return fmt.Errorf("error from API: %v", errData)
		}
	}

	return nil
}

func ResetState(baseURL string, lightID int) {
	state := map[string]interface{}{
		"on":  true,
		"bri": 254,
		"sat": 0,
	}

	if err := ChangeState(baseURL, fmt.Sprintf("%d", lightID), state); err != nil {
		log.Fatalf("Error resetting light state: %v", err)
	}
	fmt.Println("Device reset to brightest white.")
}

func UpdateState(baseURL string, lightID int, on string, hue, sat, bri int) {
	currentState, err := GetState(baseURL, fmt.Sprintf("%d", lightID))
	if err != nil {
		log.Fatalf("Failed to fetch light state: %v", err)
	}

	state := make(map[string]interface{})

	if on != "" {
		if on == "true" {
			if !currentState.On {
				state["on"] = true
			}
		} else if on == "false" {
			if currentState.On {
				state["on"] = false
			}
		} else {
			fmt.Println("Error: Invalid value for -o flag. Use 'true' or 'false'.")
			os.Exit(1)
		}
	} else if hue != -1 || sat != -1 || bri != -1 {
		if !currentState.On {
			state["on"] = true
		}
	}

	if hue != -1 {
		state["hue"] = hue
	}
	if sat != -1 {
		state["sat"] = sat
	}
	if bri != -1 {
		state["bri"] = bri
	}

	if len(state) == 0 {
		fmt.Println("Error: No state changes provided")
		os.Exit(1)
	}

	if err := ChangeState(baseURL, fmt.Sprintf("%d", lightID), state); err != nil {
		log.Fatalf("Error changing light state: %v", err)
	}
	fmt.Println("Light state updated successfully.")
}

func ListAll(baseURL string) {
	allLights, err := GetAll(baseURL)
	if err != nil {
		log.Fatalf("Failed to fetch lights: %v", err)
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	defer ui.Close()

	ids := make([]int, 0, len(allLights))
	for id := range allLights {
		intID, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalf("Failed to convert light ID to integer: %v", err)
		}
		ids = append(ids, intID)
	}
	sort.Ints(ids)

	table := widgets.NewTable()
	table.Title = "Lights in Your Environment"
	table.Rows = [][]string{
		{"ID", "Name", "On/Off", "Brightness", "Hue", "Saturation"},
	}

	for _, id := range ids {
		l := allLights[strconv.Itoa(id)]
		state := "Off"
		if l.State.On {
			state = "On"
		}
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(id),
			l.Name,
			state,
			fmt.Sprintf("%d", l.State.Bri),
			fmt.Sprintf("%d", l.State.Hue),
			fmt.Sprintf("%d", l.State.Sat),
		})
	}

	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.SetRect(0, 0, 100, len(table.Rows)+len(ids)+2)
	table.RowSeparator = true

	ui.Render(table)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent && e.ID == "q" {
			break
		}
	}
}
