package main

import (
	"log"
	"os"

	"brewls/internal/brewls"
)

func main() {
	jsonOutput, err := brewls.ExecuteBrewInfoCommand()
	if err != nil {
		log.Fatalf("Failed to execute brew command: %v", err)
	}

	brewInfo, err := brewls.ParseBrewInfoJSON(jsonOutput)
	if err != nil {
		log.Fatalf("Failed to parse brew info: %v", err)
	}

	brewls.BuildReverseDependencyGraph(brewInfo) // Call the new function

	brewls.FormatBrewOutput(brewInfo, os.Stdout)
}
