package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	lightsCmd.AddCommand(lightsOnCmd)
	lightsCmd.AddCommand(lightsOffCmd)

	rootCmd.AddCommand(lightsCmd)
}

var lightsCmd = &cobra.Command{
	Use:   "lights",
	Short: "Control lights",
}

var lightsOnCmd = &cobra.Command{
	Use:   "on [lightID]",
	Short: "Turn on a light",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lightID := args[0]
		fmt.Printf("Turning on light ID: %s\n", lightID)
		setLightState(lightID, true)
	},
}

var lightsOffCmd = &cobra.Command{
	Use:   "off [lightID]",
	Short: "Turn off a light",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lightID := args[0]
		fmt.Printf("Turning off light ID: %s\n", lightID)
		setLightState(lightID, false)
	},
}

func setLightState(lightID string, on bool) {
	hueBridgeIP := viper.GetString("hue_bridge_ip")
	applicationKey := viper.GetString("hue_application_key")

	if hueBridgeIP == "" || applicationKey == "" {
		fmt.Println("Hue Bridge IP or username not found in the config. Please run the discover and register commands first.")
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	apiEndpoint := fmt.Sprintf("https://%s/clip/v2/resource/light/%s", hueBridgeIP, lightID)

	requestBody := map[string]interface{}{
		"on": map[string]bool{
			"on": on,
		},
	}

	marshalledBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("error marshalling json: %v\n", err)
		return
	}

	req, err := http.NewRequest("PUT", apiEndpoint, bytes.NewReader(marshalledBody))
	if err != nil {
		fmt.Printf("error creating HTTP request: %v\n", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("hue-application-key", applicationKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error setting light state: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("failed to set light state. Status code: %d\n", resp.StatusCode)
		return
	}

	fmt.Printf("Light ID %s is now %s\n", lightID, map[bool]string{true: "on", false: "off"}[on])
}
