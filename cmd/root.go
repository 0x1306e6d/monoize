package cmd

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

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

		for _, source := range sources {
			u, err := url.Parse(source)
			if err != nil {
				return fmt.Errorf("`http` and `https` protocols are supported")
			}
			if u.Scheme != "http" && u.Scheme != "https" {
				return fmt.Errorf("`http` and `https` protocols are supported")
			}

			b := path.Base(u.Path)
			name := strings.TrimSuffix(b, ".git")

			cmd := exec.Command("git", "clone", source, fmt.Sprintf(".repo/%s", name))
			cmd.Dir = target
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}

			cmd = exec.Command("git", "format-patch", "--root", "-o", fmt.Sprintf("../../.patch/%s", name))
			cmd.Dir = fmt.Sprintf("%s/.repo/%s", target, name)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
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
