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
			Name:     "start-date",
			Usage:    "Start date for the query (RFC3339 format)",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "end-date",
			Usage: "End date for the query (RFC3339 format)",
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

	// Required start_time parameter
	if cmd.IsSet("start-date") {
		t, err := time.Parse(time.RFC3339, cmd.String("start-date"))
		if err == nil {
			params["start_time"] = fmt.Sprintf("%d", t.Unix())
		}
	}

	// Optional end_time parameter
	if cmd.IsSet("end-date") {
		t, err := time.Parse(time.RFC3339, cmd.String("end-date"))
		if err == nil {
			params["end_time"] = fmt.Sprintf("%d", t.Unix())
		}
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

	var allBuckets []openaiorgs.CompletionsUsageBucket

	for {
		usage, err := client.GetCompletionsUsage(params)
		if err != nil {
			return wrapError("get completions usage", err)
		}

		if paginate {
			allBuckets = append(allBuckets, usage.Data...)
			if !usage.HasMore {
				break
			}
			// Update the pagination parameter for completions which uses page= instead of after=
			if usage.NextPage != "" {
				params["page"] = usage.NextPage
			} else {
				break
			}
		} else {
			return outputCompletionsUsageResponse(usage, outputFormat, verbose)
		}
	}

	if paginate {
		response := &openaiorgs.CompletionsUsageResponse{
			Object:  "page",
			Data:    allBuckets,
			HasMore: false,
		}
		return outputCompletionsUsageResponse(response, outputFormat, verbose)
	}

	return nil
}

// outputCompletionsUsageResponse handles output formatting for the completions usage response
func outputCompletionsUsageResponse(response *openaiorgs.CompletionsUsageResponse, outputFormat string, verbose bool) error {
	switch outputFormat {
	case "json":
		return outputCompletionsUsageJSON(response, verbose)
	case "jsonl":
		return outputCompletionsUsageJSONL(response, verbose)
	case "pretty":
		return outputCompletionsUsagePretty(response, verbose)
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
}

// outputCompletionsUsageJSON outputs the completions usage response as JSON
func outputCompletionsUsageJSON(response *openaiorgs.CompletionsUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

// outputCompletionsUsageJSONL outputs the completions usage response as JSONL
func outputCompletionsUsageJSONL(response *openaiorgs.CompletionsUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)

	// First line: metadata if verbose
	if verbose {
		meta := struct {
			Total    int    `json:"total"`
			HasMore  bool   `json:"has_more"`
			NextPage string `json:"next_page"`
		}{
			Total:    len(response.Data),
			HasMore:  response.HasMore,
			NextPage: response.NextPage,
		}
		if err := encoder.Encode(meta); err != nil {
			return err
		}
	}

	// Output each bucket and its results
	for _, bucket := range response.Data {
		// Output bucket info
		bucketInfo := struct {
			StartTime int64 `json:"start_time"`
			EndTime   int64 `json:"end_time"`
			Results   int   `json:"result_count"`
		}{
			StartTime: bucket.StartTime,
			EndTime:   bucket.EndTime,
			Results:   len(bucket.Results),
		}
		if err := encoder.Encode(bucketInfo); err != nil {
			return err
		}

		// Output each result in the bucket
		for _, result := range bucket.Results {
			if err := encoder.Encode(result); err != nil {
				return err
			}
		}
	}

	return nil
}

// outputCompletionsUsagePretty outputs the completions usage response in a human-readable format
func outputCompletionsUsagePretty(response *openaiorgs.CompletionsUsageResponse, verbose bool) error {
	if verbose {
		fmt.Printf("=== Completions Usage Summary ===\n")
		fmt.Printf("Total buckets: %d\nHas more: %v\nNext page: %s\n\n",
			len(response.Data), response.HasMore, response.NextPage)
	}

	totalResults := 0
	for _, bucket := range response.Data {
		totalResults += len(bucket.Results)
	}
	fmt.Printf("Total records: %d\n", totalResults)
	fmt.Printf("Has more: %v\n\n", response.HasMore)

	for _, bucket := range response.Data {
		fmt.Printf("=== Time Bucket ===\n")
		startTime := time.Unix(bucket.StartTime, 0).Format(time.RFC3339)
		endTime := time.Unix(bucket.EndTime, 0).Format(time.RFC3339)
		fmt.Printf("Start time: %s\n", startTime)
		fmt.Printf("End time:   %s\n", endTime)
		fmt.Printf("Results:    %d\n\n", len(bucket.Results))

		for _, result := range bucket.Results {
			fmt.Printf("--- Usage Record ---\n")
			fmt.Printf("Input tokens:        %d\n", result.InputTokens)
			fmt.Printf("Output tokens:       %d\n", result.OutputTokens)
			fmt.Printf("Cached input tokens: %d\n", result.InputCachedTokens)
			fmt.Printf("Model requests:      %d\n", result.NumModelRequests)

			if result.ProjectID != "" {
				fmt.Printf("Project ID:         %s\n", result.ProjectID)
			}
			if result.UserID != "" {
				fmt.Printf("User ID:            %s\n", result.UserID)
			}
			if result.APIKeyID != "" {
				fmt.Printf("API Key ID:         %s\n", result.APIKeyID)
			}
			if result.Model != "" {
				fmt.Printf("Model:              %s\n", result.Model)
			}
			fmt.Println("")
		}
	}

	return nil
}

func getEmbeddingsUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)
	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")

	var allBuckets []openaiorgs.EmbeddingsUsageBucket

	for {
		usage, err := client.GetEmbeddingsUsage(params)
		if err != nil {
			return wrapError("get embeddings usage", err)
		}

		if paginate {
			allBuckets = append(allBuckets, usage.Data...)
			if !usage.HasMore {
				break
			}
			// Update the pagination parameter which uses page= instead of after=
			if usage.NextPage != "" {
				params["page"] = usage.NextPage
			} else {
				break
			}
		} else {
			return outputEmbeddingsUsageResponse(usage, outputFormat, verbose)
		}
	}

	if paginate {
		response := &openaiorgs.EmbeddingsUsageResponse{
			Object:  "page",
			Data:    allBuckets,
			HasMore: false,
		}
		return outputEmbeddingsUsageResponse(response, outputFormat, verbose)
	}

	return nil
}

