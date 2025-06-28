package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

func CertificatesCommand() *cli.Command {
	return &cli.Command{
		Name:  "certificates",
		Usage: "Manage organization and project certificates",
		Commands: []*cli.Command{
			certificatesOrgCommand(),
			certificatesProjectCommand(),
		},
	}
}

func certificatesOrgCommand() *cli.Command {
	return &cli.Command{
		Name:  "org",
		Usage: "Manage organization-level certificates",
		Commands: []*cli.Command{
			listOrgCertificatesCommand(),
			uploadCertificateCommand(),
			getCertificateCommand(),
			modifyCertificateCommand(),
			deleteCertificateCommand(),
			activateOrgCertificatesCommand(),
			deactivateOrgCertificatesCommand(),
		},
	}
}

func certificatesProjectCommand() *cli.Command {
	return &cli.Command{
		Name:  "project",
		Usage: "Manage project-level certificates",
		Commands: []*cli.Command{
			listProjectCertificatesCommand(),
			activateProjectCertificatesCommand(),
			deactivateProjectCertificatesCommand(),
		},
	}
}

// Organization certificate commands

func listOrgCertificatesCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List organization certificates",
		Flags: []cli.Flag{
			limitFlag,
			afterFlag,
			&cli.StringFlag{
				Name:  "order",
				Usage: "Sort order (asc or desc)",
				Value: "desc",
			},
		},
		Action: listOrgCertificates,
	}
}

func uploadCertificateCommand() *cli.Command {
	return &cli.Command{
		Name:  "upload",
		Usage: "Upload a new certificate",
		Flags: []cli.Flag{
			nameFlag,
			&cli.StringFlag{
				Name:     "content",
				Usage:    "PEM-encoded certificate content",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "content-file",
				Usage: "Path to file containing PEM-encoded certificate content",
			},
		},
		Action: uploadCertificate,
	}
}

func getCertificateCommand() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Get a specific certificate",
		Flags: []cli.Flag{
			idFlag,
			&cli.BoolFlag{
				Name:  "include-content",
				Usage: "Include certificate PEM content in response",
			},
		},
		Action: getCertificate,
	}
}

func modifyCertificateCommand() *cli.Command {
	return &cli.Command{
		Name:  "modify",
		Usage: "Modify a certificate name",
		Flags: []cli.Flag{
			idFlag,
			nameFlag,
		},
		Action: modifyCertificate,
	}
}

func deleteCertificateCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a certificate",
		Flags: []cli.Flag{
			idFlag,
		},
		Action: deleteCertificate,
	}
}

func activateOrgCertificatesCommand() *cli.Command {
	return &cli.Command{
		Name:  "activate",
		Usage: "Activate multiple certificates",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "certificate-ids",
				Usage:    "Comma-separated list of certificate IDs to activate",
				Required: true,
			},
		},
		Action: activateOrgCertificates,
	}
}

func deactivateOrgCertificatesCommand() *cli.Command {
	return &cli.Command{
		Name:  "deactivate",
		Usage: "Deactivate multiple certificates",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "certificate-ids",
				Usage:    "Comma-separated list of certificate IDs to deactivate",
				Required: true,
			},
		},
		Action: deactivateOrgCertificates,
	}
}

// Project certificate commands

func listProjectCertificatesCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List project certificates",
		Flags: []cli.Flag{
			projectIDFlag,
			limitFlag,
			afterFlag,
			&cli.StringFlag{
				Name:  "order",
				Usage: "Sort order (asc or desc)",
				Value: "desc",
			},
		},
		Action: listProjectCertificates,
	}
}

func activateProjectCertificatesCommand() *cli.Command {
	return &cli.Command{
		Name:  "activate",
		Usage: "Activate certificates for a project",
		Flags: []cli.Flag{
			projectIDFlag,
			&cli.StringSliceFlag{
				Name:     "certificate-ids",
				Usage:    "Comma-separated list of certificate IDs to activate",
				Required: true,
			},
		},
		Action: activateProjectCertificates,
	}
}

