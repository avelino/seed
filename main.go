package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/urfave/cli"
)

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

func copyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		_ = os.Remove(dst)
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			switch entry.Name() {
			case
				".git",
				"vendor",
				".github":
				continue
			}
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// not exist .go on file name continue
			if !strings.Contains(entry.Name(), ".go") {
				continue
			}
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = copyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

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
				cli.StringFlag{
					Name:  "folder, dir, d",
					Value: "vendor",
					Usage: "Install packages list by Seedfile on seed (or vendor) folder.",
				},
			},
			Action: func(c *cli.Context) error {
				ProjectFolder, _ := os.Getwd()
				SeedFile := c.String("file")
				SeedFolder := c.String("folder")

				file, err := os.Open(SeedFile)
				defer file.Close()
				if err != nil {
					return err
				}

				scanner := bufio.NewScanner(file)
				scanner.Split(bufio.ScanLines)
				for scanner.Scan() {
					repo := strings.Split(scanner.Text(), "@")
					branch := "master"
					if len(repo) == 2 {
						branch = repo[1]
					}
					fmt.Println("get: ", repo[0], " branch/commit: ", branch)

					args := []string{"get", "-u", repo[0]}
					_ = exec.Command("go", args...).Run()

					repoFolder := fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), repo[0])
					if branch != "master" {
						if err := os.Chdir(repoFolder); err != nil {
							fmt.Println("folder not exist!")
						}

						err := exec.Command("git", []string{"checkout", branch}...).Run()
						if err != nil {
							fmt.Println(err)
						}
					}

					SeedPath := SeedFolder
					if SeedFolder == "vendor" {
						SeedPath = fmt.Sprintf("%s/%s", ProjectFolder, SeedFolder)
					}
					// create SeedPath dir
					_ = os.MkdirAll(SeedPath, os.ModePerm)
					// sync folder
					dstPath := fmt.Sprintf("%s/%s", SeedPath, repo[0])
					_ = copyDir(repoFolder, dstPath)

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