// outputEmbeddingsUsageResponse handles output formatting for the embeddings usage response
func outputEmbeddingsUsageResponse(response *openaiorgs.EmbeddingsUsageResponse, outputFormat string, verbose bool) error {
	switch outputFormat {
	case "json":
		return outputEmbeddingsUsageJSON(response, verbose)
	case "jsonl":
		return outputEmbeddingsUsageJSONL(response, verbose)
	case "pretty":
		return outputEmbeddingsUsagePretty(response, verbose)
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
}

// outputEmbeddingsUsageJSON outputs the embeddings usage response as JSON
func outputEmbeddingsUsageJSON(response *openaiorgs.EmbeddingsUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

// outputEmbeddingsUsageJSONL outputs the embeddings usage response as JSONL
func outputEmbeddingsUsageJSONL(response *openaiorgs.EmbeddingsUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)

	// First line: metadata if verbose
	if verbose {
		meta := struct {
			Total    int    `json:"total"`
			HasMore  bool   `json:"has_more"`
			NextPage string `json:"next_page"`
		}{
			Total:    len(response.Data),
			HasMore:  response.HasMore,
			NextPage: response.NextPage,
		}
		if err := encoder.Encode(meta); err != nil {
			return err
		}
	}

	// Output each bucket and its results
	for _, bucket := range response.Data {
		// Output bucket info
		bucketInfo := struct {
			StartTime int64 `json:"start_time"`
			EndTime   int64 `json:"end_time"`
			Results   int   `json:"result_count"`
		}{
			StartTime: bucket.StartTime,
			EndTime:   bucket.EndTime,
			Results:   len(bucket.Results),
		}
		if err := encoder.Encode(bucketInfo); err != nil {
			return err
		}

		// Output each result in the bucket
		for _, result := range bucket.Results {
			if err := encoder.Encode(result); err != nil {
				return err
			}
		}
	}

	return nil
}

// outputEmbeddingsUsagePretty outputs the embeddings usage response in a human-readable format
func outputEmbeddingsUsagePretty(response *openaiorgs.EmbeddingsUsageResponse, verbose bool) error {
	if verbose {
		fmt.Printf("=== Embeddings Usage Summary ===\n")
		fmt.Printf("Total buckets: %d\nHas more: %v\nNext page: %s\n\n",
			len(response.Data), response.HasMore, response.NextPage)
	}

	totalResults := 0
	for _, bucket := range response.Data {
		totalResults += len(bucket.Results)
	}
	fmt.Printf("Total records: %d\n", totalResults)
	fmt.Printf("Has more: %v\n\n", response.HasMore)

	for _, bucket := range response.Data {
		fmt.Printf("=== Time Bucket ===\n")
		startTime := time.Unix(bucket.StartTime, 0).Format(time.RFC3339)
		endTime := time.Unix(bucket.EndTime, 0).Format(time.RFC3339)
		fmt.Printf("Start time: %s\n", startTime)
		fmt.Printf("End time:   %s\n", endTime)
		fmt.Printf("Results:    %d\n\n", len(bucket.Results))

		for _, result := range bucket.Results {
			fmt.Printf("--- Usage Record ---\n")
			fmt.Printf("Input tokens:   %d\n", result.InputTokens)
			fmt.Printf("Model requests: %d\n", result.NumModelRequests)

			if result.ProjectID != "" {
				fmt.Printf("Project ID:    %s\n", result.ProjectID)
			}
			if result.UserID != "" {
				fmt.Printf("User ID:       %s\n", result.UserID)
			}
			if result.APIKeyID != "" {
				fmt.Printf("API Key ID:    %s\n", result.APIKeyID)
			}
			if result.Model != "" {
				fmt.Printf("Model:         %s\n", result.Model)
			}
			fmt.Println("")
		}
	}

	return nil
}

func getModerationsUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)
	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")

	var allBuckets []openaiorgs.ModerationsUsageBucket

	for {
		usage, err := client.GetModerationsUsage(params)
		if err != nil {
			return wrapError("get moderations usage", err)
		}

		if paginate {
			allBuckets = append(allBuckets, usage.Data...)
			if !usage.HasMore {
				break
			}
			// Update the pagination parameter which uses page= instead of after=
			if usage.NextPage != "" {
				params["page"] = usage.NextPage
			} else {
				break
			}
		} else {
			return outputModerationsUsageResponse(usage, outputFormat, verbose)
		}
	}

	if paginate {
		response := &openaiorgs.ModerationsUsageResponse{
			Object:  "page",
			Data:    allBuckets,
			HasMore: false,
		}
		return outputModerationsUsageResponse(response, outputFormat, verbose)
	}

	return nil
}

