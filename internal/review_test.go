package internal

import (
	"testing"
)

func TestParseReviewResponse_ValidJSON(t *testing.T) {
	input := `{
		"summary": "Code looks good overall.",
		"issues": [
			{
				"category": "bug",
				"severity": "high",
				"file": "main.go",
				"line": 10,
				"description": "Nil pointer dereference"
			},
			{
				"category": "suggestion",
				"severity": "low",
				"file": "utils.go",
				"line": 25,
				"description": "Consider using a constant"
			}
		]
	}`

	result, err := ParseReviewResponse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Summary != "Code looks good overall." {
		t.Errorf("summary = %q, want %q", result.Summary, "Code looks good overall.")
	}
	if len(result.Issues) != 2 {
		t.Fatalf("issues count = %d, want 2", len(result.Issues))
	}
	if result.Issues[0].Category != "bug" {
		t.Errorf("issue[0].category = %q, want %q", result.Issues[0].Category, "bug")
	}
	if result.Issues[0].Severity != "high" {
		t.Errorf("issue[0].severity = %q, want %q", result.Issues[0].Severity, "high")
	}
}

func TestParseReviewResponse_WrappedInCodeBlock(t *testing.T) {
	input := "```json\n{\"summary\": \"All good.\", \"issues\": []}\n```"

	result, err := ParseReviewResponse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Summary != "All good." {
		t.Errorf("summary = %q, want %q", result.Summary, "All good.")
	}
}

func TestParseReviewResponse_InvalidJSON(t *testing.T) {
	input := "This is just plain text review."

	_, err := ParseReviewResponse(input)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestGroupIssuesByCategory(t *testing.T) {
	issues := []ReviewIssue{
		{Category: "bug", Severity: "high", File: "a.go", Line: 1, Description: "bug1"},
		{Category: "suggestion", Severity: "low", File: "b.go", Line: 2, Description: "sug1"},
		{Category: "bug", Severity: "medium", File: "c.go", Line: 3, Description: "bug2"},
	}

	grouped := GroupIssuesByCategory(issues)
	if len(grouped["bug"]) != 2 {
		t.Errorf("bug count = %d, want 2", len(grouped["bug"]))
	}
	if len(grouped["suggestion"]) != 1 {
		t.Errorf("suggestion count = %d, want 1", len(grouped["suggestion"]))
	}
}

func TestGetReviewSystemPrompt(t *testing.T) {
	cfg := &Config{Language: "zh-TW"}
	prompt := GetReviewSystemPrompt(cfg)

	if prompt == "" {
		t.Fatal("expected non-empty prompt")
	}
	if !containsSubstr(prompt, "JSON") {
		t.Error("prompt should mention JSON format")
	}
	if !containsSubstr(prompt, "Traditional Chinese") {
		t.Error("prompt should mention Traditional Chinese for zh-TW")
	}
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
