package cmd

import (
	"fmt"
	"strings"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/urfave/cli/v2"
)

func ProjectUsersCommand() *cli.Command {
	return &cli.Command{
		Name:  "project-users",
		Usage: "Manage project users",
		Subcommands: []*cli.Command{
			listProjectUsersCommand(),
			createProjectUserCommand(),
			retrieveProjectUserCommand(),
			modifyProjectUserCommand(),
			deleteProjectUserCommand(),
		},
	}
}

func listProjectUsersCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all project users",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "project-id",
				Usage:    "ID of the project",
				Required: true,
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Limit the number of users returned",
			},
			&cli.StringFlag{
				Name:  "after",
				Usage: "Return users after this ID",
			},
		},
		Action: listProjectUsers,
	}
}

func createProjectUserCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new project user",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
			roleFlag,
		},
		Action: createProjectUser,
	}
}

func retrieveProjectUserCommand() *cli.Command {
	return &cli.Command{
		Name:  "retrieve",
		Usage: "Retrieve a specific project user",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
		},
		Action: retrieveProjectUser,
	}
}

func modifyProjectUserCommand() *cli.Command {
	return &cli.Command{
		Name:  "modify",
		Usage: "Modify a project user's role",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
			roleFlag,
		},
		Action: modifyProjectUser,
	}
}

func deleteProjectUserCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a project user",
		Flags: []cli.Flag{
			projectIDFlag,
			idFlag,
		},
		Action: deleteProjectUser,
	}
}

func listProjectUsers(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	projectID := c.String("project-id")
	limit := c.Int("limit")
	after := c.String("after")

	users, err := client.ListProjectUsers(projectID, limit, after)
	if err != nil {
		return fmt.Errorf("failed to list project users: %w", err)
	}

	fmt.Println("ID | Name | Email | Role | Added At")
	fmt.Println(strings.Repeat("-", 80))
	for _, user := range users.Data {
		fmt.Printf("%s | %s | %s | %s | %s\n",
			user.ID,
			user.Name,
			user.Email,
			user.Role,
			user.AddedAt.String(),
		)
	}

	return nil
}

func createProjectUser(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	projectID := c.String("project-id")
	userID := c.String("user-id")

	user, err := client.CreateProjectUser(projectID, userID, c.String("role"))
	if err != nil {
		return fmt.Errorf("failed to create project user: %w", err)
	}

	fmt.Printf("User added to project:\n")
	fmt.Printf("ID: %s\nName: %s\nEmail: %s\nRole: %s\nAdded At: %s\n",
		user.ID,
		user.Name,
		user.Email,
		user.Role,
		user.AddedAt.String(),
	)

	return nil
}

func retrieveProjectUser(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	projectID := c.String("project-id")
	userID := c.String("user-id")

	user, err := client.RetrieveProjectUser(projectID, userID)
	if err != nil {
		return fmt.Errorf("failed to retrieve project user: %w", err)
	}

	fmt.Printf("Project user details:\n")
	fmt.Printf("ID: %s\nName: %s\nEmail: %s\nRole: %s\nAdded At: %s\n",
		user.ID,
		user.Name,
		user.Email,
		user.Role,
		user.AddedAt.String(),
	)

	return nil
}

func modifyProjectUser(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	projectID := c.String("project-id")
	userID := c.String("user-id")

	user, err := client.ModifyProjectUser(projectID, userID, c.String("role"))
	if err != nil {
		return fmt.Errorf("failed to modify project user role: %w", err)
	}

	fmt.Printf("Project user updated:\n")
	fmt.Printf("ID: %s\nName: %s\nEmail: %s\nRole: %s\nAdded At: %s\n",
		user.ID,
		user.Name,
		user.Email,
		user.Role,
		user.AddedAt.String(),
	)

	return nil
}

func deleteProjectUser(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	projectID := c.String("project-id")
	userID := c.String("user-id")

	err := client.DeleteProjectUser(projectID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete project user: %w", err)
	}

	fmt.Printf("User with ID %s has been removed from project %s\n", userID, projectID)
	return nil
}
