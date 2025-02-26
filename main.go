package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func main() {
	// Define command-line flags.
	var org, filepath string
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
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	// Set up OAuth2 authentication.
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	// Create a GitHub client.
	client := github.NewClient(tc)
	// Read the file line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		username := strings.TrimSpace(scanner.Text())
		if username == "" {
			continue
		} // Prepare invitation options for each email.
		user, _, err := client.Users.Get(ctx, username)
		if err != nil {
			log.Printf("Error fetching user %s: %v", username, err)
			continue
		}

		// Step 2: Use the user's ID to invite
		inviteOptions := &github.CreateOrgInvitationOptions{
			InviteeID: user.ID, // Correct field for GitHub user ID
			Role:      github.String("direct_member"),
		}

		// Send the invitation.
		invitation, resp, err := client.Organizations.CreateOrgInvitation(ctx, org, inviteOptions)
		if err != nil {
			log.Printf("error sending invitation to %s: %v (HTTP status: %d)", username, err, resp.StatusCode)
			continue
		}

		fmt.Printf("Invitation sent to %s (Invitation ID: %d)\n", username, *invitation.ID)

		// Check for any scanning error.
		if err := scanner.Err(); err != nil {
			log.Fatalf("error reading file: %v", err)
		}

	}
}
