package cmd

import (
	"github.com/pbidwell/hippocurl/modules/explore"

	"github.com/spf13/cobra"
)

// exploreCmd represents the explore command
var exploreCmd = &cobra.Command{
	Use:   "explore [hostname|IP]",
	Short: "Profile a host or IP",
	Long: `The 'explore' command analyzes a given hostname or IP address by performing 
a series of passive and active checks. These include DNS resolution, IP geolocation, 
open port scanning, reverse DNS lookups, HTTP probing, and more.

Example:
  hc explore example.com

This command is designed to help users get a high-level understanding of 
network assets and their exposure.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		ExecuteModule(explore.ExploreModule{}, args)
	},
}

func init() {
	rootCmd.AddCommand(exploreCmd)
}
