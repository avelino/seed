package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	_, err := exec.LookPath("go")
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Version = "0.1"
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "Install packages",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Value: "Seedfile",
					Usage: "Install from the given Seeds file. This option can be used multiple times",
				},
			},
			Action: func(c *cli.Context) error {
				Seedfile := c.String("file")

				file, err := os.Open(Seedfile)
				defer file.Close()
				if err != nil {
					return err
				}

				scanner := bufio.NewScanner(file)
				scanner.Split(bufio.ScanLines)
				for scanner.Scan() {
					repo := scanner.Text()
					fmt.Println("get: ", repo)
					args := []string{"get", "-u", repo}
					err := exec.Command("go", args...).Run()
					if err != nil {
						fmt.Println(err)
					}
				}
				return nil
			},
		},
		{
			Name:    "freeze",
			Aliases: []string{"f"},
			Usage:   "",
			Action: func(c *cli.Context) error {
				args := []string{"list", "-f", `'{{ join .Imports "\n" }}'`}
				outPut, err := exec.Command("go", args...).Output()
				if err != nil {
					fmt.Println(err)
				}

				clear := strings.Replace(string(outPut), `'`, "", -1)
				packages := strings.Split(clear, "\n")
				for _, p := range packages {
					if p != "" {
						fmt.Println(p)
					}

				}
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}
