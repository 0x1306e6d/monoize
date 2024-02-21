package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Init(base string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = base
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

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

var layout = "Mon, 2 Jan 2006 15:04:05 -0700"

func ParsePatch(p string) (Patch, error) {
	s := bufio.NewScanner(strings.NewReader(p))
	if ok := s.Scan(); !ok {
		return Patch{}, fmt.Errorf("wrong patch")
	}
	if ok := s.Scan(); !ok {
		return Patch{}, fmt.Errorf("wrong patch")
	}
	from := s.Text()
	from = strings.TrimPrefix(from, "From: ")
	if ok := s.Scan(); !ok {
		return Patch{}, fmt.Errorf("wrong patch")
	}
	d := s.Text()
	d = strings.TrimPrefix(d, "Date: ")
	date, err := time.Parse(layout, d)
	if err != nil {
		return Patch{}, err
	}
	if ok := s.Scan(); !ok {
		return Patch{}, fmt.Errorf("wrong patch")
	}
	subject := s.Text()
	subject = strings.TrimPrefix(subject, "Subject: ")
	return Patch{From: from, Date: date, Subject: subject}, nil
}

type Patch struct {
	From    string
	Date    time.Time
	Subject string
}
