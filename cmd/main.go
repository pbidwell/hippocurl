package main

import (
	"context"
	"fmt"
	"hippocurl/modules"
	"hippocurl/utils"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	registeredModules []modules.HippoModule
	logger            *log.Logger
)

var version = "dev" // Default to "dev" if not set at build time

func main() {
	var ctx context.Context

	ctx = utils.LoadLoggerIntoContext(context.Background())
	logger = ctx.Value(utils.LoggerKey).(*log.Logger)
	logger.Println("==== HippoCurl Started ====")
	defer logger.Println("==== HippoCurl Terminated ====")

	ctx, err := utils.LoadConfigIntoContext(ctx)
	if err != nil {
		log.Fatalf("error initializing config: %v", err)
	}

	rootCmd := &cobra.Command{
		Use:   "hc",
		Short: "HippoCurl - A modular HTTP utility",
		Long:  "HippoCurl (hc) is a command-line tool for exploring and interacting with HTTP and other web services.",
	}

	// Version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintFieldValuePair("version", version)
		},
	}

	rootCmd.AddCommand(versionCmd)

	// Register modules as subcommands
	registeredModules = []modules.HippoModule{
		modules.ExploreModule{},
		modules.APIModule{},
		modules.LogModule{},
		modules.ConfigModule{},
	}

	for _, module := range registeredModules {
		mod := module // Capture variable to avoid loop variable issues
		cmd := &cobra.Command{
			Use:                   mod.Use(),
			Short:                 mod.Description(),
			DisableFlagsInUseLine: true,
			Run: func(cmd *cobra.Command, args []string) {
				logger.Printf("Executing module: [%s] with arguments [%s]", mod.Name(), strings.Join(args, ", "))
				mod.Execute(ctx, args)
				logger.Printf("Module [%s] execution complete", mod.Name())
			},
		}
		rootCmd.AddCommand(cmd)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		logger.Printf("Error: %v", err)
		os.Exit(1)
	}
}
