package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"strings"
)

func init() {
	rootCmd.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register with the Philips Hue Bridge",
	Run: func(cmd *cobra.Command, args []string) {
		payload := `{"devicetype":"hue-cli"}`

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
			username := extractUsernameFromResponse(string(body))

			viper.Set("hue_username", username)
			if err := viper.WriteConfig(); err != nil {
				log.Printf("Error writing the configuration: %v\n", err)
			}
		} else if strings.Contains(string(body), `"type":101`) {
			log.Println("Please press the link button on the Philips Hue Bridge, and then run the registration command again.")
		}
	},
}

// extractUsernameFromResponse extracts the username from the response body.
func extractUsernameFromResponse(response string) string {
	var result []struct {
		Success struct {
			Username string `json:"username"`
		} `json:"success"`
	}
	if err := json.Unmarshal([]byte(response), &result); err == nil {
		return result[0].Success.Username
	}
	return ""
}
