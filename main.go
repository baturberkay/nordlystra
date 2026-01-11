package main

import (
	"flag"
	"fmt"
	"os"

	"nordlystra.app/config"
	"nordlystra.app/light"
)

func main() {
	cfg := config.LoadOrCreate()

	lightID := flag.Int("id", -1, "Device ID to change the state")
	on := flag.Bool("on", false, "Turn the device on")
	off := flag.Bool("off", false, "Turn the device off")
	reset := flag.Bool("reset", false, "Reset to brightest white")
	hue := flag.Int("hue", -1, "Hue value for the light color (0-65535)")
	sat := flag.Int("sat", -1, "Saturation value (0-255)")
	bri := flag.Int("bri", -1, "Brightness value (0-255)")
	list := flag.Bool("list", false, "List all lights and their current states")

	flag.Parse()

	baseURL := fmt.Sprintf("http://%s/api/%s", cfg.BridgeIP, cfg.Username)

	if *list {
		light.ListAll(baseURL)
		return
	}

	// If no flags provided, start interactive mode
	if *lightID == -1 && !*on && !*off && !*reset && *hue == -1 && *sat == -1 && *bri == -1 {
		light.RunInteractive(baseURL)
		return
	}

	if *lightID == -1 {
		fmt.Println("Error: Device ID is required")
		flag.Usage()
		os.Exit(1)
	}

	// Convert flags to the format UpdateState expects
	onValue := ""
	if *on {
		onValue = "true"
	} else if *off {
		onValue = "false"
	}

	if *reset {
		light.ResetState(baseURL, *lightID)
		return
	}

	light.UpdateState(baseURL, *lightID, onValue, *hue, *sat, *bri)
}
