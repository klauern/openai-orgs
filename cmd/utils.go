package cmd

import (
	"fmt"
	"strings"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/urfave/cli/v2"
)

// Constants should be grouped at the top
const (
	DefaultBaseURL = openaiorgs.DefaultBaseURL

	// Common output formats
	OutputFormatPretty = "pretty"
	OutputFormatJSON   = "json"
)

// Common flag definitions grouped together
var (
	projectIDFlag = &cli.StringFlag{
		Name:     "project-id",
		Usage:    "ID of the project",
		Required: true,
	}

	idFlag = &cli.StringFlag{
		Name:     "id",
		Usage:    "Resource ID",
		Required: true,
	}

	limitFlag = &cli.IntFlag{
		Name:  "limit",
		Usage: "Limit the number of items returned",
	}

	afterFlag = &cli.StringFlag{
		Name:  "after",
		Usage: "Return items after this ID",
	}

	beforeFlag = &cli.StringFlag{
		Name:  "before",
		Usage: "Return results before this ID",
	}

	emailFlag = &cli.StringFlag{
		Name:     "email",
		Usage:    "Email address",
		Required: true,
	}

	roleFlag = &cli.StringFlag{
		Name:     "role",
		Usage:    "Role (e.g., owner, member)",
		Required: true,
	}

	nameFlag = &cli.StringFlag{
		Name:     "name",
		Usage:    "Resource name",
		Required: true,
	}

	ValidOutputFormats = map[string]bool{
		OutputFormatPretty: true,
		OutputFormatJSON:   true,
	}
)

// Interfaces and types grouped together
type CommandProvider interface {
	Command() *cli.Command
	Subcommands() []*cli.Command
}

type TableData struct {
	Headers []string
	Rows    [][]string
}

// Base type and its methods grouped together
type BaseCommand struct {
	Name     string
	Usage    string
	Commands []*cli.Command
}

func (b *BaseCommand) Command() *cli.Command {
	return &cli.Command{
		Name:        b.Name,
		Usage:       b.Usage,
		Subcommands: b.Commands,
	}
}

func (b *BaseCommand) Subcommands() []*cli.Command {
	return b.Commands
}

func newClient(c *cli.Context) *openaiorgs.Client {
	return openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))
}

func printTable(headers []string, rows [][]string) {
	// Print headers
	fmt.Println(strings.Join(headers, " | "))
	fmt.Println(strings.Repeat("-", len(strings.Join(headers, " | "))))

	// Print rows
	for _, row := range rows {
		fmt.Println(strings.Join(row, " | "))
	}
}

func printTableData(data TableData) {
	// Print headers
	fmt.Println(strings.Join(data.Headers, " | "))
	fmt.Println(strings.Repeat("-", len(strings.Join(data.Headers, " | "))))

	// Print rows
	for _, row := range data.Rows {
		fmt.Println(strings.Join(row, " | "))
	}
}

func wrapError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("failed to %s: %w", operation, err)
}
