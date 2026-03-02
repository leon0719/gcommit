package internal

import (
	"testing"
)

func TestGetBranchDiff_InvalidBase(t *testing.T) {
	// GetBranchDiff with a nonexistent branch should return an error
	_, err := GetBranchDiff("nonexistent-branch-abc123")
	if err == nil {
		t.Fatal("expected error for nonexistent branch, got nil")
	}
}
