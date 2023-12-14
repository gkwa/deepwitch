/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		summary()
	},
}

const TmpDir = "tmp"

var TmpDirAbsPath string

func init() {
	rootCmd.AddCommand(summaryCmd)
	var err error
	TmpDirAbsPath, err = filepath.Abs(TmpDir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

const orgTemplate = `
** example

#+begin_example
{{ .Contents }}
#+end_example
`

type templateData struct {
	ExampleName string
	Contents    string
}

func summary() {
	dir := TmpDirAbsPath

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Error: The %s directory does not exist. Follow the usage instructions below:", TmpDirAbsPath)
		printHelpMessage()
		return
	}

	files, err := findMainGoFiles(dir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	sort.Strings(files)

	for _, file := range files {
		contents, err := readFileContents(file)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", file, err)
			continue
		}

		exampleName := filepath.Base(filepath.Dir(file))

		data := templateData{
			ExampleName: exampleName,
			Contents:    contents,
		}

		err = renderTemplate(orgTemplate, data)
		if err != nil {
			fmt.Printf("Error rendering template for %s: %v\n", file, err)
		}

	}
}

func findMainGoFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Base(path) != "main.go" {
			return nil
		}
		files = append(files, path)
		return nil
	})

	return files, err
}

func readFileContents(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var builder strings.Builder
	_, err = io.Copy(&builder, f)
	if err != nil {
		return "", err
	}

	contents := strings.TrimSpace(builder.String())

	return contents, nil
}

func renderTemplate(tmpl string, data templateData) error {
	t, err := template.New("orgTemplate").Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(os.Stdout, data)
}

func printHelpMessage() {
	fmt.Println(`
# usage:
cd ~/pdev/taylormonacelli/deepwitch/
rm -rf tmp && mkdir tmp
txtar-x -C tmp example.txtar
go build
./deepwitch summary
./deepwitch summary | pbcopy
	`)
}
