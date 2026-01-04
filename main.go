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

	lightID := flag.Int("light-id", -1, "Light ID to change the state")
	on := flag.String("o", "", "Turn the light on or off (true/false)")
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

	if *lightID == -1 {
		fmt.Println("Error: Light ID is required")
		flag.Usage()
		os.Exit(1)
	}

	light.UpdateState(baseURL, *lightID, *on, *hue, *sat, *bri)
}