// outputModerationsUsageResponse handles output formatting for the moderations usage response
func outputModerationsUsageResponse(response *openaiorgs.ModerationsUsageResponse, outputFormat string, verbose bool) error {
	switch outputFormat {
	case "json":
		return outputModerationsUsageJSON(response, verbose)
	case "jsonl":
		return outputModerationsUsageJSONL(response, verbose)
	case "pretty":
		return outputModerationsUsagePretty(response, verbose)
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
}

// outputModerationsUsageJSON outputs the moderations usage response as JSON
func outputModerationsUsageJSON(response *openaiorgs.ModerationsUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

// outputModerationsUsageJSONL outputs the moderations usage response as JSONL
func outputModerationsUsageJSONL(response *openaiorgs.ModerationsUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)

	// First line: metadata if verbose
	if verbose {
		meta := struct {
			Total    int    `json:"total"`
			HasMore  bool   `json:"has_more"`
			NextPage string `json:"next_page"`
		}{
			Total:    len(response.Data),
			HasMore:  response.HasMore,
			NextPage: response.NextPage,
		}
		if err := encoder.Encode(meta); err != nil {
			return err
		}
	}

	// Output each bucket and its results
	for _, bucket := range response.Data {
		// Output bucket info
		bucketInfo := struct {
			StartTime int64 `json:"start_time"`
			EndTime   int64 `json:"end_time"`
			Results   int   `json:"result_count"`
		}{
			StartTime: bucket.StartTime,
			EndTime:   bucket.EndTime,
			Results:   len(bucket.Results),
		}
		if err := encoder.Encode(bucketInfo); err != nil {
			return err
		}

		// Output each result in the bucket
		for _, result := range bucket.Results {
			if err := encoder.Encode(result); err != nil {
				return err
			}
		}
	}

	return nil
}

// outputModerationsUsagePretty outputs the moderations usage response in a human-readable format
func outputModerationsUsagePretty(response *openaiorgs.ModerationsUsageResponse, verbose bool) error {
	if verbose {
		fmt.Printf("=== Moderations Usage Summary ===\n")
		fmt.Printf("Total buckets: %d\nHas more: %v\nNext page: %s\n\n",
			len(response.Data), response.HasMore, response.NextPage)
	}

	totalResults := 0
	for _, bucket := range response.Data {
		totalResults += len(bucket.Results)
	}
	fmt.Printf("Total records: %d\n", totalResults)
	fmt.Printf("Has more: %v\n\n", response.HasMore)

	for _, bucket := range response.Data {
		fmt.Printf("=== Time Bucket ===\n")
		startTime := time.Unix(bucket.StartTime, 0).Format(time.RFC3339)
		endTime := time.Unix(bucket.EndTime, 0).Format(time.RFC3339)
		fmt.Printf("Start time: %s\n", startTime)
		fmt.Printf("End time:   %s\n", endTime)
		fmt.Printf("Results:    %d\n\n", len(bucket.Results))

		for _, result := range bucket.Results {
			fmt.Printf("--- Usage Record ---\n")
			fmt.Printf("Input tokens:   %d\n", result.InputTokens)
			fmt.Printf("Model requests: %d\n", result.NumModelRequests)

			if result.ProjectID != "" {
				fmt.Printf("Project ID:    %s\n", result.ProjectID)
			}
			if result.UserID != "" {
				fmt.Printf("User ID:       %s\n", result.UserID)
			}
			if result.APIKeyID != "" {
				fmt.Printf("API Key ID:    %s\n", result.APIKeyID)
			}
			if result.Model != "" {
				fmt.Printf("Model:         %s\n", result.Model)
			}
			fmt.Println("")
		}
	}

	return nil
}

func getImagesUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)
	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")

	var allBuckets []openaiorgs.ImagesUsageBucket

	for {
		usage, err := client.GetImagesUsage(params)
		if err != nil {
			return wrapError("get images usage", err)
		}

		if paginate {
			allBuckets = append(allBuckets, usage.Data...)
			if !usage.HasMore {
				break
			}
			// Update the pagination parameter which uses page= instead of after=
			if usage.NextPage != "" {
				params["page"] = usage.NextPage
			} else {
				break
			}
		} else {
			return outputImagesUsageResponse(usage, outputFormat, verbose)
		}
	}

	if paginate {
		response := &openaiorgs.ImagesUsageResponse{
			Object:  "page",
			Data:    allBuckets,
			HasMore: false,
		}
		return outputImagesUsageResponse(response, outputFormat, verbose)
	}

	return nil
}

