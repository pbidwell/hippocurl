package modules

import (
	"context"
	"fmt"
	"hippocurl/utils"
	"os"
)

// ConfigModule implements the HippoModule interface
type ConfigModule struct{}

func (l ConfigModule) Name() string {
	return "config"
}

func (l ConfigModule) Description() string {
	return "Displays the location of the config file and the contents (up to 100 lines)."
}

func (l ConfigModule) Use() string {
	return l.Name()
}

func (l ConfigModule) Execute(ctx context.Context, args []string) {
	configFilePath, ok := ctx.Value(utils.ConfigFilePathKey).(string)
	if !ok {
		fmt.Println("Config file path not found in context.")
		return
	}

	utils.Print("Config File Location", utils.Header2)
	utils.Print(configFilePath, utils.NormalText)

	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Printf("Error opening config file: %v\n", err)
		return
	}
	defer file.Close()

	utils.Print("Config File Contents", utils.Header2)
	lines := readLastLines(file, 100)
	for _, line := range lines {
		utils.Print(line, utils.NormalText)
	}
}

func (l ConfigModule) Logo() string {
	return "🛠️"
}
