package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/urfave/cli/v3"
)

func UsageCommand() *cli.Command {
	return &cli.Command{
		Name:  "usage",
		Usage: "Access usage data for the organization",
		Commands: []*cli.Command{
			{
				Name:   "completions",
				Usage:  "Get completions usage data",
				Action: getCompletionsUsage,
				Flags:  commonUsageFlags(),
			},
			{
				Name:   "embeddings",
				Usage:  "Get embeddings usage data",
				Action: getEmbeddingsUsage,
				Flags:  commonUsageFlags(),
			},
			{
				Name:   "moderations",
				Usage:  "Get moderations usage data",
				Action: getModerationsUsage,
				Flags:  commonUsageFlags(),
			},
			{
				Name:   "images",
				Usage:  "Get image creation usage data",
				Action: getImagesUsage,
				Flags:  commonUsageFlags(),
			},
			{
				Name:   "audio-speeches",
				Usage:  "Get audio speeches usage data",
				Action: getAudioSpeechesUsage,
				Flags:  commonUsageFlags(),
			},
			{
				Name:   "audio-transcriptions",
				Usage:  "Get audio transcriptions usage data",
				Action: getAudioTranscriptionsUsage,
				Flags:  commonUsageFlags(),
			},
			{
				Name:   "vector-stores",
				Usage:  "Get vector stores usage data",
				Action: getVectorStoresUsage,
				Flags:  commonUsageFlags(),
			},
			{
				Name:   "code-interpreter",
				Usage:  "Get code interpreter usage data",
				Action: getCodeInterpreterUsage,
				Flags:  commonUsageFlags(),
			},
			{
				Name:   "costs",
				Usage:  "Get costs usage data",
				Action: getCostsUsage,
				Flags:  commonUsageFlags(),
			},
		},
	}
}

func commonUsageFlags() []cli.Flag {
	return []cli.Flag{
		limitFlag,
		afterFlag,
		&cli.StringFlag{
			Name:  "start-date",
			Usage: "Start date for the query (YYYY-MM-DD format)",
		},
		&cli.StringFlag{
			Name:  "end-date",
			Usage: "End date for the query (YYYY-MM-DD format)",
		},
		&cli.StringFlag{
			Name:  "project-id",
			Usage: "Filter by project ID",
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
		&cli.BoolFlag{
			Name:  "paginate",
			Usage: "Automatically paginate through all results",
		},
	}
}

func buildUsageQueryParams(cmd *cli.Command) map[string]string {
	params := make(map[string]string)

	if cmd.IsSet("limit") {
		params["limit"] = fmt.Sprintf("%d", cmd.Int("limit"))
	}

	if cmd.IsSet("after") {
		params["after"] = cmd.String("after")
	}

	if cmd.IsSet("start-date") {
		params["start_date"] = cmd.String("start-date")
	}

	if cmd.IsSet("end-date") {
		params["end_date"] = cmd.String("end-date")
	}

	if cmd.IsSet("project-id") {
		params["project_id"] = cmd.String("project-id")
	}

	return params
}

func getCompletionsUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")

	var allUsage []openaiorgs.UsageRecord

	for {
		usage, err := client.GetCompletionsUsage(params)
		if err != nil {
			return wrapError("get completions usage", err)
		}

		if paginate {
			allUsage = append(allUsage, usage.Data...)
			if !usage.HasMore {
				break
			}
			params["after"] = usage.LastID
		} else {
			return outputUsageResponse(usage, outputFormat, verbose)
		}
	}

	if paginate {
		response := &openaiorgs.UsageResponse{
			Object: "list",
			Data:   allUsage,
		}
		return outputUsageResponse(response, outputFormat, verbose)
	}

	return nil
}

func getEmbeddingsUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")

	var allUsage []openaiorgs.UsageRecord

	for {
		usage, err := client.GetEmbeddingsUsage(params)
		if err != nil {
			return wrapError("get embeddings usage", err)
		}

		if paginate {
			allUsage = append(allUsage, usage.Data...)
			if !usage.HasMore {
				break
			}
			params["after"] = usage.LastID
		} else {
			return outputUsageResponse(usage, outputFormat, verbose)
		}
	}

	if paginate {
		response := &openaiorgs.UsageResponse{
			Object: "list",
			Data:   allUsage,
		}
		return outputUsageResponse(response, outputFormat, verbose)
	}

	return nil
}

func getModerationsUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")

	var allUsage []openaiorgs.UsageRecord

	for {
		usage, err := client.GetModerationsUsage(params)
		if err != nil {
			return wrapError("get moderations usage", err)
		}

		if paginate {
			allUsage = append(allUsage, usage.Data...)
			if !usage.HasMore {
				break
			}
			params["after"] = usage.LastID
		} else {
			return outputUsageResponse(usage, outputFormat, verbose)
		}
	}

	if paginate {
		response := &openaiorgs.UsageResponse{
			Object: "list",
			Data:   allUsage,
		}
		return outputUsageResponse(response, outputFormat, verbose)
	}

	return nil
}

func getImagesUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")

	var allUsage []openaiorgs.UsageRecord

	for {
		usage, err := client.GetImagesUsage(params)
		if err != nil {
			return wrapError("get images usage", err)
		}

		if paginate {
			allUsage = append(allUsage, usage.Data...)
			if !usage.HasMore {
				break
			}
			params["after"] = usage.LastID
		} else {
			return outputUsageResponse(usage, outputFormat, verbose)
		}
	}

	if paginate {
		response := &openaiorgs.UsageResponse{
			Object: "list",
			Data:   allUsage,
		}
		return outputUsageResponse(response, outputFormat, verbose)
	}

	return nil
}

func getAudioSpeechesUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")

	var allUsage []openaiorgs.UsageRecord

	for {
		usage, err := client.GetAudioSpeechesUsage(params)
		if err != nil {
			return wrapError("get audio speeches usage", err)
		}

		if paginate {
			allUsage = append(allUsage, usage.Data...)
			if !usage.HasMore {
				break
			}
			params["after"] = usage.LastID
		} else {
			return outputUsageResponse(usage, outputFormat, verbose)
		}
	}

	if paginate {
		response := &openaiorgs.UsageResponse{
			Object: "list",
			Data:   allUsage,
		}
		return outputUsageResponse(response, outputFormat, verbose)
	}

	return nil
}

func getAudioTranscriptionsUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")

	var allUsage []openaiorgs.UsageRecord

	for {
		usage, err := client.GetAudioTranscriptionsUsage(params)
		if err != nil {
			return wrapError("get audio transcriptions usage", err)
		}

		if paginate {
			allUsage = append(allUsage, usage.Data...)
			if !usage.HasMore {
				break
			}
			params["after"] = usage.LastID
		} else {
			return outputUsageResponse(usage, outputFormat, verbose)
		}
	}

	if paginate {
		response := &openaiorgs.UsageResponse{
			Object: "list",
			Data:   allUsage,
		}
		return outputUsageResponse(response, outputFormat, verbose)
	}

	return nil
}

func getVectorStoresUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")

	var allUsage []openaiorgs.UsageRecord

	for {
		usage, err := client.GetVectorStoresUsage(params)
		if err != nil {
			return wrapError("get vector stores usage", err)
		}

		if paginate {
			allUsage = append(allUsage, usage.Data...)
			if !usage.HasMore {
				break
			}
			params["after"] = usage.LastID
		} else {
			return outputUsageResponse(usage, outputFormat, verbose)
		}
	}

	if paginate {
		response := &openaiorgs.UsageResponse{
			Object: "list",
			Data:   allUsage,
		}
		return outputUsageResponse(response, outputFormat, verbose)
	}

	return nil
}

func getCodeInterpreterUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")

	var allUsage []openaiorgs.UsageRecord

	for {
		usage, err := client.GetCodeInterpreterUsage(params)
		if err != nil {
			return wrapError("get code interpreter usage", err)
		}

		if paginate {
			allUsage = append(allUsage, usage.Data...)
			if !usage.HasMore {
				break
			}
			params["after"] = usage.LastID
		} else {
			return outputUsageResponse(usage, outputFormat, verbose)
		}
	}

	if paginate {
		response := &openaiorgs.UsageResponse{
			Object: "list",
			Data:   allUsage,
		}
		return outputUsageResponse(response, outputFormat, verbose)
	}

	return nil
}

func getCostsUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")

	var allUsage []openaiorgs.UsageRecord

	for {
		usage, err := client.GetCostsUsage(params)
		if err != nil {
			return wrapError("get costs usage", err)
		}

		if paginate {
			allUsage = append(allUsage, usage.Data...)
			if !usage.HasMore {
				break
			}
			params["after"] = usage.LastID
		} else {
			return outputUsageResponse(usage, outputFormat, verbose)
		}
	}

	if paginate {
		response := &openaiorgs.UsageResponse{
			Object: "list",
			Data:   allUsage,
		}
		return outputUsageResponse(response, outputFormat, verbose)
	}

	return nil
}

