package log

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pbidwell/hippocurl/internal/config"
	"github.com/pbidwell/hippocurl/utils"
)

// LogModule implements the HippoModule interface
type LogModule struct{}

func (l LogModule) Name() string {
	return "log"
}

func (l LogModule) Description() string {
	return "Displays the location of the log file and the last 20 lines."
}

func (l LogModule) Use() string {
	return l.Name()
}

func (l LogModule) Execute(app *config.App, args []string) {
	logFilePath := app.LogFilePath

	utils.Print("Log File Location", utils.Header2)
	utils.Print(logFilePath, utils.NormalText)

	file, err := os.Open(logFilePath)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer file.Close()

	utils.Print("Log File Contents", utils.Header2)
	lines := readLastLines(file, 20)
	for _, line := range lines {
		utils.Print(line, utils.NormalText)
	}
}

func (l LogModule) Logo() string {
	return "📝"
}

// readLastLines reads the last n lines from a file
func readLastLines(file *os.File, n int) []string {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > n {
			lines = lines[1:]
		}
	}
	return lines
}
