package cmd

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/cobra"
)

var (
	username string
	password string
)

var rootCmd = &cobra.Command{
	Use:   "monoize [source repository] [target repository]",
	Short: "monoize makes your git repositories monorepo",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sources := args[:len(args)-1]
		// target := args[len(args)-1]

		for _, source := range sources {
			repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
				Auth: &http.BasicAuth{
					Username: username, Password: password,
				},
				URL: source,
			})
			if err != nil {
				panic(err)
			}

			ref, err := repo.Head()
			if err != nil {
				panic(err)
			}

			fmt.Println(ref)

			cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
			if err != nil {
				panic(err)
			}

			cIter.ForEach(func(c *object.Commit) error {
				fmt.Println(c)
				return nil
			})
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&username, "username", "u", "", "username for auth")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "password for auth")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
