package cmd

import (
	"fmt"
	"strings"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/urfave/cli/v2"
)

func ProjectsCommand() *cli.Command {
	return &cli.Command{
		Name:  "projects",
		Usage: "Manage organization projects",
		Subcommands: []*cli.Command{
			listProjectsCommand(),
			createProjectCommand(),
			retrieveProjectCommand(),
			modifyProjectCommand(),
			archiveProjectCommand(),
		},
	}
}

func listProjectsCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all projects",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Limit the number of projects returned",
			},
			&cli.StringFlag{
				Name:  "after",
				Usage: "Return projects after this ID",
			},
			&cli.BoolFlag{
				Name:  "include-archived",
				Usage: "Include archived projects in the list",
			},
		},
		Action: listProjects,
	}
}

func createProjectCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new project",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Name of the project to create",
				Required: true,
			},
		},
		Action: createProject,
	}
}

func retrieveProjectCommand() *cli.Command {
	return &cli.Command{
		Name:  "retrieve",
		Usage: "Retrieve a specific project",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "ID of the project to retrieve",
				Required: true,
			},
		},
		Action: retrieveProject,
	}
}

func modifyProjectCommand() *cli.Command {
	return &cli.Command{
		Name:  "modify",
		Usage: "Modify a project",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "ID of the project to modify",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "New name for the project",
				Required: true,
			},
		},
		Action: modifyProject,
	}
}

func archiveProjectCommand() *cli.Command {
	return &cli.Command{
		Name:  "archive",
		Usage: "Archive a project",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "ID of the project to archive",
				Required: true,
			},
		},
		Action: archiveProject,
	}
}

func listProjects(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	limit := c.Int("limit")
	after := c.String("after")
	includeArchived := c.Bool("include-archived")

	projects, err := client.ListProjects(limit, after, includeArchived)
	if err != nil {
		return fmt.Errorf("failed to list projects: %w", err)
	}

	fmt.Println("ID | Name | Created At | Archived At | Status")
	fmt.Println(strings.Repeat("-", 80))
	for _, project := range projects.Data {
		archivedAt := "N/A"
		if project.ArchivedAt != nil {
			archivedAt = project.ArchivedAt.String()
		}
		fmt.Printf("%s | %s | %s | %s | %s\n",
			project.ID,
			project.Name,
			project.CreatedAt.String(),
			archivedAt,
			project.Status,
		)
	}

	return nil
}

func createProject(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	name := c.String("name")

	project, err := client.CreateProject(name)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	fmt.Printf("Project created:\n")
	fmt.Printf("ID: %s\nName: %s\nCreated At: %s\nStatus: %s\n",
		project.ID,
		project.Name,
		project.CreatedAt.String(),
		project.Status,
	)

	return nil
}

func retrieveProject(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	id := c.String("id")

	project, err := client.RetrieveProject(id)
	if err != nil {
		return fmt.Errorf("failed to retrieve project: %w", err)
	}

	fmt.Printf("Project details:\n")
	fmt.Printf("ID: %s\nName: %s\nCreated At: %s\nStatus: %s\n",
		project.ID,
		project.Name,
		project.CreatedAt.String(),
		project.Status,
	)
	if project.ArchivedAt != nil {
		fmt.Printf("Archived At: %s\n", project.ArchivedAt.String())
	}

	return nil
}

func modifyProject(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	id := c.String("id")
	name := c.String("name")

	project, err := client.ModifyProject(id, name)
	if err != nil {
		return fmt.Errorf("failed to modify project: %w", err)
	}

	fmt.Printf("Project modified:\n")
	fmt.Printf("ID: %s\nNew Name: %s\nCreated At: %s\nStatus: %s\n",
		project.ID,
		project.Name,
		project.CreatedAt.String(),
		project.Status,
	)

	return nil
}

func archiveProject(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))

	id := c.String("id")

	project, err := client.ArchiveProject(id)
	if err != nil {
		return fmt.Errorf("failed to archive project: %w", err)
	}

	fmt.Printf("Project archived:\n")
	fmt.Printf("ID: %s\nName: %s\nCreated At: %s\nArchived At: %s\nStatus: %s\n",
		project.ID,
		project.Name,
		project.CreatedAt.String(),
		project.ArchivedAt.String(),
		project.Status,
	)

	return nil
}
