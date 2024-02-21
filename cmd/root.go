package cmd

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/cobra"
)

var (
	username string
	password string
	force    bool
)

var rootCmd = &cobra.Command{
	Use:   "monoize [source repository] [target repository]",
	Short: "monoize makes your git repositories monorepo",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		sources := args[:len(args)-1]
		target := args[len(args)-1]

		for _, source := range sources {
			u, err := url.Parse(source)
			if err != nil {
				return fmt.Errorf("`http` and `https` protocols are supported")
			}
			if u.Scheme != "http" && u.Scheme != "https" {
				return fmt.Errorf("`http` and `https` protocols are supported")
			}
		}

		if force {
			err := os.RemoveAll(target)
			if err != nil {
				return err
			}
		}
		err := os.Mkdir(target, os.ModePerm)
		if err != nil {
			if os.IsExist(err) {
				return fmt.Errorf("`%s` already exists", target)
			}
			return err
		}

		commits := []commitSrc{}
		for _, src := range sources {
			repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
				Auth: &http.BasicAuth{
					Username: username, Password: password,
				},
				URL: src,
			})
			if err != nil {
				return err
			}

			cIter, err := repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime, All: true})
			if err != nil {
				return err
			}

			for {
				c, err := cIter.Next()
				if err == io.EOF {
					break
				}
				commits = append(commits, commitSrc{src, *c})
			}
		}

		sort.Sort(byCTime(commits))

		for _, c := range commits {
			fmt.Printf("%s: [%s] %s\n", c.src, &c.commit.Committer.When, c.commit.Message)
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&username, "username", "u", "", "username for auth")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "password for auth")
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "force to overwrite the target directory")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type commitSrc struct {
	src    string
	commit object.Commit
}

func (c commitSrc) String() string {
	return fmt.Sprintf("%s: %s", c.src, c.commit.Hash)
}

type byCTime []commitSrc

func (c byCTime) Len() int { return len(c) }
func (c byCTime) Less(i, j int) bool {
	return c[i].commit.Committer.When.Before(c[j].commit.Committer.When)
}
func (c byCTime) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
