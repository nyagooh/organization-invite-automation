package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	m "org/utils"

	"github.com/google/go-github/v53/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it; proceeding with system environment variables")
	}
	var org, filepath string
	args := os.Args[1:]
	if len(args) < 1 || len(args) > 2 {
		log.Fatalf("Usage: %s <github_organization_name> <file_path>", os.Args[0])
	}
	org = args[0]
	filepath = args[1]


	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}


	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	
	client := github.NewClient(tc)
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		username := strings.TrimSpace(scanner.Text())
		if username == "" {
			continue
		} 
		shouldInvite, err := m.ShouldInviteUser(ctx, client, org, username)
		if err != nil {
			log.Printf("error checking invitation status for %s: %v", username, err)
			continue
		}
		if !shouldInvite {
			fmt.Printf("Skipping %s: already a member or invitation pending\n", username)
			continue
		}

		if err := m.InviteUser(ctx, client, org, username); err != nil {
			log.Printf("error inviting %s: %v", username, err)
			continue
		}
		
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
}
}
