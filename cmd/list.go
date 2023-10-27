package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		applicationKey := viper.GetString("hue_application_key")

		if hueBridgeIP == "" || applicationKey == "" {
			// Hue Bridge IP or username doesn't exist in the configuration, discover and register first
			log.Println("Hue Bridge IP or username not found in the config." +
				" Please run the discover and register commands first.")
			return
		}

		// Perform the device listing using the discovered Hue Bridge IP and username
		listDevices(hueBridgeIP, applicationKey)
	},
}

// Device represents the device structure in the new JSON response.
type Device struct {
	ID    string `json:"id"`
	IDV1  string `json:"id_v1"`
	Owner struct {
		RID   string `json:"rid"`
		RType string `json:"rtype"`
	} `json:"owner"`
	Metadata struct {
		Name      string `json:"name"`
		Archetype string `json:"archetype"`
	} `json:"metadata"`
	On struct {
		On bool `json:"on"`
	} `json:"on"`
	Dimming struct {
		Brightness  float64 `json:"brightness"`
		MinDimLevel float64 `json:"min_dim_level"`
	} `json:"dimming"`
	ColorTemperature struct {
		Mirek int  `json:"mirek"`
		Valid bool `json:"mirek_valid"`
	} `json:"color_temperature"`
	Color struct {
		XY struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"xy"`
		Gamut struct {
			Red struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"red"`
			Green struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"green"`
			Blue struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"blue"`
		} `json:"gamut"`
		GamutType string `json:"gamut_type"`
	} `json:"color"`
}

// Response represents the structure of the new JSON response.
type Response struct {
	Errors []interface{} `json:"errors"`
	Data   []Device      `json:"data"`
}

// listDevices lists all devices connected to the Philips Hue Bridge.
func listDevices(bridgeIP, applicationKey string) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	apiEndpoint := fmt.Sprintf("https://%s/clip/v2/resource/light", bridgeIP)

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		log.Printf("Error creating HTTP request: %v\n", err)
		return
	}
	req.Header.Add("hue-application-key", applicationKey)

	resp, err := client.Do(req)
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

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Error parsing JSON response: %v\n", err)
		return
	}

	fmt.Println("Devices connected to the Philips Hue Bridge:")
	for _, device := range response.Data {
		fmt.Printf("Device ID: %s\n", device.ID)
		fmt.Printf("Name: %s\n", device.Metadata.Name) // Extract and print the name field
		fmt.Printf("State: %s\n", getLightStateString(device))
		fmt.Println("--------------")
	}
}

// getLightStateString formats the light state as a string.
func getLightStateString(device Device) string {
	if device.On.On {
		return fmt.Sprintf("On (Brightness: %.2f, XY: (%.4f, %.4f))",
			device.Dimming.Brightness, device.Color.XY.X, device.Color.XY.Y)
	}
	return "Off"
}
