package cmd

import (
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
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Limit the number of results",
			},
			&cli.StringFlag{
				Name:  "after",
				Usage: "Return results after this ID",
			},
			&cli.StringFlag{
				Name:  "before",
				Usage: "Return results before this ID",
			},
			&cli.StringFlag{
				Name:  "start-date",
				Usage: "Start date for the query (RFC3339 format)",
			},
			&cli.StringFlag{
				Name:  "end-date",
				Usage: "End date for the query (RFC3339 format)",
			},
		},
		Action: listAuditLogs,
	}
}

func listAuditLogs(c *cli.Context) error {
	client := openaiorgs.NewClient("https://api.openai.com/v1", os.Getenv("OPENAI_API_KEY"))

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

	// Print the audit logs
	for _, log := range response.Data {
		fmt.Printf("ID: %s, Type: %s, Timestamp: %s\n", log.ID, log.Type, log.Timestamp)
		payload, err := openaiorgs.ParseAuditLogPayload(&log)
		if err != nil {
			fmt.Printf("Error parsing payload: %v\n", err)
		} else {
			fmt.Printf("Payload: %+v\n", payload)
		}
		fmt.Println("---")
	}

	return nil
}
