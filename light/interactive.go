package light

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printTable(lights map[string]Light) {
	ids := make([]int, 0, len(lights))
	for id := range lights {
		intID, _ := strconv.Atoi(id)
		ids = append(ids, intID)
	}
	sort.Ints(ids)

	// Calculate column widths
	nameWidth := 4 // "Name"
	for _, id := range ids {
		l := lights[strconv.Itoa(id)]
		if len(l.Name) > nameWidth {
			nameWidth = len(l.Name)
		}
	}

	// Print header
	fmt.Println()
	fmt.Println("  Lights in Your Environment")
	fmt.Println("  " + strings.Repeat("-", 60))
	fmt.Printf("  %-4s | %-*s | %-6s | %-10s | %-5s | %-10s\n", "ID", nameWidth, "Name", "On/Off", "Brightness", "Hue", "Saturation")
	fmt.Println("  " + strings.Repeat("-", 60))

	// Print rows
	for _, id := range ids {
		l := lights[strconv.Itoa(id)]
		state := "Off"
		if l.State.On {
			state = "On"
		}
		fmt.Printf("  %-4d | %-*s | %-6s | %-10d | %-5d | %-10d\n",
			id, nameWidth, l.Name, state, l.State.Bri, l.State.Hue, l.State.Sat)
	}
	fmt.Println("  " + strings.Repeat("-", 60))
	fmt.Println()
}

func printHelp() {
	fmt.Println()
	fmt.Println("  Available Commands:")
	fmt.Println("  -------------------")
	fmt.Println("  id {id} on               - Turn device on")
	fmt.Println("  id {id} off              - Turn device off")
	fmt.Println("  id {id} reset            - Reset to brightest white")
	fmt.Println("  id {id} bri {0-255}      - Set brightness")
	fmt.Println("  id {id} hue {0-65535}    - Set hue (color)")
	fmt.Println("  id {id} sat {0-255}      - Set saturation")
	fmt.Println()
	fmt.Println("  help, man                - Show this help")
	fmt.Println("  quit, q                  - Exit the application")
	fmt.Println()
}

func parseCommand(input string) (deviceID int, property, value string, err error) {
	parts := strings.Fields(input)

	if len(parts) < 3 {
		return 0, "", "", fmt.Errorf("invalid command format\n\n  Usage:\n    id {id} on/off           - Turn device on or off\n    id {id} bri/hue/sat {val} - Set property value\n\n  Examples:\n    id 1 on\n    id 1 bri 200\n\n  Type 'man' or 'help' for more info")
	}

	if parts[0] != "id" {
		return 0, "", "", fmt.Errorf("command must start with 'id'\n\n  Example: id 1 on\n\n  Type 'man' or 'help' for more info")
	}

	deviceID, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, "", "", fmt.Errorf("invalid device ID: %s", parts[1])
	}

	property = parts[2]

	// Handle on/off/reset (no value needed)
	if property == "on" || property == "off" || property == "reset" {
		return deviceID, property, "", nil
	}

	// Validate property before checking for value
	validProps := map[string]bool{"bri": true, "hue": true, "sat": true}
	if !validProps[property] {
		return 0, "", "", fmt.Errorf("invalid property '%s'\n\n  Available properties:\n    on/off             - Turn device on or off\n    bri  {0-255}       - Set brightness\n    hue  {0-65535}     - Set hue (color)\n    sat  {0-255}       - Set saturation", property)
	}

	// Valid properties need a value
	if len(parts) != 4 {
		return 0, "", "", fmt.Errorf("property '%s' requires a value\n\n  Example: id %d %s 200", property, deviceID, property)
	}

	value = parts[3]

	return deviceID, property, value, nil
}

func executeCommand(baseURL string, lightID int, property, value string) error {
	state := make(map[string]interface{})

	switch property {
	case "on":
		state["on"] = true
	case "off":
		state["on"] = false
	case "reset":
		state["on"] = true
		state["bri"] = 254
		state["sat"] = 0
	case "bri":
		val, err := strconv.Atoi(value)
		if err != nil || val < 0 || val > 255 {
			return fmt.Errorf("brightness must be between 0 and 255")
		}
		state["on"] = true
		state["bri"] = val
	case "hue":
		val, err := strconv.Atoi(value)
		if err != nil || val < 0 || val > 65535 {
			return fmt.Errorf("hue must be between 0 and 65535")
		}
		state["on"] = true
		state["hue"] = val
		state["sat"] = 254
	case "sat":
		val, err := strconv.Atoi(value)
		if err != nil || val < 0 || val > 255 {
			return fmt.Errorf("saturation must be between 0 and 255")
		}
		state["on"] = true
		state["sat"] = val
	}

	return ChangeState(baseURL, strconv.Itoa(lightID), state)
}

func RunInteractive(baseURL string) {
	reader := bufio.NewReader(os.Stdin)
	needsRefresh := true

	for {
		if needsRefresh {
			clearScreen()

			// Fetch and display lights
			lights, err := GetAll(baseURL)
			if err != nil {
				fmt.Printf("Error fetching lights: %v\n", err)
				return
			}
			printTable(lights)
		}
		needsRefresh = true // Default to refresh for next iteration

		// Show prompt
		fmt.Print("nordlystra> ")

		// Read input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return
		}

		input = strings.TrimSpace(input)
		if input == "" {
			needsRefresh = false
			continue
		}

		// Handle special commands
		if input == "quit" || input == "q" {
			fmt.Println("Goodbye!")
			return
		}

		if input == "help" || input == "man" {
			printHelp()
			needsRefresh = false
			continue
		}

		// Parse and execute command
		deviceID, property, value, err := parseCommand(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			needsRefresh = false
			continue
		}

		if err := executeCommand(baseURL, deviceID, property, value); err != nil {
			fmt.Printf("Error: %v\n", err)
			needsRefresh = false
			continue
		}

		fmt.Println("Done!")
	}
}
