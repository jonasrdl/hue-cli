package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	lightsOnCmd.Flags().String("light-id", "", "ID of the light to control")
	lightsOffCmd.Flags().String("light-id", "", "ID of the light to control")

	lightsCmd.AddCommand(lightsOnCmd)
	lightsCmd.AddCommand(lightsOffCmd)

	rootCmd.AddCommand(lightsCmd)
}

var lightsCmd = &cobra.Command{
	Use:   "lights",
	Short: "Control lights",
}

var lightsOnCmd = &cobra.Command{
	Use:   "on [lightName]",
	Short: "Turn on a light",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lightName := args[0]
		lightID, _ := cmd.Flags().GetString("light-id")

		if lightID == "" {
			fmt.Printf("Turning on light by name: %s\n", lightName)
			id, err := getLightIDByName(lightName)
			if err != nil {
				fmt.Printf("Error finding light ID: %v\n", err)
				return
			}
			lightID = id
		} else {
			fmt.Printf("Turning on light by ID: %s\n", lightID)
		}

		setLightState(lightID, true)
	},
}

var lightsOffCmd = &cobra.Command{
	Use:   "off [lightName]",
	Short: "Turn off a light",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lightName := args[0]
		lightID, _ := cmd.Flags().GetString("light-id")

		if lightID == "" {
			fmt.Printf("Turning off light by name: %s\n", lightName)
			id, err := getLightIDByName(lightName)
			if err != nil {
				fmt.Printf("Error finding light ID: %v\n", err)
				return
			}
			lightID = id
		} else {
			fmt.Printf("Turning off light by ID: %s\n", lightID)
		}

		setLightState(lightID, false)
	},
}

func getLightIDByName(name string) (string, error) {
	hueBridgeIP := viper.GetString("hue_bridge_ip")
	applicationKey := viper.GetString("hue_application_key")

	if hueBridgeIP == "" || applicationKey == "" {
		return "", fmt.Errorf("hue bridge IP or username not found in the config")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	apiEndpoint := fmt.Sprintf("https://%s/clip/v2/resource/light", hueBridgeIP)

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Add("hue-application-key", applicationKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error listing devices: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to list devices. Status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON response: %v", err)
	}

	for _, device := range response.Data {
		if device.Metadata.Name == name {
			return device.ID, nil
		}
	}

	return "", fmt.Errorf("light with name %s not found", name)
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
