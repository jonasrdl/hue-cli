package cmd

import (
	"context"
	"fmt"
	"github.com/grandcat/zeroconf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {
	rootCmd.AddCommand(discoverCmd)
}

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover Philips Hue Bridge on the local network",
	Run: func(cmd *cobra.Command, args []string) {
		bridge, err := discoverBridgeUsingMDNS()
		if err != nil {
			fmt.Printf("Error discovering Philips Hue Bridge: %v\n", err)
		} else if bridge != nil {
			fmt.Printf("Discovered Philips Hue Bridge:\nID: %s\nInternal IP: %s\n", bridge.ID, bridge.InternalIP)

			// Store the IP address in the config file
			viper.Set("hue_bridge_ip", bridge.InternalIP)
			if err := viper.WriteConfig(); err != nil {
				fmt.Printf("Error writing config file: %v\n", err)
			}

			// Exit the program after a bridge is found
			os.Exit(0)
		} else {
			fmt.Println("No Philips Hue Bridge found on the local network.")
		}
	},
}

// Bridge exposes a hardware bridge through a struct
type Bridge struct {
	ID         string `json:"id"`
	InternalIP string `json:"internalipaddress"`
}

// discoverBridgeUsingMDNS discovers a Philips Hue Bridge using mDNS.
func discoverBridgeUsingMDNS() (*Bridge, error) {
	var bridges []Bridge

	// Initialize a new mDNS resolver
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}

	// Define a channel to receive discovered services
	entries := make(chan *zeroconf.ServiceEntry)

	// Start resolving "_hue._tcp" services using mDNS
	err = resolver.Browse(context.Background(), "_hue._tcp", "local.", entries)
	if err != nil {
		return nil, err
	}

	// Listen for discovered services
	for entry := range entries {
		var addr string
		for _, ip := range entry.AddrIPv4 {
			addr = ip.String()
		}
		if addr == "" {
			for _, ip := range entry.AddrIPv6 {
				addr = ip.String()
			}
		}

		bridges = append(bridges, Bridge{
			ID:         entry.Instance,
			InternalIP: addr,
		})

		// Assuming the first bridge found is the correct one, update the configuration
		if len(bridges) > 0 {
			viper.Set("hue_bridge_id", bridges[0].ID)
			viper.Set("hue_bridge_ip", bridges[0].InternalIP)
			if err := viper.WriteConfig(); err != nil {
				fmt.Printf("Error writing config file: %v\n", err)
			}

			// Log the discovered Hue Bridge
			log.Printf("Discovered Philips Hue Bridge:\nID: %s\nInternal IP: %s", bridges[0].ID, bridges[0].InternalIP)

			// Exit the program after a bridge is found
			os.Exit(0)
		}
	}

	if len(bridges) == 0 {
		return nil, nil
	}

	return &bridges[0], nil
}
