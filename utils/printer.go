package utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// HeadingLevel defines different levels of text emphasis
type HeadingLevel int

const (
	// Heading levels
	ModuleTitle HeadingLevel = iota
	Header1
	Header1_WithOpenDelimeter
	Header1_ClosedDelimeter
	Header1_Alternate
	Header2
	Hint
	NormalText
)

// Print prints a string with formatting based on the heading level
func Print(text string, level HeadingLevel) {
	switch level {
	case ModuleTitle:
		color.New(color.FgHiBlue, color.Bold).Printf("\n===== %s =====\n\n", strings.ToUpper(text))
	case Header1:
		color.New(color.FgMagenta, color.Bold).Printf("\n%s\n", text)
	case Header1_Alternate:
		color.New(color.FgHiGreen, color.Bold).Printf("\n%s\n", text)
	case Header2:
		color.New(color.FgCyan).Printf("\n%s\n", text)
	case NormalText:
		color.New(color.FgWhite).Println(text)
	case Hint:
		color.New(color.FgYellow).Printf("\n* Hint: %s *\n", text)
	}
}

func PrintTitle() {
	color.New(color.FgHiBlue, color.Bold).Print(`
===================================
	HIPPOCURL (hc)           
===================================
`)
}

// PrintHeaders prints key-value pairs from a map where values are arrays, consolidating them into a space-delimited string
// and ensuring the keys are printed in alphabetical order
func PrintHeaders(data map[string][]string) {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		values := data[key]
		if len(values) > 0 {
			fmt.Printf("%s: %s\n", key, strings.Join(values, " "))
		}
	}
}