// outputImagesUsageResponse handles output formatting for the images usage response
func outputImagesUsageResponse(response *openaiorgs.ImagesUsageResponse, outputFormat string, verbose bool) error {
	switch outputFormat {
	case "json":
		return outputImagesUsageJSON(response, verbose)
	case "jsonl":
		return outputImagesUsageJSONL(response, verbose)
	case "pretty":
		return outputImagesUsagePretty(response, verbose)
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
}

// outputImagesUsageJSON outputs the images usage response as JSON
func outputImagesUsageJSON(response *openaiorgs.ImagesUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

// outputImagesUsageJSONL outputs the images usage response as JSONL
func outputImagesUsageJSONL(response *openaiorgs.ImagesUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)

	// First line: metadata if verbose
	if verbose {
		meta := struct {
			Total    int    `json:"total"`
			HasMore  bool   `json:"has_more"`
			NextPage string `json:"next_page"`
		}{
			Total:    len(response.Data),
			HasMore:  response.HasMore,
			NextPage: response.NextPage,
		}
		if err := encoder.Encode(meta); err != nil {
			return err
		}
	}

	// Output each bucket and its results
	for _, bucket := range response.Data {
		// Output bucket info
		bucketInfo := struct {
			StartTime int64 `json:"start_time"`
			EndTime   int64 `json:"end_time"`
			Results   int   `json:"result_count"`
		}{
			StartTime: bucket.StartTime,
			EndTime:   bucket.EndTime,
			Results:   len(bucket.Results),
		}
		if err := encoder.Encode(bucketInfo); err != nil {
			return err
		}

		// Output each result in the bucket
		for _, result := range bucket.Results {
			if err := encoder.Encode(result); err != nil {
				return err
			}
		}
	}

	return nil
}

// outputImagesUsagePretty outputs the images usage response in a human-readable format
func outputImagesUsagePretty(response *openaiorgs.ImagesUsageResponse, verbose bool) error {
	if verbose {
		fmt.Printf("=== Images Usage Summary ===\n")
		fmt.Printf("Total buckets: %d\nHas more: %v\nNext page: %s\n\n",
			len(response.Data), response.HasMore, response.NextPage)
	}

	totalResults := 0
	for _, bucket := range response.Data {
		totalResults += len(bucket.Results)
	}
	fmt.Printf("Total records: %d\n", totalResults)
	fmt.Printf("Has more: %v\n\n", response.HasMore)

	for _, bucket := range response.Data {
		fmt.Printf("=== Time Bucket ===\n")
		startTime := time.Unix(bucket.StartTime, 0).Format(time.RFC3339)
		endTime := time.Unix(bucket.EndTime, 0).Format(time.RFC3339)
		fmt.Printf("Start time: %s\n", startTime)
		fmt.Printf("End time:   %s\n", endTime)
		fmt.Printf("Results:    %d\n\n", len(bucket.Results))

		for _, result := range bucket.Results {
			fmt.Printf("--- Usage Record ---\n")
			fmt.Printf("Images:        %d\n", result.Images)
			fmt.Printf("Model requests: %d\n", result.NumModelRequests)

			if result.Size != "" {
				fmt.Printf("Size:          %s\n", result.Size)
			}
			if result.Source != "" {
				fmt.Printf("Source:        %s\n", result.Source)
			}
			if result.ProjectID != "" {
				fmt.Printf("Project ID:    %s\n", result.ProjectID)
			}
			if result.UserID != "" {
				fmt.Printf("User ID:       %s\n", result.UserID)
			}
			if result.APIKeyID != "" {
				fmt.Printf("API Key ID:    %s\n", result.APIKeyID)
			}
			if result.Model != "" {
				fmt.Printf("Model:         %s\n", result.Model)
			}
			fmt.Println("")
		}
	}

	return nil
}

func getAudioSpeechesUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)
	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")
	var allBuckets []openaiorgs.AudioSpeechesUsageBucket
	for {
		usage, err := client.GetAudioSpeechesUsage(params)
		if err != nil {
			return wrapError("get audio speeches usage", err)
		}
		if paginate {
			allBuckets = append(allBuckets, usage.Data...)
			if !usage.HasMore {
				break
			}
			// Update the pagination parameter which uses page= instead of after=
			if usage.NextPage != "" {
				params["page"] = usage.NextPage
			} else {
				break
			}
		} else {
			return outputAudioSpeechesUsageResponse(usage, outputFormat, verbose)
		}
	}
	if paginate {
		response := &openaiorgs.AudioSpeechesUsageResponse{
			Object:  "page",
			Data:    allBuckets,
			HasMore: false,
		}
		return outputAudioSpeechesUsageResponse(response, outputFormat, verbose)
	}
	return nil
}

