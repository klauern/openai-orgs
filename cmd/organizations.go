package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/urfave/cli/v3"
)

func OrganizationsCommand() *cli.Command {
	return &cli.Command{
		Name:  "organizations",
		Usage: "Manage organizations",
		Commands: []*cli.Command{
			listOrganizationsCommand(),
			getOrganizationCommand(),
			updateOrganizationCommand(),
			deleteOrganizationCommand(),
		},
	}
}

func listOrganizationsCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List organizations",
		Flags: []cli.Flag{
			limitFlag,
			afterFlag,
		},
		Action: listOrganizations,
	}
}

func getOrganizationCommand() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Get organization details",
		Flags: []cli.Flag{
			idFlag,
		},
		Action: getOrganization,
	}
}

func updateOrganizationCommand() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Update organization settings",
		Flags: []cli.Flag{
			idFlag,
			&cli.StringFlag{
				Name:  "name",
				Usage: "Organization name",
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Organization description",
			},
			&cli.StringFlag{
				Name:  "billing-address",
				Usage: "Default billing address (JSON format)",
			},
		},
		Action: updateOrganization,
	}
}

func deleteOrganizationCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete organization",
		Flags: []cli.Flag{
			idFlag,
		},
		Action: deleteOrganization,
	}
}

func listOrganizations(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	limit := int(cmd.Int("limit"))
	orgs, err := client.ListOrganizations(
		limit,
		cmd.String("after"),
	)
	if err != nil {
		return wrapError("list organizations", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Description", "Created", "Blocked"},
		Rows:    make([][]string, len(orgs.Data)),
	}

	for i, org := range orgs.Data {
		data.Rows[i] = []string{
			org.ID,
			org.Name,
			org.Description,
			org.CreatedAt.String(),
			fmt.Sprintf("%v", org.IsBlocked),
		}
	}

	printTableData(data)
	return nil
}

func getOrganization(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	org, err := client.GetOrganization(cmd.String("id"))
	if err != nil {
		return wrapError("get organization", err)
	}

	fmt.Printf("Organization details:\n")
	fmt.Printf("ID: %s\nName: %s\nDescription: %s\nCreated: %s\nBlocked: %v\n",
		org.ID,
		org.Name,
		org.Description,
		org.CreatedAt.String(),
		org.IsBlocked,
	)

	if org.Settings.DefaultBillingAddress != nil {
		fmt.Printf("\nDefault Billing Address:\n")
		fmt.Printf("Line1: %s\n", org.Settings.DefaultBillingAddress.Line1)
		if org.Settings.DefaultBillingAddress.Line2 != "" {
			fmt.Printf("Line2: %s\n", org.Settings.DefaultBillingAddress.Line2)
		}
		fmt.Printf("City: %s\nState: %s\nCountry: %s\nPostal Code: %s\n",
			org.Settings.DefaultBillingAddress.City,
			org.Settings.DefaultBillingAddress.State,
			org.Settings.DefaultBillingAddress.Country,
			org.Settings.DefaultBillingAddress.PostalCode,
		)
	}

	return nil
}

func updateOrganization(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	req := &openaiorgs.UpdateOrganizationRequest{}

	if cmd.IsSet("name") {
		name := cmd.String("name")
		req.Name = &name
	}

	if cmd.IsSet("description") {
		desc := cmd.String("description")
		req.Description = &desc
	}

	if cmd.IsSet("billing-address") {
		var address openaiorgs.BillingAddress
		if err := json.Unmarshal([]byte(cmd.String("billing-address")), &address); err != nil {
			return wrapError("parse billing address", err)
		}
		req.Settings = &openaiorgs.OrgSettings{
			DefaultBillingAddress: &address,
		}
	}

	org, err := client.UpdateOrganization(cmd.String("id"), req)
	if err != nil {
		return wrapError("update organization", err)
	}

	fmt.Printf("Organization updated successfully:\n")
	fmt.Printf("ID: %s\nName: %s\nDescription: %s\nCreated: %s\nBlocked: %v\n",
		org.ID,
		org.Name,
		org.Description,
		org.CreatedAt.String(),
		org.IsBlocked,
	)

	return nil
}

func deleteOrganization(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	if err := client.DeleteOrganization(cmd.String("id")); err != nil {
		return wrapError("delete organization", err)
	}

	fmt.Printf("Organization %s deleted successfully\n", cmd.String("id"))
	return nil
}
