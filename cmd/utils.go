package cmd

import (
	"fmt"
	"strings"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/urfave/cli/v2"
)

// Common flag definitions
var (
	projectIDFlag = &cli.StringFlag{
		Name:     "project-id",
		Usage:    "ID of the project",
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

	idFlag = &cli.StringFlag{
		Name:     "id",
		Usage:    "ID of the resource",
		Required: true,
	}
)

// newClient creates a new OpenAI client from context
func newClient(c *cli.Context) *openaiorgs.Client {
	return openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))
}

// printTable prints a formatted table with headers and rows
func printTable(headers []string, rows [][]string) {
	// Print headers
	fmt.Println(strings.Join(headers, " | "))
	fmt.Println(strings.Repeat("-", len(strings.Join(headers, " | "))))

	// Print rows
	for _, row := range rows {
		fmt.Println(strings.Join(row, " | "))
	}
}
