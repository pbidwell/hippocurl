package main

import (
	"context"
	"fmt"
	"hippocurl/modules"
	"hippocurl/utils"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var registeredModules []modules.HippoModule

var logger *log.Logger

func main() {
	ctx := utils.LoadLoggerIntoContext(context.Background())
	ctx = utils.LoadConfig(ctx)

	logger = ctx.Value(utils.LoggerKey).(*log.Logger)
	logger.Println("==== HippoCurl Started ====")
	defer logger.Println("==== HippoCurl Terminated ====")

	var rootCmd = &cobra.Command{
		Use:   "hc",
		Short: "HippoCurl - A modular HTTP utility",
		Long:  "HippoCurl (hc) is a command-line tool for exploring and interacting with HTTP services.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				utils.PrintTitle()
				utils.Print("A modular command-line tool for exploring and interacting with HTTP services.", utils.NormalText)
				utils.Print("Available Modules", utils.Header1)
				for _, module := range registeredModules {
					utils.Print(fmt.Sprintf("%s %s - %s", module.Logo(), module.Name(), module.Description()), utils.NormalText)
				}
				utils.Print("\nUsage: hc <module> [arguments]\n", utils.NormalText)
				return
			}

			moduleName := args[0]
			for _, module := range registeredModules {
				if module.Name() == moduleName {
					logger.Printf("Executing module: [%s] with arguments [%s]", moduleName, strings.Join(args[1:], ", "))
					module.Execute(ctx, args[1:]) // Pass remaining arguments
					logger.Printf("Module [%s] execution complete", moduleName)
					return
				}
			}

			fmt.Printf("Error: Module '%s' not found.\n", moduleName)
		},
	}

	registeredModules = append(registeredModules, modules.ExploreModule{})
	registeredModules = append(registeredModules, modules.APIModule{})
	registeredModules = append(registeredModules, modules.LogModule{})
	registeredModules = append(registeredModules, modules.ConfigModule{})

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		logger.Printf("Error: %v\n", err)
		return
	}
}
