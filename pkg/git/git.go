package git

import (
	"os"
	"os/exec"
)

func Clone(base string, repo string, dir string) error {
	cmd := exec.Command("git", "clone", repo, dir)
	cmd.Dir = base
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func FormatPatch(repo string, output string) error {
	cmd := exec.Command("git", "format-patch", "--root", "-o", output)
	cmd.Dir = repo
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
