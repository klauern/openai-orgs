package cmd

import (
	"fmt"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/urfave/cli/v2"
)

func UsersCommand() *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "Manage organization users",
		Subcommands: []*cli.Command{
			listUsersCommand(),
			retrieveUserCommand(),
			deleteUserCommand(),
			modifyUserRoleCommand(),
		},
	}
}

func listUsersCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all users",
		Flags: []cli.Flag{
			limitFlag,
			afterFlag,
		},
		Action: listUsers,
	}
}

func retrieveUserCommand() *cli.Command {
	return &cli.Command{
		Name:  "retrieve",
		Usage: "Retrieve a specific user",
		Flags: []cli.Flag{
			idFlag,
		},
		Action: retrieveUser,
	}
}

func deleteUserCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a user",
		Flags: []cli.Flag{
			idFlag,
		},
		Action: deleteUser,
	}
}

func modifyUserRoleCommand() *cli.Command {
	return &cli.Command{
		Name:  "modify-role",
		Usage: "Modify a user's role",
		Flags: []cli.Flag{
			idFlag,
			roleFlag,
		},
		Action: modifyUserRole,
	}
}

func listUsers(c *cli.Context) error {
	client := newClient(c)

	users, err := client.ListUsers(
		c.Int("limit"),
		c.String("after"),
	)
	if err != nil {
		return wrapError("list users", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Email", "Role", "Added At"},
		Rows:    make([][]string, len(users.Data)),
	}

	for i, user := range users.Data {
		data.Rows[i] = []string{
			user.ID,
			user.Name,
			user.Email,
			user.Role,
			user.AddedAt.String(),
		}
	}

	printTableData(data)
	return nil
}

func retrieveUser(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	id := c.String("id")

	user, err := client.RetrieveUser(id)
	if err != nil {
		return fmt.Errorf("failed to retrieve user: %w", err)
	}

	fmt.Printf("User details:\n")
	fmt.Printf("ID: %s\nName: %s\nEmail: %s\nRole: %s\nAdded At: %s\n",
		user.ID,
		user.Name,
		user.Email,
		user.Role,
		user.AddedAt.String(),
	)

	return nil
}

func deleteUser(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	id := c.String("id")

	err := client.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	fmt.Printf("User with ID %s has been deleted\n", id)
	return nil
}

func modifyUserRole(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	id := c.String("id")
	role := c.String("role")

	err := client.ModifyUserRole(id, role)
	if err != nil {
		return fmt.Errorf("failed to modify user role: %w", err)
	}

	fmt.Printf("User with ID %s has been updated with role: %s\n", id, role)
	return nil
}
