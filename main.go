package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func main() {
	// Define command-line flags.
	var org, filePath string
	args := os.Args[1:]
	if len(args) < 1 || len(args) > 2 {
		log.Fatalf("Usage: %s <github_organization_name> <file_path>", os.Args[0])
	}
	org = args[0]
	filepath = args[2]
	//takesgithub organization name
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	// Read the file line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		username := strings.TrimSpace(scanner.Text())
		if username == "" {
			continue // Skip empty lines.
		}

	}
}