// outputAudioSpeechesUsageResponse handles output formatting for the audio speeches usage response
func outputAudioSpeechesUsageResponse(response *openaiorgs.AudioSpeechesUsageResponse, outputFormat string, verbose bool) error {
	switch outputFormat {
	case "json":
		return outputAudioSpeechesUsageJSON(response, verbose)
	case "jsonl":
		return outputAudioSpeechesUsageJSONL(response, verbose)
	case "pretty":
		return outputAudioSpeechesUsagePretty(response, verbose)
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
}

// outputAudioSpeechesUsageJSON outputs the audio speeches usage response as JSON
func outputAudioSpeechesUsageJSON(response *openaiorgs.AudioSpeechesUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

// outputAudioSpeechesUsageJSONL outputs the audio speeches usage response as JSONL
func outputAudioSpeechesUsageJSONL(response *openaiorgs.AudioSpeechesUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	// First line: metadata if verbose
	if verbose {
		meta := struct {
			Total    int    `json:"total"`
			HasMore  bool   `json:"has_more"`
			NextPage string `json:"next_page"`
		}{
			Total:    len(response.Data),
			HasMore:  response.HasMore,
			NextPage: response.NextPage,
		}
		if err := encoder.Encode(meta); err != nil {
			return err
		}
	}
	// Output each bucket and its results
	for _, bucket := range response.Data {
		// Output bucket info
		bucketInfo := struct {
			StartTime int64 `json:"start_time"`
			EndTime   int64 `json:"end_time"`
			Results   int   `json:"result_count"`
		}{
			StartTime: bucket.StartTime,
			EndTime:   bucket.EndTime,
			Results:   len(bucket.Results),
		}
		if err := encoder.Encode(bucketInfo); err != nil {
			return err
		}
		// Output each result in the bucket
		for _, result := range bucket.Results {
			if err := encoder.Encode(result); err != nil {
				return err
			}
		}
	}
	return nil
}

// outputAudioSpeechesUsagePretty outputs the audio speeches usage response in a human-readable format
func outputAudioSpeechesUsagePretty(response *openaiorgs.AudioSpeechesUsageResponse, verbose bool) error {
	if verbose {
		fmt.Printf("=== Audio Speeches Usage Summary ===\n")
		fmt.Printf("Total buckets: %d\nHas more: %v\nNext page: %s\n\n",
			len(response.Data), response.HasMore, response.NextPage)
	}
	totalResults := 0
	for _, bucket := range response.Data {
		totalResults += len(bucket.Results)
	}
	fmt.Printf("Total records: %d\n", totalResults)
	fmt.Printf("Has more: %v\n\n", response.HasMore)
	for _, bucket := range response.Data {
		fmt.Printf("=== Time Bucket ===\n")
		startTime := time.Unix(bucket.StartTime, 0).Format(time.RFC3339)
		endTime := time.Unix(bucket.EndTime, 0).Format(time.RFC3339)
		fmt.Printf("Start time: %s\n", startTime)
		fmt.Printf("End time:   %s\n", endTime)
		fmt.Printf("Results:    %d\n\n", len(bucket.Results))
		for _, result := range bucket.Results {
			fmt.Printf("--- Usage Record ---\n")
			fmt.Printf("Characters:     %d\n", result.Characters)
			fmt.Printf("Model requests: %d\n", result.NumModelRequests)
			if result.ProjectID != "" {
				fmt.Printf("Project ID:     %s\n", result.ProjectID)
			}
			if result.UserID != "" {
				fmt.Printf("User ID:        %s\n", result.UserID)
			}
			if result.APIKeyID != "" {
				fmt.Printf("API Key ID:     %s\n", result.APIKeyID)
			}
			if result.Model != "" {
				fmt.Printf("Model:          %s\n", result.Model)
			}
			fmt.Println("")
		}
	}
	return nil
}

func getAudioTranscriptionsUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)
	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")
	var allBuckets []openaiorgs.AudioTranscriptionsUsageBucket
	for {
		usage, err := client.GetAudioTranscriptionsUsage(params)
		if err != nil {
			return wrapError("get audio transcriptions usage", err)
		}
		if paginate {
			allBuckets = append(allBuckets, usage.Data...)
			if !usage.HasMore {
				break
			}
			// Update the pagination parameter which uses page= instead of after=
			if usage.NextPage != "" {
				params["page"] = usage.NextPage
			} else {
				break
			}
		} else {
			return outputAudioTranscriptionsUsageResponse(usage, outputFormat, verbose)
		}
	}
	if paginate {
		response := &openaiorgs.AudioTranscriptionsUsageResponse{
			Object:  "page",
			Data:    allBuckets,
			HasMore: false,
		}
		return outputAudioTranscriptionsUsageResponse(response, outputFormat, verbose)
	}
	return nil
}

