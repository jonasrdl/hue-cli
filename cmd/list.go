package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List devices connected to the Philips Hue Bridge",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if Hue Bridge IP is in the configuration
		hueBridgeIP := viper.GetString("hue_bridge_ip")
		hueUsername := viper.GetString("hue_username")

		if hueBridgeIP == "" || hueUsername == "" {
			// Hue Bridge IP or username doesn't exist in the configuration, discover and register first
			log.Println("Hue Bridge IP or username not found in the configuration. Please run the discover and register commands first.")
			return
		}

		// Perform the device listing using the discovered Hue Bridge IP and username
		listDevices(hueBridgeIP, hueUsername)
	},
}

// listDevices lists all devices connected to the Philips Hue Bridge.
func listDevices(bridgeIP, username string) {
	apiEndpoint := fmt.Sprintf("http://%s/api/%s/lights", bridgeIP, username)

	resp, err := http.Get(apiEndpoint)
	if err != nil {
		log.Printf("Error listing devices: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to list devices. Status code: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return
	}

	var lights map[string]struct {
		State struct {
			On         bool
			Brightness int
			XY         []float64
		} `json:"state"`
		Name string `json:"name"`
	}
	err = json.Unmarshal(body, &lights)
	if err != nil {
		log.Printf("Error parsing JSON response: %v\n", err)
		return
	}

	fmt.Println("Devices connected to the Philips Hue Bridge:")
	for id, light := range lights {
		fmt.Printf("Light ID: %s\n", id)
		fmt.Printf("Name: %s\n", light.Name)
		fmt.Printf("State: %s\n", getLightStateString(light.State))
		fmt.Println("--------------")
	}
}

// getLightStateString formats the light state as a string.
func getLightStateString(state struct {
	On         bool
	Brightness int
	XY         []float64
}) string {
	if state.On {
		return fmt.Sprintf("On (Brightness: %d, XY: %v)", state.Brightness, state.XY)
	}
	return "Off"
}