func deactivateProjectCertificatesCommand() *cli.Command {
	return &cli.Command{
		Name:  "deactivate",
		Usage: "Deactivate certificates for a project",
		Flags: []cli.Flag{
			projectIDFlag,
			&cli.StringSliceFlag{
				Name:     "certificate-ids",
				Usage:    "Comma-separated list of certificate IDs to deactivate",
				Required: true,
			},
		},
		Action: deactivateProjectCertificates,
	}
}

// Action handlers

func listOrgCertificates(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	limit := cmd.Int("limit")
	after := cmd.String("after")
	order := cmd.String("order")

	certificates, err := client.ListOrganizationCertificates(limit, after, order)
	if err != nil {
		return wrapError("list organization certificates", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Active", "Valid At", "Expires At"},
		Rows:    make([][]string, len(certificates.Data)),
	}

	for i, cert := range certificates.Data {
		active := "N/A"
		if cert.Active != nil {
			if *cert.Active {
				active = "Yes"
			} else {
				active = "No"
			}
		}
		data.Rows[i] = []string{
			cert.ID,
			cert.Name,
			active,
			cert.CertificateDetails.ValidAt.String(),
			cert.CertificateDetails.ExpiresAt.String(),
		}
	}

	printTableData(data)
	return nil
}

func uploadCertificate(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	name := cmd.String("name")
	content := cmd.String("content")
	contentFile := cmd.String("content-file")

	if content == "" && contentFile == "" {
		return fmt.Errorf("either --content or --content-file must be provided")
	}

	if content != "" && contentFile != "" {
		return fmt.Errorf("only one of --content or --content-file can be provided")
	}

	if contentFile != "" {
		data, err := os.ReadFile(contentFile)
		if err != nil {
			return fmt.Errorf("failed to read certificate file: %v", err)
		}
		content = string(data)
	}

	certificate, err := client.UploadCertificate(content, name)
	if err != nil {
		return wrapError("upload certificate", err)
	}

	fmt.Printf("Certificate uploaded:\n")
	fmt.Printf("ID: %s\nName: %s\nValid At: %s\nExpires At: %s\n",
		certificate.ID,
		certificate.Name,
		certificate.CertificateDetails.ValidAt.String(),
		certificate.CertificateDetails.ExpiresAt.String())
	return nil
}

func getCertificate(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	id := cmd.String("id")
	includeContent := cmd.Bool("include-content")

	certificate, err := client.GetCertificate(id, includeContent)
	if err != nil {
		return wrapError("get certificate", err)
	}

	fmt.Printf("Certificate details:\n")
	fmt.Printf("ID: %s\nName: %s\nValid At: %s\nExpires At: %s\n",
		certificate.ID,
		certificate.Name,
		certificate.CertificateDetails.ValidAt.String(),
		certificate.CertificateDetails.ExpiresAt.String())
	if certificate.CertificateDetails.Content != nil {
		fmt.Printf("Content:\n%s\n", *certificate.CertificateDetails.Content)
	}
	return nil
}

func modifyCertificate(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	id := cmd.String("id")
	name := cmd.String("name")

	certificate, err := client.ModifyCertificate(id, name)
	if err != nil {
		return wrapError("modify certificate", err)
	}

	fmt.Printf("Certificate modified:\n")
	fmt.Printf("ID: %s\nName: %s\nValid At: %s\nExpires At: %s\n",
		certificate.ID,
		certificate.Name,
		certificate.CertificateDetails.ValidAt.String(),
		certificate.CertificateDetails.ExpiresAt.String())
	return nil
}

func deleteCertificate(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	id := cmd.String("id")

	response, err := client.DeleteCertificate(id)
	if err != nil {
		return wrapError("delete certificate", err)
	}

	if response.Deleted {
		fmt.Printf("Certificate %s deleted successfully\n", response.ID)
	} else {
		fmt.Printf("Failed to delete certificate %s\n", response.ID)
	}
	return nil
}

func activateOrgCertificates(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	certificateIDs := cmd.StringSlice("certificate-ids")
	if len(certificateIDs) == 0 {
		return fmt.Errorf("at least one certificate ID must be provided")
	}

	// Handle comma-separated values in a single string
	var allIDs []string
	for _, id := range certificateIDs {
		if strings.Contains(id, ",") {
			allIDs = append(allIDs, strings.Split(id, ",")...)
		} else {
			allIDs = append(allIDs, id)
		}
	}

	response, err := client.ActivateOrganizationCertificates(allIDs)
	if err != nil {
		return wrapError("activate certificates", err)
	}

	if response.Success {
		fmt.Printf("Successfully activated %d certificates\n", len(allIDs))
	} else {
		fmt.Printf("Failed to activate certificates\n")
	}
	return nil
}

func deactivateOrgCertificates(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	certificateIDs := cmd.StringSlice("certificate-ids")
	if len(certificateIDs) == 0 {
		return fmt.Errorf("at least one certificate ID must be provided")
	}

	// Handle comma-separated values in a single string
	var allIDs []string
	for _, id := range certificateIDs {
		if strings.Contains(id, ",") {
			allIDs = append(allIDs, strings.Split(id, ",")...)
		} else {
			allIDs = append(allIDs, id)
		}
	}

	response, err := client.DeactivateOrganizationCertificates(allIDs)
	if err != nil {
		return wrapError("deactivate certificates", err)
	}

	if response.Success {
		fmt.Printf("Successfully deactivated %d certificates\n", len(allIDs))
	} else {
		fmt.Printf("Failed to deactivate certificates\n")
	}
	return nil
}

func listProjectCertificates(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	projectID := cmd.String("project-id")
	limit := cmd.Int("limit")
	after := cmd.String("after")
	order := cmd.String("order")

	certificates, err := client.ListProjectCertificates(projectID, limit, after, order)
	if err != nil {
		return wrapError("list project certificates", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Active", "Valid At", "Expires At"},
		Rows:    make([][]string, len(certificates.Data)),
	}

	for i, cert := range certificates.Data {
		active := "N/A"
		if cert.Active != nil {
			if *cert.Active {
				active = "Yes"
			} else {
				active = "No"
			}
		}
		data.Rows[i] = []string{
			cert.ID,
			cert.Name,
			active,
			cert.CertificateDetails.ValidAt.String(),
			cert.CertificateDetails.ExpiresAt.String(),
		}
	}

	printTableData(data)
	return nil
}

func activateProjectCertificates(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	projectID := cmd.String("project-id")
	certificateIDs := cmd.StringSlice("certificate-ids")
	if len(certificateIDs) == 0 {
		return fmt.Errorf("at least one certificate ID must be provided")
	}

	// Handle comma-separated values in a single string
	var allIDs []string
	for _, id := range certificateIDs {
		if strings.Contains(id, ",") {
			allIDs = append(allIDs, strings.Split(id, ",")...)
		} else {
			allIDs = append(allIDs, id)
		}
	}

	response, err := client.ActivateProjectCertificates(projectID, allIDs)
	if err != nil {
		return wrapError("activate project certificates", err)
	}

	if response.Success {
		fmt.Printf("Successfully activated %d certificates for project %s\n", len(allIDs), projectID)
	} else {
		fmt.Printf("Failed to activate certificates for project %s\n", projectID)
	}
	return nil
}

func deactivateProjectCertificates(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	projectID := cmd.String("project-id")
	certificateIDs := cmd.StringSlice("certificate-ids")
	if len(certificateIDs) == 0 {
		return fmt.Errorf("at least one certificate ID must be provided")
	}

	// Handle comma-separated values in a single string
	var allIDs []string
	for _, id := range certificateIDs {
		if strings.Contains(id, ",") {
			allIDs = append(allIDs, strings.Split(id, ",")...)
		} else {
			allIDs = append(allIDs, id)
		}
	}

	response, err := client.DeactivateProjectCertificates(projectID, allIDs)
	if err != nil {
		return wrapError("deactivate project certificates", err)
	}

	if response.Success {
		fmt.Printf("Successfully deactivated %d certificates for project %s\n", len(allIDs), projectID)
	} else {
		fmt.Printf("Failed to deactivate certificates for project %s\n", projectID)
	}
	return nil
}
