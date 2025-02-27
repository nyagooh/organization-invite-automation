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

