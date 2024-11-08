package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/urfave/cli/v2"
)

func AuditLogsCommand() *cli.Command {
	return &cli.Command{
		Name:  "audit-logs",
		Usage: "List audit logs",
		Flags: []cli.Flag{
			limitFlag,
			afterFlag,
			&cli.StringFlag{
				Name:  "start-date",
				Usage: "Start date for the query (RFC3339 format)",
			},
			&cli.StringFlag{
				Name:  "end-date",
				Usage: "End date for the query (RFC3339 format)",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output format (pretty, json, jsonl)",
				Value:   "pretty",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Show verbose output",
			},
		},
		Action: listAuditLogs,
	}
}

func listAuditLogs(c *cli.Context) error {
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, c.String("api-key"))
	verbose := c.Bool("verbose")
	outputFormat := c.String("output")

	params := &openaiorgs.AuditLogListParams{
		Limit:  c.Int("limit"),
		After:  c.String("after"),
		Before: c.String("before"),
	}

	if startDate := c.String("start-date"); startDate != "" {
		parsedStartDate, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			return fmt.Errorf("invalid start-date format: %w", err)
		}
		params.StartDate = parsedStartDate
	}

	if endDate := c.String("end-date"); endDate != "" {
		parsedEndDate, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			return fmt.Errorf("invalid end-date format: %w", err)
		}
		params.EndDate = parsedEndDate
	}

	response, err := client.ListAuditLogs(params)
	if err != nil {
		return fmt.Errorf("failed to list audit logs: %w", err)
	}

	switch outputFormat {
	case "json":
		return outputJSON(response, verbose)
	case "jsonl":
		return outputJSONL(response, verbose)
	case "pretty":
		return outputPretty(response, verbose)
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
}

func outputJSON(response *openaiorgs.ListResponse[openaiorgs.AuditLog], verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

func outputJSONL(response *openaiorgs.ListResponse[openaiorgs.AuditLog], verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	// First line: metadata if verbose
	if verbose {
		meta := struct {
			Total   int    `json:"total"`
			HasMore bool   `json:"has_more"`
			FirstID string `json:"first_id"`
			LastID  string `json:"last_id"`
		}{
			Total:   len(response.Data),
			HasMore: response.HasMore,
			FirstID: response.FirstID,
			LastID:  response.LastID,
		}
		if err := encoder.Encode(meta); err != nil {
			return err
		}
	}

	// Output each log entry on its own line
	for _, log := range response.Data {
		if err := encoder.Encode(log); err != nil {
			return err
		}
	}
	return nil
}

func outputPretty(response *openaiorgs.ListResponse[openaiorgs.AuditLog], verbose bool) error {
	if verbose {
		fmt.Printf("=== Audit Log Summary ===\n")
		fmt.Printf("Total logs: %d\nHas more: %v\nFirst ID: %s\nLast ID: %s\n\n",
			len(response.Data), response.HasMore, response.FirstID, response.LastID)
	}

	for _, log := range response.Data {
		fmt.Printf("=== Audit Log Entry ===\n")
		fmt.Printf("ID:        %s\n", log.ID)
		fmt.Printf("Type:      %s\n", log.Type)
		fmt.Printf("Timestamp: %s\n", log.Timestamp)

		if verbose {
			fmt.Printf("\nActor Details:\n")
			fmt.Printf("  ID:   %s\n", log.Actor.ID)
			fmt.Printf("  Name: %s\n", log.Actor.Name)
			fmt.Printf("  Type: %s\n", log.Actor.Type)

			fmt.Printf("\nEvent Details:\n")
			fmt.Printf("  ID:     %s\n", log.Event.ID)
			fmt.Printf("  Type:   %s\n", log.Event.Type)
			fmt.Printf("  Action: %s\n", log.Event.Action)
			if log.Event.Auth.Type != "" {
				fmt.Printf("  Auth:    %s (%s)\n", log.Event.Auth.Type, log.Event.Auth.Transport)
			}
		}

		payload, err := openaiorgs.ParseAuditLogPayload(&log)
		if err != nil {
			fmt.Printf("\nPayload Error: %v\n", err)
		} else {
			fmt.Printf("\nPayload Details:\n")
			switch p := payload.(type) {
			case *openaiorgs.LoginSucceeded:
				fmt.Printf("  Login Success\n")
				// Add specific login success details if available
			// Add other payload types as needed
			default:
				fmt.Printf("  %#v\n", p)
			}
		}
		fmt.Println("\n---")
	}

	return nil
}