// outputAudioTranscriptionsUsageResponse handles output formatting for the audio transcriptions usage response
func outputAudioTranscriptionsUsageResponse(response *openaiorgs.AudioTranscriptionsUsageResponse, outputFormat string, verbose bool) error {
	switch outputFormat {
	case "json":
		return outputAudioTranscriptionsUsageJSON(response, verbose)
	case "jsonl":
		return outputAudioTranscriptionsUsageJSONL(response, verbose)
	case "pretty":
		return outputAudioTranscriptionsUsagePretty(response, verbose)
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
}

// outputAudioTranscriptionsUsageJSON outputs the audio transcriptions usage response as JSON
func outputAudioTranscriptionsUsageJSON(response *openaiorgs.AudioTranscriptionsUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

// outputAudioTranscriptionsUsageJSONL outputs the audio transcriptions usage response as JSONL
func outputAudioTranscriptionsUsageJSONL(response *openaiorgs.AudioTranscriptionsUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	// First line: metadata if verbose
	if verbose {
		meta := struct {
			Total    int    `json:"total"`
			HasMore  bool   `json:"has_more"`
			NextPage string `json:"next_page"`
		}{
			Total:    len(response.Data),
			HasMore:  response.HasMore,
			NextPage: response.NextPage,
		}
		if err := encoder.Encode(meta); err != nil {
			return err
		}
	}
	// Output each bucket and its results
	for _, bucket := range response.Data {
		// Output bucket info
		bucketInfo := struct {
			StartTime int64 `json:"start_time"`
			EndTime   int64 `json:"end_time"`
			Results   int   `json:"result_count"`
		}{
			StartTime: bucket.StartTime,
			EndTime:   bucket.EndTime,
			Results:   len(bucket.Results),
		}
		if err := encoder.Encode(bucketInfo); err != nil {
			return err
		}
		// Output each result in the bucket
		for _, result := range bucket.Results {
			if err := encoder.Encode(result); err != nil {
				return err
			}
		}
	}
	return nil
}

// outputAudioTranscriptionsUsagePretty outputs the audio transcriptions usage response in a human-readable format
func outputAudioTranscriptionsUsagePretty(response *openaiorgs.AudioTranscriptionsUsageResponse, verbose bool) error {
	if verbose {
		fmt.Printf("=== Audio Transcriptions Usage Summary ===\n")
		fmt.Printf("Total buckets: %d\nHas more: %v\nNext page: %s\n\n",
			len(response.Data), response.HasMore, response.NextPage)
	}
	totalResults := 0
	for _, bucket := range response.Data {
		totalResults += len(bucket.Results)
	}
	fmt.Printf("Total records: %d\n", totalResults)
	fmt.Printf("Has more: %v\n\n", response.HasMore)
	for _, bucket := range response.Data {
		fmt.Printf("=== Time Bucket ===\n")
		startTime := time.Unix(bucket.StartTime, 0).Format(time.RFC3339)
		endTime := time.Unix(bucket.EndTime, 0).Format(time.RFC3339)
		fmt.Printf("Start time: %s\n", startTime)
		fmt.Printf("End time:   %s\n", endTime)
		fmt.Printf("Results:    %d\n\n", len(bucket.Results))
		for _, result := range bucket.Results {
			fmt.Printf("--- Usage Record ---\n")
			fmt.Printf("Seconds:        %d\n", result.Seconds)
			fmt.Printf("Model requests: %d\n", result.NumModelRequests)
			if result.ProjectID != "" {
				fmt.Printf("Project ID:     %s\n", result.ProjectID)
			}
			if result.UserID != "" {
				fmt.Printf("User ID:        %s\n", result.UserID)
			}
			if result.APIKeyID != "" {
				fmt.Printf("API Key ID:     %s\n", result.APIKeyID)
			}
			if result.Model != "" {
				fmt.Printf("Model:          %s\n", result.Model)
			}
			fmt.Println("")
		}
	}
	return nil
}

func getVectorStoresUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)
	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")
	var allBuckets []openaiorgs.VectorStoresUsageBucket
	for {
		usage, err := client.GetVectorStoresUsage(params)
		if err != nil {
			return wrapError("get vector stores usage", err)
		}
		if paginate {
			allBuckets = append(allBuckets, usage.Data...)
			if !usage.HasMore {
				break
			}
			// Update the pagination parameter which uses page= instead of after=
			if usage.NextPage != "" {
				params["page"] = usage.NextPage
			} else {
				break
			}
		} else {
			return outputVectorStoresUsageResponse(usage, outputFormat, verbose)
		}
	}
	if paginate {
		response := &openaiorgs.VectorStoresUsageResponse{
			Object:  "page",
			Data:    allBuckets,
			HasMore: false,
		}
		return outputVectorStoresUsageResponse(response, outputFormat, verbose)
	}
	return nil
}

