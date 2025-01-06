package cmd

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"

	openaiorgs "github.com/klauern/openai-orgs"
)

func ProjectRateLimitsCommand() *cli.Command {
	return &cli.Command{
		Name:  "project-rate-limits",
		Usage: "Manage organization project rate limits",
		Subcommands: []*cli.Command{
			listProjectRateLimitsCommand(),
		},
	}
}

func listProjectRateLimitsCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all project rate limits for a given project ID",
		Flags: []cli.Flag{
			limitFlag,
			afterFlag,
			&cli.StringFlag{
				Name:     "project-id",
				Usage:    "ID of the project whose rate limits will be listed",
				Required: true,
			},
		},
		Action: listProjectRateLimits,
	}
}

func printProjectRateLimitsJson(projectRateLimits *openaiorgs.ListResponse[openaiorgs.ProjectRateLimit]) error {
	marshalled, err := json.Marshal(projectRateLimits.Data)
	if err != nil {
		return wrapError("json marshalling error", err)
	}

	os.Stdout.Write(marshalled)

	return nil
}

func printProjectRateLimitsTable(projectRateLimits *openaiorgs.ListResponse[openaiorgs.ProjectRateLimit]) error {
	data := TableData{
		Headers: []string{
			"ID",
			"Model",
			"Max Requests Per 1 Minute",
			"Max Tokens Per 1 Minute",
			"Max Images Per 1 Minute",
			"Max Audio Megabytes Per 1 Minute",
			"Max Requests Per 1 Day",
			"Batch 1 Day Max Input Tokens",
		},
		Rows: make(
			[][]string,
			len(projectRateLimits.Data),
		),
	}

	for i, projectRateLimit := range projectRateLimits.Data {
		data.Rows[i] = []string{
			projectRateLimit.ID,
			projectRateLimit.Model,
			strconv.Itoa(projectRateLimit.MaxRequestsPer1Minute),
			strconv.Itoa(projectRateLimit.MaxTokensPer1Minute),
			strconv.Itoa(projectRateLimit.MaxImagesPer1Minute),
			strconv.Itoa(projectRateLimit.MaxAudioMegabytesPer1Minute),
			strconv.Itoa(projectRateLimit.MaxRequestsPer1Day),
			strconv.Itoa(projectRateLimit.Batch1DayMaxInputTokens),
		}
	}

	printTableData(data)
	return nil
}

func listProjectRateLimits(c *cli.Context) error {
	client := newClient(c)

	projectRateLimits, err := client.ListProjectRateLimits(
		c.Int("limit"),
		c.String("after"),
		c.String("project-id"),
	)
	if err != nil {
		return wrapError("list project rate limits", err)
	}

	output := c.String("output")
	switch output {
	case OutputFormatJSON:
		return printProjectRateLimitsJson(projectRateLimits)
	default:
		return printProjectRateLimitsTable(projectRateLimits)
	}
}