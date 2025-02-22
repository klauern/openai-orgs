package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/urfave/cli/v3"

	openaiorgs "github.com/klauern/openai-orgs"
)

func ProjectRateLimitsCommand() *cli.Command {
	return &cli.Command{
		Name:  "project-rate-limits",
		Usage: "Manage organization project rate limits",
		Commands: []*cli.Command{
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
			&cli.StringFlag{
				Name:     "project-id",
				Usage:    "ID of the project whose rate limits will be modified",
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
			strconv.FormatInt(projectRateLimit.MaxRequestsPer1Minute, 10),
			strconv.FormatInt(projectRateLimit.MaxTokensPer1Minute, 10),
			strconv.FormatInt(projectRateLimit.MaxImagesPer1Minute, 10),
			strconv.FormatInt(projectRateLimit.MaxAudioMegabytesPer1Minute, 10),
			strconv.FormatInt(projectRateLimit.MaxRequestsPer1Day, 10),
			strconv.FormatInt(projectRateLimit.Batch1DayMaxInputTokens, 10),
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
				strconv.FormatInt(projectRateLimit.MaxRequestsPer1Minute, 10),
				strconv.FormatInt(projectRateLimit.MaxTokensPer1Minute, 10),
				strconv.FormatInt(projectRateLimit.MaxImagesPer1Minute, 10),
				strconv.FormatInt(projectRateLimit.MaxAudioMegabytesPer1Minute, 10),
				strconv.FormatInt(projectRateLimit.MaxRequestsPer1Day, 10),
				strconv.FormatInt(projectRateLimit.Batch1DayMaxInputTokens, 10),
			},
		},
	}

	printTableData(data)
	return nil
}

func listProjectRateLimits(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	limit := int(cmd.Int("limit"))
	projectRateLimits, err := client.ListProjectRateLimits(
		limit,
		cmd.String("after"),
		cmd.String("project-id"),
	)
	if err != nil {
		return wrapError("list project rate limits", err)
	}

	switch cmd.String("output") {
	case "json":
		return printProjectRateLimitsJson(projectRateLimits)
	default:
		return printProjectRateLimitsTable(projectRateLimits)
	}
}

func validateModifyProjectRateLimitContext(ctx context.Context, cmd *cli.Command) error {
	if cmd.String("project-id") == "" {
		return errors.New("project-id is required")
	}
	if cmd.String("rate-limit-id") == "" {
		return errors.New("rate-limit-id is required")
	}

	// At least one rate limit field must be set
	if cmd.Int("max-requests-per-1-minute") == 0 &&
		cmd.Int("max-tokens-per-1-minute") == 0 &&
		cmd.Int("max-images-per-1-minute") == 0 &&
		cmd.Int("max-audio-megabytes-per-1-minute") == 0 &&
		cmd.Int("max-requests-per-1-day") == 0 &&
		cmd.Int("batch-1-day-max-input-tokens") == 0 {
		return errors.New("must set at least one rate limit field to modify")
	}

	return nil
}

func modifyProjectRateLimit(ctx context.Context, cmd *cli.Command) error {
	if err := validateModifyProjectRateLimitContext(ctx, cmd); err != nil {
		return err
	}

	client := newClient(ctx, cmd)

	fields := openaiorgs.ProjectRateLimitRequestFields{
		MaxRequestsPer1Minute:       int64(cmd.Int("max-requests-per-1-minute")),
		MaxTokensPer1Minute:         int64(cmd.Int("max-tokens-per-1-minute")),
		MaxImagesPer1Minute:         int64(cmd.Int("max-images-per-1-minute")),
		MaxAudioMegabytesPer1Minute: int64(cmd.Int("max-audio-megabytes-per-1-minute")),
		MaxRequestsPer1Day:          int64(cmd.Int("max-requests-per-1-day")),
		Batch1DayMaxInputTokens:     int64(cmd.Int("batch-1-day-max-input-tokens")),
	}

	projectRateLimit, err := client.ModifyProjectRateLimit(
		cmd.String("project-id"),
		cmd.String("rate-limit-id"),
		fields,
	)
	if err != nil {
		return wrapError("modify project rate limit", err)
	}

	switch cmd.String("output") {
	case "json":
		return printProjectRateLimitJson(projectRateLimit)
	default:
		return printProjectRateLimitTable(projectRateLimit)
	}
}
