/*
Copyright © 2025 Pablo Bidwell <bidwell.pablo@gmail.com>
*/
package cmd

import (
	"github.com/pbidwell/hippocurl/modules/api"

	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api [service_name] [route_name] [env_name]",
	Short: "Send requests to a configured API endpoint and view responses",
	Long: `The 'api' command sends HTTP requests to an API defined in your configuration file.
You can specify the service, route, and environment to automatically construct and execute
the request with any headers or authentication you’ve defined.

Example:
  hc api ServiceOne GetUser staging

If run without any arguments, the command enters an interactive mode, allowing you to 
select a service, route, and environment through a guided prompt.

This command is ideal for quickly testing or exploring API routes during development.`,
	Run: func(cmd *cobra.Command, args []string) {
		ExecuteModule(api.APIModule{}, args)
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
