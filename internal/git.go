package internal

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	ErrNoStagedChanges = errors.New("no staged changes found")
	ErrNotGitRepo      = errors.New("not a git repository")
)

func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func GetStagedDiff() (string, error) {
	if !IsGitRepo() {
		return "", ErrNotGitRepo
	}

	cmd := exec.Command("git", "diff", "--cached", "--no-color")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	diff := strings.TrimSpace(string(output))
	if diff == "" {
		return "", ErrNoStagedChanges
	}

	return diff, nil
}

func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.TrimSpace(string(output))
	if lines == "" {
		return nil, nil
	}

	return strings.Split(lines, "\n"), nil
}

func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}
