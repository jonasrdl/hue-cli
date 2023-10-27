package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register with the Philips Hue Bridge",
	Run: func(cmd *cobra.Command, args []string) {
		payload := `{"devicetype":"hue-cli", "generateclientkey": true}`

		bridgeIP := viper.GetString("hue_bridge_ip")
		if bridgeIP == "" {
			log.Println("No Philips Hue Bridge IP found in the configuration. Please run the discover command first.")
			return
		}

		apiEndpoint := fmt.Sprintf("http://%s/api", bridgeIP)
		resp, err := http.Post(apiEndpoint, "application/json", strings.NewReader(payload))
		if err != nil {
			log.Printf("Error registering with Philips Hue Bridge: %v\n", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v\n", err)
			return
		}

		fmt.Println("Registration response:")
		fmt.Println(string(body))

		if strings.Contains(string(body), `"success":{"username":`) {
			username, clientkey := extractUsernameAndClientkeyFromResponse(string(body))

			viper.Set("hue_application_key", username)
			viper.Set("hue_client_key", clientkey)

			// Username is used as hue_application_key, to set it as header for auth

			if err := viper.WriteConfig(); err != nil {
				log.Printf("Error writing the configuration: %v\n", err)
			}
		} else if strings.Contains(string(body), `"type":101`) {
			log.Println("Please press the link button on the Philips Hue Bridge, and then run the registration command again.")
		}
	},
}

// extractUsernameAndClientkeyFromResponse extracts the username and clientkey from the response body.
func extractUsernameAndClientkeyFromResponse(response string) (string, string) {
	var result []struct {
		Success struct {
			Username  string `json:"username"`
			ClientKey string `json:"clientkey"`
		} `json:"success"`
	}
	if err := json.Unmarshal([]byte(response), &result); err == nil {
		return result[0].Success.Username, result[0].Success.ClientKey
	}
	return "", ""
}
