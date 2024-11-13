package cmd

import (
	"fmt"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/urfave/cli/v2"
)

func InvitesCommand() *cli.Command {
	return &cli.Command{
		Name:  "invites",
		Usage: "Manage organization invites",
		Subcommands: []*cli.Command{
			listInvitesCommand(),
			createInviteCommand(),
			retrieveInviteCommand(),
			deleteInviteCommand(),
		},
	}
}

func listInvitesCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "List all invites",
		Action: listInvites,
	}
}

func createInviteCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new invite",
		Flags: []cli.Flag{
			emailFlag,
			roleFlag,
		},
		Action: createInvite,
	}
}

func deleteInviteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete an invite",
		Flags: []cli.Flag{
			idFlag,
		},
		Action: deleteInvite,
	}
}

func retrieveInviteCommand() *cli.Command {
	return &cli.Command{
		Name:  "retrieve",
		Usage: "Retrieve an invite",
		Flags: []cli.Flag{
			idFlag,
		},
		Action: retrieveInvite,
	}
}

func listInvites(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	invites, err := client.ListInvites()
	if err != nil {
		return fmt.Errorf("failed to list invites: %w", err)
	}

	data := TableData{
		Headers: []string{"ID", "Email", "Role", "Status", "Created At", "Expires At", "Accepted At"},
		Rows:    make([][]string, len(invites)),
	}

	for i, invite := range invites {
		acceptedAt := "N/A"
		if invite.AcceptedAt != nil {
			acceptedAt = invite.AcceptedAt.String()
		}
		data.Rows[i] = []string{
			invite.ID,
			invite.Email,
			invite.Role,
			invite.Status,
			invite.CreatedAt.String(),
			invite.ExpiresAt.String(),
			acceptedAt,
		}
	}

	printTable(data.Headers, data.Rows)

	return nil
}

func createInvite(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	email := c.String("email")

	invite, err := client.CreateInvite(email, c.String("role"))
	if err != nil {
		return fmt.Errorf("failed to create invite: %w", err)
	}

	fmt.Printf("Invite created: ID: %s, Email: %s, Role: %s, Status: %s\n", invite.ID, invite.Email, invite.Role, invite.Status)
	return nil
}

func deleteInvite(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	id := c.String("id")

	err := client.DeleteInvite(id)
	if err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	fmt.Printf("Invite with ID %s has been deleted\n", id)
	return nil
}

func retrieveInvite(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	id := c.String("id")

	invite, err := client.RetrieveInvite(id)
	if err != nil {
		return fmt.Errorf("failed to retrieve invite: %w", err)
	}

	acceptedAt := "N/A"
	if invite.AcceptedAt != nil {
		acceptedAt = invite.AcceptedAt.String()
	}

	fmt.Printf("Invite retrieved: ID: %s, Email: %s, Role: %s, Status: %s, Created At: %s, Expires At: %s, Accepted At: %s\n",
		invite.ID, invite.Email, invite.Role, invite.Status, invite.CreatedAt.String(), invite.ExpiresAt.String(), acceptedAt)
	return nil
}