func outputUsageResponse(response *openaiorgs.UsageResponse, outputFormat string, verbose bool) error {
	switch outputFormat {
	case "json":
		return outputUsageJSON(response, verbose)
	case "jsonl":
		return outputUsageJSONL(response, verbose)
	case "pretty":
		return outputUsagePretty(response, verbose)
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
}

func outputUsageJSON(response *openaiorgs.UsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

func outputUsageJSONL(response *openaiorgs.UsageResponse, verbose bool) error {
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

	// Output each usage entry on its own line
	for _, usage := range response.Data {
		if err := encoder.Encode(usage); err != nil {
			return err
		}
	}
	return nil
}

func outputUsagePretty(response *openaiorgs.UsageResponse, verbose bool) error {
	if verbose {
		fmt.Printf("=== Usage Summary ===\n")
		fmt.Printf("Total records: %d\nHas more: %v\nFirst ID: %s\nLast ID: %s\n\n",
			len(response.Data), response.HasMore, response.FirstID, response.LastID)
	}

	for _, usage := range response.Data {
		fmt.Printf("=== Usage Record ===\n")
		fmt.Printf("ID:        %s\n", usage.ID)
		fmt.Printf("Type:      %s\n", usage.Type)
		fmt.Printf("Timestamp: %s\n", usage.Timestamp.Format(time.RFC3339))
		fmt.Printf("Cost:      $%.4f\n", usage.Cost)
		fmt.Printf("Project:   %s\n", usage.ProjectID)
		if usage.UserID != "" {
			fmt.Printf("User:      %s\n", usage.UserID)
		}

		if usage.UsageDetails != nil {
			fmt.Printf("\nUsage Details:\n")

			switch usage.Type {
			case openaiorgs.UsageTypeCompletions:
				if details, ok := usage.UsageDetails.(map[string]interface{}); ok {
					if model, ok := details["model"].(string); ok {
						fmt.Printf("  Model:            %s\n", model)
					}
					if prompt, ok := details["prompt_tokens"].(float64); ok {
						fmt.Printf("  Prompt tokens:    %.0f\n", prompt)
					}
					if completion, ok := details["completion_tokens"].(float64); ok {
						fmt.Printf("  Completion tokens: %.0f\n", completion)
					}
					if total, ok := details["total_tokens"].(float64); ok {
						fmt.Printf("  Total tokens:     %.0f\n", total)
					}
				}

			case openaiorgs.UsageTypeEmbeddings:
				if details, ok := usage.UsageDetails.(map[string]interface{}); ok {
					if model, ok := details["model"].(string); ok {
						fmt.Printf("  Model:         %s\n", model)
					}
					if prompt, ok := details["prompt_tokens"].(float64); ok {
						fmt.Printf("  Prompt tokens: %.0f\n", prompt)
					}
				}

			case openaiorgs.UsageTypeModerations:
				if details, ok := usage.UsageDetails.(map[string]interface{}); ok {
					if model, ok := details["model"].(string); ok {
						fmt.Printf("  Model:         %s\n", model)
					}
					if prompt, ok := details["prompt_tokens"].(float64); ok {
						fmt.Printf("  Prompt tokens: %.0f\n", prompt)
					}
				}

			case openaiorgs.UsageTypeImages:
				if details, ok := usage.UsageDetails.(map[string]interface{}); ok {
					if model, ok := details["model"].(string); ok {
						fmt.Printf("  Model:  %s\n", model)
					}
					if images, ok := details["images"].(float64); ok {
						fmt.Printf("  Images: %.0f\n", images)
					}
					if size, ok := details["size"].(string); ok {
						fmt.Printf("  Size:   %s\n", size)
					}
				}

			case openaiorgs.UsageTypeAudioSpeeches:
				if details, ok := usage.UsageDetails.(map[string]interface{}); ok {
					if model, ok := details["model"].(string); ok {
						fmt.Printf("  Model:      %s\n", model)
					}
					if chars, ok := details["characters"].(float64); ok {
						fmt.Printf("  Characters: %.0f\n", chars)
					}
				}

			case openaiorgs.UsageTypeAudioTranscriptions:
				if details, ok := usage.UsageDetails.(map[string]interface{}); ok {
					if model, ok := details["model"].(string); ok {
						fmt.Printf("  Model:   %s\n", model)
					}
					if seconds, ok := details["seconds"].(float64); ok {
						fmt.Printf("  Seconds: %.0f\n", seconds)
					}
				}

			case openaiorgs.UsageTypeVectorStores:
				if details, ok := usage.UsageDetails.(map[string]interface{}); ok {
					if model, ok := details["model"].(string); ok {
						fmt.Printf("  Model:   %s\n", model)
					}
					if vectors, ok := details["vectors"].(float64); ok {
						fmt.Printf("  Vectors: %.0f\n", vectors)
					}
					if size, ok := details["size"].(float64); ok {
						fmt.Printf("  Size:    %.0f bytes\n", size)
					}
				}

			case openaiorgs.UsageTypeCodeInterpreter:
				if details, ok := usage.UsageDetails.(map[string]interface{}); ok {
					if model, ok := details["model"].(string); ok {
						fmt.Printf("  Model:           %s\n", model)
					}
					if duration, ok := details["session_duration"].(float64); ok {
						fmt.Printf("  Session duration: %.0f seconds\n", duration)
					}
				}

			default:
				// Display raw JSON for other types or costs
				jsonData, err := json.MarshalIndent(usage.UsageDetails, "  ", "  ")
				if err == nil {
					fmt.Println(string(jsonData))
				} else {
					fmt.Printf("  %v\n", usage.UsageDetails)
				}
			}
		}
		fmt.Println("\n---")
	}

	return nil
}