// outputVectorStoresUsageResponse handles output formatting for the vector stores usage response
func outputVectorStoresUsageResponse(response *openaiorgs.VectorStoresUsageResponse, outputFormat string, verbose bool) error {
	switch outputFormat {
	case "json":
		return outputVectorStoresUsageJSON(response, verbose)
	case "jsonl":
		return outputVectorStoresUsageJSONL(response, verbose)
	case "pretty":
		return outputVectorStoresUsagePretty(response, verbose)
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
}

// outputVectorStoresUsageJSON outputs the vector stores usage response as JSON
func outputVectorStoresUsageJSON(response *openaiorgs.VectorStoresUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

// outputVectorStoresUsageJSONL outputs the vector stores usage response as JSONL
func outputVectorStoresUsageJSONL(response *openaiorgs.VectorStoresUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	// First line: metadata if verbose
	if verbose {
		meta := struct {
			Total    int    `json:"total"`
			HasMore  bool   `json:"has_more"`
			NextPage string `json:"next_page"`
		}{
			Total:    len(response.Data),
			HasMore:  response.HasMore,
			NextPage: response.NextPage,
		}
		if err := encoder.Encode(meta); err != nil {
			return err
		}
	}
	// Output each bucket and its results
	for _, bucket := range response.Data {
		// Output bucket info
		bucketInfo := struct {
			StartTime int64 `json:"start_time"`
			EndTime   int64 `json:"end_time"`
			Results   int   `json:"result_count"`
		}{
			StartTime: bucket.StartTime,
			EndTime:   bucket.EndTime,
			Results:   len(bucket.Results),
		}
		if err := encoder.Encode(bucketInfo); err != nil {
			return err
		}
		// Output each result in the bucket
		for _, result := range bucket.Results {
			if err := encoder.Encode(result); err != nil {
				return err
			}
		}
	}
	return nil
}

// outputVectorStoresUsagePretty outputs the vector stores usage response in a human-readable format
func outputVectorStoresUsagePretty(response *openaiorgs.VectorStoresUsageResponse, verbose bool) error {
	if verbose {
		fmt.Printf("=== Vector Stores Usage Summary ===\n")
		fmt.Printf("Total buckets: %d\nHas more: %v\nNext page: %s\n\n",
			len(response.Data), response.HasMore, response.NextPage)
	}
	totalResults := 0
	for _, bucket := range response.Data {
		totalResults += len(bucket.Results)
	}
	fmt.Printf("Total records: %d\n", totalResults)
	fmt.Printf("Has more: %v\n\n", response.HasMore)
	for _, bucket := range response.Data {
		fmt.Printf("=== Time Bucket ===\n")
		startTime := time.Unix(bucket.StartTime, 0).Format(time.RFC3339)
		endTime := time.Unix(bucket.EndTime, 0).Format(time.RFC3339)
		fmt.Printf("Start time: %s\n", startTime)
		fmt.Printf("End time:   %s\n", endTime)
		fmt.Printf("Results:    %d\n\n", len(bucket.Results))
		for _, result := range bucket.Results {
			fmt.Printf("--- Usage Record ---\n")
			fmt.Printf("Usage bytes:  %d\n", result.UsageBytes)
			if result.ProjectID != "" {
				fmt.Printf("Project ID:   %s\n", result.ProjectID)
			}
			fmt.Println("")
		}
	}
	return nil
}

func getCodeInterpreterUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)
	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")
	var allBuckets []openaiorgs.CodeInterpreterUsageBucket
	for {
		usage, err := client.GetCodeInterpreterUsage(params)
		if err != nil {
			return wrapError("get code interpreter usage", err)
		}
		if paginate {
			allBuckets = append(allBuckets, usage.Data...)
			if !usage.HasMore {
				break
			}
			// Update the pagination parameter which uses page= instead of after=
			if usage.NextPage != "" {
				params["page"] = usage.NextPage
			} else {
				break
			}
		} else {
			return outputCodeInterpreterUsageResponse(usage, outputFormat, verbose)
		}
	}
	if paginate {
		response := &openaiorgs.CodeInterpreterUsageResponse{
			Object:  "page",
			Data:    allBuckets,
			HasMore: false,
		}
		return outputCodeInterpreterUsageResponse(response, outputFormat, verbose)
	}
	return nil
}

