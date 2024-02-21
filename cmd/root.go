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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	force bool
)

func init() {
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "force to overwrite the target directory")
}

var rootCmd = &cobra.Command{
	Use:   "monoize [srcRepo{>>targetDir}] [target repository]",
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
			if !os.IsExist(err) {
				return err
			}
		}

		if !exists(fmt.Sprintf("%s/.git", target)) {
			if err := git.Init(target); err != nil {
				return err
			}
		}

		patches := []patchFile{}
		for _, src := range sources {
			srcRepo, srcRepoName, targetDir, err := parseSrcRepo(src)
			if err != nil {
				return err
			}

			if err := git.Clone(target, srcRepo, fmt.Sprintf(".repo/%s", srcRepoName)); err != nil {
				return err
			}

			if err := git.FormatPatch(fmt.Sprintf("%s/.repo/%s", target, srcRepoName), fmt.Sprintf("../../.patch/%s", srcRepoName)); err != nil {
				return err
			}

			dir := os.DirFS(fmt.Sprintf("%s/.patch/%s", target, srcRepoName))
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

				path := filepath.Join(target, ".patch", srcRepoName, e.Name())
				path, err = filepath.Abs(path)
				if err != nil {
					return nil
				}

				pf := patchFile{
					directory: targetDir,
					path:      path,
					patch:     p,
				}
				patches = append(patches, pf)
			}
		}

		sort.Sort(byPatchDate(patches))

		for _, p := range patches {
			git.Am(target, p.path, p.directory)
			fmt.Printf("[%s] Applying %s to %s\n", p.patch.Date, p.patch.Subject, p.directory)
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

func parseSrcRepo(src string) (string, string, string, error) {
	var srcRepo, srcRepoName, targetDir string
	split := strings.Split(src, ">>")
	if len(split) == 1 {
		srcRepo = split[0]
	} else if len(split) == 2 {
		srcRepo, targetDir = split[0], split[1]
	} else {
		return "", "", "", fmt.Errorf("wrong source repository")
	}

	u, err := url.Parse(srcRepo)
	if err != nil {
		return "", "", "", fmt.Errorf("wrong source repository")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", "", "", fmt.Errorf("wrong source repository")
	}

	b := path.Base(u.Path)
	srcRepoName = strings.TrimSuffix(b, ".git")

	return srcRepo, srcRepoName, targetDir, nil
}

type patchFile struct {
	directory string
	path      string
	patch     git.Patch
}

type byPatchDate []patchFile

func (b byPatchDate) Len() int           { return len(b) }
func (b byPatchDate) Less(i, j int) bool { return b[i].patch.Date.Before(b[j].patch.Date) }
func (b byPatchDate) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func exists(name string) bool {
	_, err := os.Stat(name)
	return os.IsExist(err)
}
