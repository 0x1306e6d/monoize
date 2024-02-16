package cmd

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "monoize",
	Short: "monoize makes your git repositories monorepo",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL: "https://github.com/go-git/go-billy",
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
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
