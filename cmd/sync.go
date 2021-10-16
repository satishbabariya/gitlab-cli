/*
Copyright © 2021 Satish Babariya

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
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	. "github.com/logrusorgru/aurora"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Clone all repositories from gitlab to current directory",
	Long:  `Clone all repositories from gitlab to current directory`,
	Run: func(cmd *cobra.Command, args []string) {

		cloneAllProjects()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringVarP(&token, "token", "t", "", "Enter your personal access token")
	syncCmd.MarkFlagRequired("token")
}

func cloneAllProjects() {

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
		fmt.Println(Bold(Red("No Projects Found")))
		// fmt.Println("No Projects Found")
		os.Exit(0)
	}

	fmt.Println(Bold(Cyan("Cloning all your projects ")))

	for _, project := range projects {
		clone(project)
	}
}

func clone(project *gitlab.Project) {
	dir, err := os.Getwd()

	//check if you need to panic, fallback or report
	if err != nil {
		fmt.Println(Bold(Red("Error Cloning: " + err.Error())))
		return
	}

	_, err = git.PlainClone(dir+"/"+project.Name, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "abc123", // yes, this can be anything except an empty string
			Password: token,
		},
		URL:      project.HTTPURLToRepo,
		Progress: os.Stdout,
	})

	if err != nil {
		fmt.Println(Bold(Red("Error Cloning " + project.Name + ": " + err.Error())))
	}
}
