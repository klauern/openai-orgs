package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func InvitesCommand() *cli.Command {
	return &cli.Command{
		Name:  "invites",
		Usage: "Manage organization invites",
		Commands: []*cli.Command{
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
		Name:   "create",
		Usage:  "Create a new invite",
		Flags:  []cli.Flag{emailFlag, roleFlag},
		Action: createInvite,
	}
}

func retrieveInviteCommand() *cli.Command {
	return &cli.Command{
		Name:   "retrieve",
		Usage:  "Retrieve a specific invite",
		Flags:  []cli.Flag{idFlag},
		Action: retrieveInvite,
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

func listInvites(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	invites, err := client.ListInvites()
	if err != nil {
		return wrapError("list invites", err)
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

	printTableData(data)
	return nil
}

func createInvite(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	invite, err := client.CreateInvite(
		cmd.String("email"),
		cmd.String("role"),
	)
	if err != nil {
		return wrapError("create invite", err)
	}

	fmt.Printf("Invite created:\n")
	fmt.Printf("ID: %s\nEmail: %s\nRole: %s\nCreated At: %s\nExpires At: %s\n",
		invite.ID,
		invite.Email,
		invite.Role,
		invite.CreatedAt.String(),
		invite.ExpiresAt.String(),
	)

	return nil
}

func retrieveInvite(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	invite, err := client.RetrieveInvite(cmd.String("id"))
	if err != nil {
		return wrapError("retrieve invite", err)
	}

	fmt.Printf("Invite details:\n")
	fmt.Printf("ID: %s\nEmail: %s\nRole: %s\nCreated At: %s\nExpires At: %s\n",
		invite.ID,
		invite.Email,
		invite.Role,
		invite.CreatedAt.String(),
		invite.ExpiresAt.String(),
	)

	return nil
}

func deleteInvite(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	err := client.DeleteInvite(cmd.String("id"))
	if err != nil {
		return wrapError("delete invite", err)
	}

	fmt.Printf("Invite %s deleted successfully\n", cmd.String("id"))
	return nil
}
