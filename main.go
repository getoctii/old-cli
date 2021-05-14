package main

import (
	"github.com/urfave/cli/v2"
	"os"
	apps "github.com/innatical/octii-cli/apps"
)


func main() {
	app := &cli.App {
		Name: "octii",
		Usage: "Octii DevTools",
		UsageText: "octii [global options] command [command options] [arguments...]",
		Commands: []*cli.Command {
			{
				Name: "account",
				Usage: "Manage your Octii account",
				Subcommands: []*cli.Command {
					{
						Name: "login",
						Usage: "Login into Octii",
						UsageText: "octii account login",
						Action: apps.Login,
					},
					{
						Name: "info",
						Usage: "Show account info for the logged in user",
						UsageText: "octii account info",
						Action: apps.Account,
					},
				},
			},
			{
				Name: "organizations",
				Usage: "Manage your developer organizations",
				Subcommands: []*cli.Command{
					{
						Name: "list",
						Usage: "List your Octii organizations",
						UsageText: "octii organizations list",
						Action: apps.Organizations,
					},
					{
						Name: "products",
						Usage: "Get the products belonging to an organization",
						UsageText: "octii organizations products <id>",
						Subcommands: []*cli.Command{
							{
								Name:      "list",
								Usage:     "List all products",
								UsageText: "octii organizations products list <organization_id>",
								Action:    apps.Products,
							},
						},
					},
				},
			},
			{
				Name: "products",
				Usage: "Manage your Octii products",
				Subcommands: []*cli.Command{
					{
						Name: "resources",
						Usage: "Manage your Octii resources",
						Subcommands: []*cli.Command{
							{
								Name: "get",
								Usage: "Get the payload of a resource",
								UsageText: "octii products resources get <product_id> <resource_id> <output_file>",
								Action: apps.GetResource,
							},
							{
								Name: "put",
								Usage: "Put the payload of a resource",
								UsageText: "octii products resources put <product_id> <resource_id> <input_file>",
								Action: apps.PutResource,
							},
							{
								Name: "list",
								Usage: "Get all resources for a product",
								UsageText: "octii products resources list <product_id>",
								Action: apps.Resources,
							},
						},
					},
				},
			},
			{
				Name: "templates",
				Usage: "Scaffold your product from a template",
				Subcommands: []*cli.Command {
					{
						Name: "theme",
						Usage: "Write a theme template",
						UsageText: "octii templates theme <output_file>",
						Action: apps.Theme,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
