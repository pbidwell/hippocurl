package main

import (
	"fmt"
	"hippocurl/modules"
	"os"

	"github.com/spf13/cobra"
)

var registeredModules []modules.HippoModule

func main() {
	var rootCmd = &cobra.Command{
		Use:   "hc",
		Short: "HippoCurl - A modular HTTP utility",
		Long:  "HippoCurl (hc) is a command-line tool for exploring and interacting with HTTP services.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("===================================")
				fmt.Println("         HIPPOCURL (hc)           ")
				fmt.Println("===================================")
				fmt.Println("A modular command-line tool for exploring and interacting with HTTP services.")
				fmt.Println("\nAvailable Modules:")
				for _, module := range registeredModules {
					fmt.Printf("%s %s: %s\n", module.Logo(), module.Name(), module.Description())
				}
				fmt.Println("\nUsage: hc <module> [arguments]")
				return
			}

			moduleName := args[0]
			for _, module := range registeredModules {
				if module.Name() == moduleName {
					module.Execute(args[1:]) // Pass remaining arguments
					return
				}
			}

			fmt.Printf("Error: Module '%s' not found.\n", moduleName)
		},
	}

	registeredModules = append(registeredModules, modules.ExploreModule{})
	registeredModules = append(registeredModules, modules.APIModule{})

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
