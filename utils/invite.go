package org

import (
	"context"
	"fmt"
	"github.com/google/go-github/v53/github"
	"strings"
)
func ShouldInviteUser(ctx context.Context, client *github.Client, org, username string) (bool, error) {
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

func InviteUser(ctx context.Context, client *github.Client, org, username string) error {

	user, _, err := client.Users.Get(ctx, username)
	if err != nil {
		return fmt.Errorf("error fetching user %s: %v", username, err)
	}


	inviteOptions := &github.CreateOrgInvitationOptions{
		InviteeID: user.ID,
		Role:      github.String("direct_member"),
	}

	
	invitation, resp, err := client.Organizations.CreateOrgInvitation(ctx, org, inviteOptions)
	if err != nil {
		return fmt.Errorf("error sending invitation to %s: %v (HTTP status: %d)", username, err, resp.StatusCode)
	}

	fmt.Printf("Invitation sent to %s (Invitation ID: %d)\n", username, *invitation.ID)
	return nil
}
