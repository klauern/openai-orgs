package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
		if params.EffectiveAt == nil {
			params.EffectiveAt = &struct {
				Gte int64 `json:"gte,omitempty"`
				Gt  int64 `json:"gt,omitempty"`
				Lte int64 `json:"lte,omitempty"`
				Lt  int64 `json:"lt,omitempty"`
			}{}
		}
		params.EffectiveAt.Gte = parsedStartDate.Unix()
	}

	if endDate := c.String("end-date"); endDate != "" {
		parsedEndDate, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			return fmt.Errorf("invalid end-date format: %w", err)
		}
		if params.EffectiveAt == nil {
			params.EffectiveAt = &struct {
				Gte int64 `json:"gte,omitempty"`
				Gt  int64 `json:"gt,omitempty"`
				Lte int64 `json:"lte,omitempty"`
				Lt  int64 `json:"lt,omitempty"`
			}{}
		}
		params.EffectiveAt.Lte = parsedEndDate.Unix()
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
		fmt.Printf("Effective: %s\n", log.EffectiveAt.Time().Format(time.RFC3339))

		if verbose {
			fmt.Printf("\nActor Details:\n")
			fmt.Printf("  Type: %s\n", log.Actor.Type)
			if log.Actor.Session != nil {
				fmt.Printf("  User ID:    %s\n", log.Actor.Session.User.ID)
				fmt.Printf("  User Email: %s\n", log.Actor.Session.User.Email)
				fmt.Printf("  IP:         %s\n", log.Actor.Session.IPAddress)
				fmt.Printf("  User Agent: %s\n", log.Actor.Session.UserAgent)
			}
			if log.Actor.APIKey != nil {
				fmt.Printf("  API Key Type: %s\n", log.Actor.APIKey.Type)
				fmt.Printf("  User ID:      %s\n", log.Actor.APIKey.User.ID)
				fmt.Printf("  User Email:   %s\n", log.Actor.APIKey.User.Email)
			}
		}

		if log.Details != nil {
			fmt.Printf("\nPayload Details:\n")
			switch details := log.Details.(type) {
			case *openaiorgs.APIKeyCreated:
				fmt.Printf("  API Key created with ID: %s\n", details.ID)
				if len(details.Data.Scopes) > 0 {
					fmt.Printf("  Scopes: %s\n", strings.Join(details.Data.Scopes, ", "))
				}
			case *openaiorgs.APIKeyUpdated:
				fmt.Printf("  API Key updated with ID: %s\n", details.ID)
				if len(details.ChangesRequested.Scopes) > 0 {
					fmt.Printf("  New scopes: %s\n", strings.Join(details.ChangesRequested.Scopes, ", "))
				}
			case *openaiorgs.APIKeyDeleted:
				fmt.Printf("  API Key deleted with ID: %s\n", details.ID)
			case *openaiorgs.InviteSent:
				fmt.Printf("  Invite sent with ID: %s\n  Email: %s\n",
					details.ID, details.Data.Email)
			case *openaiorgs.InviteAccepted:
				fmt.Printf("  Invite accepted with ID: %s\n", details.ID)
			case *openaiorgs.InviteDeleted:
				fmt.Printf("  Invite deleted with ID: %s\n", details.ID)
			case *openaiorgs.LoginFailed:
				fmt.Printf("  Login failed\n  Error code: %s\n  Error message: %s\n",
					details.ErrorCode, details.ErrorMessage)
			case *openaiorgs.LogoutFailed:
				fmt.Printf("  Logout failed\n  Error code: %s\n  Error message: %s\n",
					details.ErrorCode, details.ErrorMessage)
			case *openaiorgs.OrganizationUpdated:
				fmt.Printf("  Organization updated with ID: %s\n", details.ID)
				if details.ChangesRequested.Name != "" {
					fmt.Printf("  New name: %s\n", details.ChangesRequested.Name)
				}
			case *openaiorgs.ProjectCreated:
				fmt.Printf("  Project created with ID: %s\n  Name: %s\n  Title: %s\n",
					details.ID, details.Data.Name, details.Data.Title)
			case *openaiorgs.ProjectUpdated:
				fmt.Printf("  Project updated with ID: %s\n  New title: %s\n",
					details.ID, details.ChangesRequested.Title)
			case *openaiorgs.ProjectArchived:
				fmt.Printf("  Project archived with ID: %s\n", details.ID)
			case *openaiorgs.RateLimitUpdated:
				fmt.Printf("  Rate limit updated with ID: %s\n", details.ID)
				changes := details.ChangesRequested
				if changes.MaxRequestsPer1Minute > 0 {
					fmt.Printf("  Max requests per minute: %d\n", changes.MaxRequestsPer1Minute)
				}
				if changes.MaxTokensPer1Minute > 0 {
					fmt.Printf("  Max tokens per minute: %d\n", changes.MaxTokensPer1Minute)
				}
				if changes.MaxImagesPer1Minute > 0 {
					fmt.Printf("  Max images per minute: %d\n", changes.MaxImagesPer1Minute)
				}
				if changes.MaxAudioMegabytesPer1Minute > 0 {
					fmt.Printf("  Max audio MB per minute: %d\n", changes.MaxAudioMegabytesPer1Minute)
				}
				if changes.MaxRequestsPer1Day > 0 {
					fmt.Printf("  Max requests per day: %d\n", changes.MaxRequestsPer1Day)
				}
				if changes.Batch1DayMaxInputTokens > 0 {
					fmt.Printf("  Batch max input tokens per day: %d\n", changes.Batch1DayMaxInputTokens)
				}
			case *openaiorgs.RateLimitDeleted:
				fmt.Printf("  Rate limit deleted with ID: %s\n", details.ID)
			case *openaiorgs.ServiceAccountCreated:
				fmt.Printf("  Service account created with ID: %s\n  Role: %s\n",
					details.ID, details.Data.Role)
			case *openaiorgs.ServiceAccountUpdated:
				fmt.Printf("  Service account updated with ID: %s\n  New role: %s\n",
					details.ID, details.ChangesRequested.Role)
			case *openaiorgs.ServiceAccountDeleted:
				fmt.Printf("  Service account deleted with ID: %s\n", details.ID)
			case *openaiorgs.UserAdded:
				fmt.Printf("  User added with ID: %s\n  Role: %s\n",
					details.ID, details.Data.Role)
			case *openaiorgs.UserUpdated:
				fmt.Printf("  User updated with ID: %s\n  New role: %s\n",
					details.ID, details.ChangesRequested.Role)
			case *openaiorgs.UserDeleted:
				fmt.Printf("  User deleted with ID: %s\n", details.ID)
			default:
				fmt.Printf("  %#v\n", details)
			}
		}
		fmt.Println("\n---")
	}

	return nil
}
