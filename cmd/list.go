/*
Copyright © 2019 Satish Babariya

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	"github.com/jedib0t/go-pretty/list"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all your projects",
	Long:  `Get a list of all visible projects across GitLab for the authenticated user. When accessed without authentication, only public projects with “simple” fields are returned.`,
	Run: func(cmd *cobra.Command, args []string) {

		listAllProjects()
	},
}

var token string

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&token, "token", "t", "", "Enter your personal access token")
	listCmd.MarkFlagRequired("token")
}

func listAllProjects() {
	var client = gitlab.NewClient(nil, token)
	opt := &gitlab.ListProjectsOptions{
		Owned:      gitlab.Bool(true),
		Membership: gitlab.Bool(true),
	}
	projects, _, err := client.Projects.ListProjects(opt)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(projects) == 0 {
		fmt.Println("No Projects Found")
		os.Exit(0)
	}

	l := list.NewWriter()
	for _, project := range projects {
		l.AppendItem(project.Name)
	}
	l.SetStyle(list.StyleBulletCircle)
	fmt.Println("\n")
	consoleLog("List all your projects", l.Render(), "")
}
