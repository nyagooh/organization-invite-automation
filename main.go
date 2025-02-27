package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
		shouldInvite, err := shouldInviteUser(ctx, client, org, username)
		if err != nil {
			log.Printf("error checking invitation status for %s: %v", username, err)
			continue
		}
		if !shouldInvite {
			fmt.Printf("Skipping %s: already a member or invitation pending\n", username)
			continue
		}
		user, _, err := client.Users.Get(ctx, username)
		if err != nil {
			log.Printf("Error fetching user %s: %v", username, err)
			continue
		}

		inviteOptions := &github.CreateOrgInvitationOptions{
			InviteeID: user.ID,
			Role:      github.String("direct_member"),
		}

		invitation, resp, err := client.Organizations.CreateOrgInvitation(ctx, org, inviteOptions)
		if err != nil {
			log.Printf("error sending invitation to %s: %v (HTTP status: %d)", username, err, resp.StatusCode)
			continue
		}

		fmt.Printf("Invitation sent to %s (Invitation ID: %d)\n", username, *invitation.ID)

		if err := scanner.Err(); err != nil {
			log.Fatalf("error reading file: %v", err)
		}

	}
}

func shouldInviteUser(ctx context.Context, client *github.Client, org, username string) (bool, error) {
	isMember, _, err := client.Organizations.IsMember(ctx, org, username)
	if err != nil {
		return false, fmt.Errorf("error checking membership for %s: %v", username, err)
	}
	if isMember {
		return false, nil
	}

	invitations, _, err := client.Organizations.ListPendingOrgInvitations(ctx, org, nil)
	if err != nil {
		return false, fmt.Errorf("error listing pending invitations for organization %s: %v", org, err)
	}

	for _, invite := range invitations {
		if invite != nil && strings.EqualFold(invite.GetLogin(), username) {
			return false, nil
		}
	}
	return true, nil
}
