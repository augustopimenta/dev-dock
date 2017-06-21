package main

import (
	"os"
	"fmt"
	"strings"
	"os/user"

	"devdock/docker"
	"devdock/configs"

	"github.com/urfave/cli"
	"github.com/olekukonko/tablewriter"
	"github.com/lextoumbourou/goodhosts"
)

func main() {
	app := cli.NewApp()
	app.Name = "Dev Dock"
	app.Version = "0.1.2"
	app.Usage = "Organize your Docker Development Containers"

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
			Usage: "List all configs",
			Action: handlerListAction,
		},
		{
			Name: "up",
			Aliases: []string{"u"},
			Usage: "Start a project",
			UsageText: "new name - Create a new project",
			ArgsUsage: "[name]",
			Action: handlerUpAction,
		},
		{
			Name: "down",
			Aliases: []string{"d"},
			Usage: "Finish a started project",
			Action: handleDownAction,
		},
	}

	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c)
		return nil
	}

	checkForRootAccess()

	updateHosts()

	app.Run(os.Args)
}

func handlerListAction(c *cli.Context) error {
	confs := configs.NewConfigFile()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "DOMAIN", "IMAGE", "STATUS", "VOLUMES", "PORTS"})

	for _, project := range confs.Projects {
		if project.Name == configs.ExampleProjectName {
			continue
		}

		container := docker.GetProjectContainer(project);

		if container != nil {
			project.Status = strings.ToUpper(container.State.Status)
		} else {
			project.Status = "DOWN"
		}

		table.Append(project.ToSlice())
	}

	table.Render()

	return nil
}

func handlerUpAction(c *cli.Context) error {
	name := c.Args().First()
	if (name == "") {
		fmt.Println("Project name require")
		return nil
	}

	conf := configs.NewConfigFile()
	project := conf.FindProject(name)
	if (project == nil) {
		fmt.Printf("Project \"%s\" not found\n", name)
		return nil;
	}

	fmt.Printf("Starting \"%v\"...\n\n", name)

	if conf.UseVirtualHost {
		docker.StartProxyContainer()
	} else {
		docker.FinishProxyContainer()
	}

	docker.StartProjectContainer(*project)

	return nil
}

func handleDownAction(c *cli.Context) error {
	name := c.Args().First()
	if (name == "") {
		fmt.Println("Project name require")
		return nil
	}

	project := configs.NewConfigFile().FindProject(name)
	if (project == nil) {
		fmt.Printf("Project \"%s\" not found\n", name)
		return nil;
	}

	fmt.Printf("Finishing \"%v\"...\n\n", name)

	docker.FinishProjectContainer(*project)

	return nil
}

func checkForRootAccess() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	if user.Uid != "0" {
		fmt.Println("Error: You need root access to manage Docker containers and hosts file!")
		os.Exit(0)
	}
}

func updateHosts() {
	hosts, err := goodhosts.NewHosts()

	if err != nil {
		panic(err)
	}

	confs := configs.NewConfigFile()

	for _, project := range confs.Projects {
		if project.Name == configs.ExampleProjectName {
			continue
		}
		hosts.Add("127.0.0.1", project.Domain)
	}

	hosts.Flush()
}

