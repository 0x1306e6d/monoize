package cmd

import (
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/0x1306e6d/monoize/pkg/git"
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

		if err := git.Init(target); err != nil {
			return err
		}

		patches := []patchFile{}
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

			if err := git.Clone(target, source, fmt.Sprintf(".repo/%s", name)); err != nil {
				return err
			}

			if err := git.FormatPatch(fmt.Sprintf("%s/.repo/%s", target, name), fmt.Sprintf("../../.patch/%s", name)); err != nil {
				return err
			}

			dir := os.DirFS(fmt.Sprintf("%s/.patch/%s", target, name))
			entries, err := fs.ReadDir(dir, ".")
			if err != nil {
				return err
			}

			for _, e := range entries {
				d, err := fs.ReadFile(dir, e.Name())
				if err != nil {
					return err
				}

				p, err := git.ParsePatch(string(d))
				if err != nil {
					return err
				}

				path := filepath.Join(target, ".patch", name, e.Name())
				path, err = filepath.Abs(path)
				if err != nil {
					return nil
				}

				pf := patchFile{
					Repository: name,
					Path:       path,
					Patch:      p,
				}
				patches = append(patches, pf)
			}
		}

		sort.Sort(byPatchDate(patches))

		for _, p := range patches {
			git.Am(target, p.Path, p.Repository)
			fmt.Printf("[%s] Applying %s to %s\n", p.Patch.Date, p.Patch.Subject, p.Repository)
		}

		p := filepath.Join(target, ".repo")
		if err := os.RemoveAll(p); err != nil {
			return err
		}

		p = filepath.Join(target, ".patch")
		if err := os.RemoveAll(p); err != nil {
			return err
		}

		return nil
	},
}

type patchFile struct {
	Repository string
	Path       string
	Patch      git.Patch
}

type byPatchDate []patchFile

func (b byPatchDate) Len() int           { return len(b) }
func (b byPatchDate) Less(i, j int) bool { return b[i].Patch.Date.Before(b[j].Patch.Date) }
func (b byPatchDate) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

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
