package main

import (
	"os"
	"devdock/docker"
	"devdock/configs"

	"github.com/urfave/cli"
	"github.com/olekukonko/tablewriter"
	"fmt"
)

func main() {
	app := cli.NewApp()
	app.Name = "Dev Dock"
	app.Version = "0.0.1"
	app.Usage = "Organize your Docker DEV Containers"

	app.Commands = []cli.Command{
		//{
		//	Name: "new",
		//	Aliases: []string{"n"},
		//	Usage: "Create a new project",
		//	UsageText: "new name - Create a new project",
		//	ArgsUsage: "[name]",
		//	Flags: []cli.Flag{
		//		cli.BoolFlag{Name: "compose, c"},
		//	},
		//	Action: func(c *cli.Context) error {
		//		return nil
		//	},
		//},
		{
			Name: "list",
			Aliases: []string{"l"},
			Usage: "List all projects",
			Action: func(c *cli.Context) error {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"NAME", "DOMAIN", "IMAGE", "STATUS", "VOLUMES", "PORTS"})

				confs := configs.Read()
				for _, config := range confs.Projects {
					table.Append(config.ToSlice())
				}

				table.Render()

				return nil
			},
		},
		{
			Name: "up",
			Aliases: []string{"u"},
			Usage: "Start a project",
			UsageText: "new name - Create a new project",
			ArgsUsage: "[name]",
			Action: func(c *cli.Context) error {
				name := c.Args().First()
				if (name == "") {
					fmt.Println("Project name require")
					return nil
				}

				project := configs.Find(name)
				if (project == nil) {
					fmt.Printf("Project \"%s\" not found\n", name)
					return nil;
				}

				fmt.Printf("Starting \"%v\"...\n\n", name)

				//docker.StartProxyContainer()
				docker.StartProjectContainer(*project)

				return nil
			},
		},
		{
			Name: "down",
			Aliases: []string{"d"},
			Usage: "Finish a started project",
			Action: func(c *cli.Context) error {
				name := c.Args().First()
				if (name == "") {
					fmt.Println("Project name require")
					return nil
				}

				project := configs.Find(name)
				if (project == nil) {
					fmt.Printf("Project \"%s\" not found\n", name)
					return nil;
				}

				fmt.Printf("Finishing \"%v\"...\n\n", name)

				docker.FinishProjectContainer(*project)

				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c)
		return nil
	}

	if !configs.Exists() {
		configs.Create()
	}

	app.Run(os.Args)
}


