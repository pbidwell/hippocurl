package cmd

import (
	"fmt"
	"hippocurl/internal/config"
	"hippocurl/modules"
	"hippocurl/utils"
	"strings"

	// "log"
	"os"

	"github.com/spf13/cobra"
)

var registeredModules []modules.HippoModule

// logger            *log.Logger

var version = "dev" // Default to "dev" if not set at build time

var rootCmd = &cobra.Command{
	Use:   "hc",
	Short: "HippoCurl - A modular HTTP utility",
	Long:  fmt.Sprintf("%vHippoCurl (hc) is a command-line tool for exploring and interacting with HTTP and other web services.", utils.GetTitle()),
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func ExecuteModule(mod modules.HippoModule, args []string) {
	cfg := config.Load()
	logger := cfg.Logger

	logger.Printf("Executing module: [%s] with arguments [%s]", mod.Name(), strings.Join(args, ", "))
	mod.Execute(cfg, args)
	logger.Printf("Module [%s] execution complete", mod.Name())
}

func init() {
	rootCmd.Flags().String("configFilePath", "", "Global config file location. Defaults to $HOME/.hc/config.yml")
}
