package helpers

import (
	"os/exec"
	"path/filepath"
	"strings"
)

type GitService struct{}

func (git GitService) Clone(repoURL string) (repoName string, output []byte, status bool) {
	parts := strings.Split(strings.TrimSuffix(repoURL, ".git"), "/")
	repoName = parts[len(parts)-1]
	targetDir := filepath.Join("storage/repos", repoName)

	// Clone the repository
	status = true
	cmd := exec.Command("git", "clone", repoURL, targetDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		status = false
		return
	}

	// Track all remote branches
	script := `
	cd "` + targetDir + `" && \
	for branch in $(git branch -r | grep -v '\->'); do \
	git branch --track "${branch#origin/}" "$branch" 2>/dev/null || true; \
	done
	`
	cmd = exec.Command("bash", "-c", script)
	output, err = cmd.CombinedOutput()
	if err != nil {
		status = false
		return
	}

	return
}

func (git GitService) Branches(repoPath string) (output []byte, branches []string, status bool) {
	// Use --git-dir to point to the bare/mirrored repo
	cmd := exec.Command("git", "-C", repoPath, "branch", "-a")
	output, err := cmd.Output()
	if err != nil {
		status = false
	} else {
		status = true
		rawBranches := strings.Split(string(output), "\n")
		for _, b := range rawBranches {
			b = strings.TrimSpace(b)
			if b != "" {
				// remove "* " from current branch
				if strings.HasPrefix(b, "* ") {
					b = strings.TrimPrefix(b, "* ")
				}
				branches = append(branches, b)
			}
		}
	}

	return
}

func (git GitService) Pull(repoPath string, branch string) (output []byte, status bool) {
	status = true
	// Checkout the branch first
	checkoutCmd := exec.Command("git", "-C", repoPath, "checkout", branch)
	output, err := checkoutCmd.CombinedOutput()
	if err != nil {
		status = false
		return
	}
	// Then pull the latest changes
	pullCmd := exec.Command("git", "-C", repoPath, "pull", "origin", branch)
	output, err = pullCmd.CombinedOutput()
	if err != nil {
		status = false
		return
	}
	return
}
