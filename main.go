package main

import (
	"context"
	"fmt"
	"hippocurl/modules"
	"hippocurl/utils"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var registeredModules []modules.HippoModule

const (
	HC_CONFIG_DIR = ".hcconfig"
	HC_LOG_FILE   = "hc.log"
)

var logger *log.Logger

func main() {
	ctx := loadLoggingIntoContext(context.Background())
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

func loadLoggingIntoContext(ctx context.Context) context.Context {
	// Expand home directory if necessary
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configDir := filepath.Join(homeDir, HC_CONFIG_DIR)
		logFilePath := filepath.Join(configDir, HC_LOG_FILE)

		// Ensure config directory exists
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf("Failed to create config directory: %v\n", err)
		}

		// Setup logging
		logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Failed to open log file: %v\n", err)
			os.Exit(1)
		}
		logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)

		// Store logger and config in context
		ctx = context.WithValue(ctx, utils.LoggerKey, logger)
		ctx = context.WithValue(ctx, utils.LogFilePath, logFilePath)
	}
	return ctx
}
