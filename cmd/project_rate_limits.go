package cmd

import (
	"encoding/json"
	"errors"
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
			modifyProjectRateLimitsCommand(),
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

func modifyProjectRateLimitsCommand() *cli.Command {
	return &cli.Command{
		Name:  "modify",
		Usage: "Modify a project rate limit identified by rate limit ID for a given project ID",
		Flags: []cli.Flag{
			limitFlag,
			afterFlag,
			&cli.StringFlag{
				Name:     "project-id",
				Usage:    "ID of the project whose rate limits will be listed",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "rate-limit-id",
				Usage:    "ID of the rate limit to modify",
				Required: true,
			},
			&cli.IntFlag{
				Name:  "max-requests-per-1-minute",
				Usage: "The new max-requests-per-1-minute to set on the project rate limit",
			},
			&cli.IntFlag{
				Name:  "max-tokens-per-1-minute",
				Usage: "The new max-tokens-per-1-minute to set on the project rate limit",
			},
			&cli.IntFlag{
				Name:  "max-images-per-1-minute",
				Usage: "The new max-images-per-1-minute to set on the project rate limit",
			},
			&cli.IntFlag{
				Name:  "max-audio-megabytes-per-1-minute",
				Usage: "The new max-audio-megabytes-per-1-minute to set on the project rate limit",
			},
			&cli.IntFlag{
				Name:  "max-requests-per-1-day",
				Usage: "The new max-requests-per-1-day to set on the project rate limit",
			},
			&cli.IntFlag{
				Name:  "batch-1-day-max-input-tokens",
				Usage: "The new batch-1-day-max-input-tokens to set on the project rate limit",
			},
		},
		Action: modifyProjectRateLimit,
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

func printProjectRateLimitJson(projectRateLimit *openaiorgs.ProjectRateLimit) error {
	marshalled, err := json.Marshal(projectRateLimit)
	if err != nil {
		return wrapError("json marshalling error", err)
	}

	os.Stdout.Write(marshalled)

	return nil
}

func printProjectRateLimitTable(projectRateLimit *openaiorgs.ProjectRateLimit) error {
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
		Rows: [][]string{
			{
				projectRateLimit.ID,
				projectRateLimit.Model,
				strconv.Itoa(projectRateLimit.MaxRequestsPer1Minute),
				strconv.Itoa(projectRateLimit.MaxTokensPer1Minute),
				strconv.Itoa(projectRateLimit.MaxImagesPer1Minute),
				strconv.Itoa(projectRateLimit.MaxAudioMegabytesPer1Minute),
				strconv.Itoa(projectRateLimit.MaxRequestsPer1Day),
				strconv.Itoa(projectRateLimit.Batch1DayMaxInputTokens),
			},
		},
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

func validateModifyProjectRateLimitContext(c *cli.Context) error {
	if c.Int("max-requests-per-1-minute") != 0 {
		return nil
	}

	if c.Int("max-tokens-per-1-minute") != 0 {
		return nil
	}

	if c.Int("max-images-per-1-minute") != 0 {
		return nil
	}

	if c.Int("max-audio-megabytes-per-1-minute") != 0 {
		return nil
	}

	if c.Int("max-requests-per-1-day") != 0 {
		return nil
	}

	if c.Int("batch-1-day-max-input-tokens") != 0 {
		return nil
	}

	return errors.New("must set at least one field to modify for the project rate limit")
}

func modifyProjectRateLimit(c *cli.Context) error {
	client := newClient(c)

	if err := validateModifyProjectRateLimitContext(c); err != nil {
		return err
	}

	projectRateLimit, err := client.ModifyProjectRateLimit(
		c.Int("limit"),
		c.String("after"),
		c.String("project-id"),
		c.String("rate-limit-id"),
		openaiorgs.ProjectRateLimitRequestFields{
			MaxRequestsPer1Minute:       c.Int("max-requests-per-1-minute"),
			MaxTokensPer1Minute:         c.Int("max-tokens-per-1-minute"),
			MaxImagesPer1Minute:         c.Int("max-images-per-1-minute"),
			MaxAudioMegabytesPer1Minute: c.Int("max-audio-megabytes-per-1-minute"),
			MaxRequestsPer1Day:          c.Int("max-requests-per-1-day"),
			Batch1DayMaxInputTokens:     c.Int("batch-1-day-max-input-tokens"),
		},
	)
	if err != nil {
		return wrapError("modify project rate limit", err)
	}

	output := c.String("output")
	switch output {
	case OutputFormatJSON:
		return printProjectRateLimitJson(projectRateLimit)
	default:
		return printProjectRateLimitTable(projectRateLimit)
	}
}
