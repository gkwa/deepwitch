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

func init() {
	rootCmd.AddCommand(summaryCmd)
}

const orgTemplate = `
** example

#+begin_example{{"\n"}}{{- .Contents -}}{{"\n"}}#+end_example
`

type templateData struct {
	ExampleName string
	Contents    string
}

func summary() {
	// Define the directory path
	dir := "/tmp/new"

	// Find all main.go files in the specified directory
	files, err := findMainGoFiles(dir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Sort the file list for consistent output
	sort.Strings(files)

	// Loop through each file and print its contents in org mode format
	for _, file := range files {
		contents, err := readFileContents(file)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", file, err)
			continue
		}

		// Extract example name from the file path
		exampleName := filepath.Base(filepath.Dir(file))

		// Render the template
		data := templateData{
			ExampleName: exampleName,
			Contents:    contents,
		}

		// Output in org mode format
		err = renderTemplate(orgTemplate, data)
		if err != nil {
			fmt.Printf("Error rendering template for %s: %v\n", file, err)
		}

	}
}

// findMainGoFiles finds all main.go files in the specified directory and its subdirectories.
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

// readFileContents reads the contents of a file and returns them as a string.
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

	// Trim leading and trailing whitespaces
	contents := strings.TrimSpace(builder.String())

	return contents, nil
}

// renderTemplate renders the Go template with the provided data.
func renderTemplate(tmpl string, data templateData) error {
	t, err := template.New("orgTemplate").Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(os.Stdout, data)
}
