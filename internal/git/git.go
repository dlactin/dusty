package git

import (
	"bytes"
	"fmt"
	"math"
	"os/exec"
	"slices"
	"strings"
	"time"
)

var (
	protectedBranches = []string{"main", "master"}
)

type GitBranch struct {
	Remote string
	Name   string
	Author string
	Age    int
	Merged bool
}

// Use Git CLI to get local branches with name, author, date, and upstream
func GetLocalBranches() ([]GitBranch, error) {
	cmd := exec.Command("git", "for-each-ref",
		"--format=%(refname:short)|%(authorname)|%(authordate:iso)|%(upstream:short)",
		"refs/heads/")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error checking git branches: %w", err)
	}

	lines := bytes.Split(bytes.TrimSpace(out), []byte("\n"))
	var branches []GitBranch

	for _, line := range lines {
		parts := strings.Split(string(line), "|")
		if len(parts) < 3 {
			continue
		}

		name := parts[0]
		author := parts[1]
		dateStr := parts[2]
		remote := ""
		if len(parts) > 3 {
			remote = parts[3]
		}

		// Parse the date (we'll calculate age as Duration since commit)
		t, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
		if err != nil {
			continue
		}
		duration := time.Since(t).Hours()

		daysOld := int(math.Round(duration / 24))

		// Check if merged into HEAD
		merged := isMerged(name)

		// Exclude protected branches
		if !slices.Contains(protectedBranches, name) {
			branches = append(branches, GitBranch{
				Remote: remote,
				Name:   name,
				Author: author,
				Age:    daysOld,
				Merged: merged,
			})
		}
	}

	return branches, nil
}

// Helper to check if branch is merged into HEAD
func isMerged(branch string) bool {
	cmd := exec.Command("git", "branch", "--merged")
	out, _ := cmd.Output()
	return strings.Contains(string(out), branch)
}

// Helper to delete the local branch
// We use -d unless the force flag is used
func DelBranch(branch string, force bool) error {
	delFlag := "-d"

	if force {
		delFlag = "-D"
	}

	cmd := exec.Command("git", "branch", delFlag, branch)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("git", "branch", delFlag, branch)
		return fmt.Errorf("error deleting git branch: %w", err)
	}
	fmt.Printf("deleted %s branch\n", branch)
	return nil
}