// outputCodeInterpreterUsageResponse handles output formatting for the code interpreter usage response
func outputCodeInterpreterUsageResponse(response *openaiorgs.CodeInterpreterUsageResponse, outputFormat string, verbose bool) error {
	switch outputFormat {
	case "json":
		return outputCodeInterpreterUsageJSON(response, verbose)
	case "jsonl":
		return outputCodeInterpreterUsageJSONL(response, verbose)
	case "pretty":
		return outputCodeInterpreterUsagePretty(response, verbose)
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
}

// outputCodeInterpreterUsageJSON outputs the code interpreter usage response as JSON
func outputCodeInterpreterUsageJSON(response *openaiorgs.CodeInterpreterUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

// outputCodeInterpreterUsageJSONL outputs the code interpreter usage response as JSONL
func outputCodeInterpreterUsageJSONL(response *openaiorgs.CodeInterpreterUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	// First line: metadata if verbose
	if verbose {
		meta := struct {
			Total    int    `json:"total"`
			HasMore  bool   `json:"has_more"`
			NextPage string `json:"next_page"`
		}{
			Total:    len(response.Data),
			HasMore:  response.HasMore,
			NextPage: response.NextPage,
		}
		if err := encoder.Encode(meta); err != nil {
			return err
		}
	}
	// Output each bucket and its results
	for _, bucket := range response.Data {
		// Output bucket info
		bucketInfo := struct {
			StartTime int64 `json:"start_time"`
			EndTime   int64 `json:"end_time"`
			Results   int   `json:"result_count"`
		}{
			StartTime: bucket.StartTime,
			EndTime:   bucket.EndTime,
			Results:   len(bucket.Results),
		}
		if err := encoder.Encode(bucketInfo); err != nil {
			return err
		}
		// Output each result in the bucket
		for _, result := range bucket.Results {
			if err := encoder.Encode(result); err != nil {
				return err
			}
		}
	}
	return nil
}

// outputCodeInterpreterUsagePretty outputs the code interpreter usage response in a human-readable format
func outputCodeInterpreterUsagePretty(response *openaiorgs.CodeInterpreterUsageResponse, verbose bool) error {
	if verbose {
		fmt.Printf("=== Code Interpreter Usage Summary ===\n")
		fmt.Printf("Total buckets: %d\nHas more: %v\nNext page: %s\n\n",
			len(response.Data), response.HasMore, response.NextPage)
	}
	totalResults := 0
	for _, bucket := range response.Data {
		totalResults += len(bucket.Results)
	}
	fmt.Printf("Total records: %d\n", totalResults)
	fmt.Printf("Has more: %v\n\n", response.HasMore)
	for _, bucket := range response.Data {
		fmt.Printf("=== Time Bucket ===\n")
		startTime := time.Unix(bucket.StartTime, 0).Format(time.RFC3339)
		endTime := time.Unix(bucket.EndTime, 0).Format(time.RFC3339)
		fmt.Printf("Start time: %s\n", startTime)
		fmt.Printf("End time:   %s\n", endTime)
		fmt.Printf("Results:    %d\n\n", len(bucket.Results))
		for _, result := range bucket.Results {
			fmt.Printf("--- Usage Record ---\n")
			fmt.Printf("Number of sessions: %d\n", result.NumSessions)
			if result.ProjectID != "" {
				fmt.Printf("Project ID:         %s\n", result.ProjectID)
			}
			fmt.Println("")
		}
	}
	return nil
}

func getCostsUsage(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)
	params := buildUsageQueryParams(cmd)
	outputFormat := cmd.String("output")
	verbose := cmd.Bool("verbose")
	paginate := cmd.Bool("paginate")
	var allBuckets []openaiorgs.CostsUsageBucket
	for {
		usage, err := client.GetCostsUsage(params)
		if err != nil {
			return wrapError("get costs usage", err)
		}
		if paginate {
			allBuckets = append(allBuckets, usage.Data...)
			if !usage.HasMore {
				break
			}
			// Update the pagination parameter which uses page= instead of after=
			if usage.NextPage != "" {
				params["page"] = usage.NextPage
			} else {
				break
			}
		} else {
			return outputCostsUsageResponse(usage, outputFormat, verbose)
		}
	}
	if paginate {
		response := &openaiorgs.CostsUsageResponse{
			Object:  "page",
			Data:    allBuckets,
			HasMore: false,
		}
		return outputCostsUsageResponse(response, outputFormat, verbose)
	}
	return nil
}

// outputCostsUsageResponse handles output formatting for the costs usage response
func outputCostsUsageResponse(response *openaiorgs.CostsUsageResponse, outputFormat string, verbose bool) error {
	switch outputFormat {
	case "json":
		return outputCostsUsageJSON(response, verbose)
	case "jsonl":
		return outputCostsUsageJSONL(response, verbose)
	case "pretty":
		return outputCostsUsagePretty(response, verbose)
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
}

// outputCostsUsageJSON outputs the costs usage response as JSON
func outputCostsUsageJSON(response *openaiorgs.CostsUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(response)
}

// outputCostsUsageJSONL outputs the costs usage response as JSONL
func outputCostsUsageJSONL(response *openaiorgs.CostsUsageResponse, verbose bool) error {
	encoder := json.NewEncoder(os.Stdout)
	// First line: metadata if verbose
	if verbose {
		meta := struct {
			Total    int    `json:"total"`
			HasMore  bool   `json:"has_more"`
			NextPage string `json:"next_page"`
		}{
			Total:    len(response.Data),
			HasMore:  response.HasMore,
			NextPage: response.NextPage,
		}
		if err := encoder.Encode(meta); err != nil {
			return err
		}
	}
	// Output each bucket and its results
	for _, bucket := range response.Data {
		// Output bucket info
		bucketInfo := struct {
			StartTime int64 `json:"start_time"`
			EndTime   int64 `json:"end_time"`
			Results   int   `json:"result_count"`
		}{
			StartTime: bucket.StartTime,
			EndTime:   bucket.EndTime,
			Results:   len(bucket.Results),
		}
		if err := encoder.Encode(bucketInfo); err != nil {
			return err
		}
		// Output each result in the bucket
		for _, result := range bucket.Results {
			if err := encoder.Encode(result); err != nil {
				return err
			}
		}
	}
	return nil
}

// outputCostsUsagePretty outputs the costs usage response in a human-readable format
func outputCostsUsagePretty(response *openaiorgs.CostsUsageResponse, verbose bool) error {
	if verbose {
		fmt.Printf("=== Costs Usage Summary ===\n")
		fmt.Printf("Total buckets: %d\nHas more: %v\nNext page: %s\n\n",
			len(response.Data), response.HasMore, response.NextPage)
	}
	totalResults := 0
	for _, bucket := range response.Data {
		totalResults += len(bucket.Results)
	}
	fmt.Printf("Total records: %d\n", totalResults)
	fmt.Printf("Has more: %v\n\n", response.HasMore)
	for _, bucket := range response.Data {
		fmt.Printf("=== Time Bucket ===\n")
		startTime := time.Unix(bucket.StartTime, 0).Format(time.RFC3339)
		endTime := time.Unix(bucket.EndTime, 0).Format(time.RFC3339)
		fmt.Printf("Start time: %s\n", startTime)
		fmt.Printf("End time:   %s\n", endTime)
		fmt.Printf("Results:    %d\n\n", len(bucket.Results))
		for _, result := range bucket.Results {
			fmt.Printf("--- Usage Record ---\n")
			fmt.Printf("Amount: %.2f %s\n", result.Amount.Value, result.Amount.Currency)
			if result.ProjectID != "" {
				fmt.Printf("Project ID: %s\n", result.ProjectID)
			}
			if result.LineItem != nil {
				fmt.Printf("Line item: %v\n", result.LineItem)
			}
			fmt.Println("")
		}
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
