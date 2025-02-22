package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func ProjectsCommand() *cli.Command {
	return &cli.Command{
		Name:  "projects",
		Usage: "Manage organization projects",
		Commands: []*cli.Command{
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
			limitFlag,
			afterFlag,
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
			nameFlag,
		},
		Action: createProject,
	}
}

func retrieveProjectCommand() *cli.Command {
	return &cli.Command{
		Name:  "retrieve",
		Usage: "Retrieve a specific project",
		Flags: []cli.Flag{
			idFlag,
		},
		Action: retrieveProject,
	}
}

func modifyProjectCommand() *cli.Command {
	return &cli.Command{
		Name:  "modify",
		Usage: "Modify a project",
		Flags: []cli.Flag{
			idFlag,
			nameFlag,
		},
		Action: modifyProject,
	}
}

func archiveProjectCommand() *cli.Command {
	return &cli.Command{
		Name:  "archive",
		Usage: "Archive a project",
		Flags: []cli.Flag{
			idFlag,
		},
		Action: archiveProject,
	}
}

func listProjects(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	limit := int(cmd.Int("limit"))
	projects, err := client.ListProjects(
		limit,
		cmd.String("after"),
		cmd.Bool("include-archived"),
	)
	if err != nil {
		return wrapError("list projects", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Created At", "Archived At", "Status"},
		Rows:    make([][]string, len(projects.Data)),
	}

	for i, project := range projects.Data {
		archivedAt := "N/A"
		if project.ArchivedAt != nil {
			archivedAt = project.ArchivedAt.String()
		}
		data.Rows[i] = []string{
			project.ID,
			project.Name,
			project.CreatedAt.String(),
			archivedAt,
			project.Status,
		}
	}

	printTableData(data)
	return nil
}

func createProject(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	name := cmd.String("name")

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

func retrieveProject(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	id := cmd.String("id")

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

func modifyProject(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	id := cmd.String("id")
	name := cmd.String("name")

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

func archiveProject(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	id := cmd.String("id")

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
