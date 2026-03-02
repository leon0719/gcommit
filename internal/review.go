package internal

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type ReviewIssue struct {
	Category    string `json:"category"`
	Severity    string `json:"severity"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Description string `json:"description"`
}

type ReviewResult struct {
	Summary string        `json:"summary"`
	Issues  []ReviewIssue `json:"issues"`
}

func ParseReviewResponse(raw string) (*ReviewResult, error) {
	cleaned := cleanJSONResponse(raw)

	var result ReviewResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("failed to parse review JSON: %w", err)
	}

	return &result, nil
}

func cleanJSONResponse(raw string) string {
	re := regexp.MustCompile("(?s)```(?:json)?\\s*\n?(.*?)```")
	if matches := re.FindStringSubmatch(raw); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return strings.TrimSpace(raw)
}

func GroupIssuesByCategory(issues []ReviewIssue) map[string][]ReviewIssue {
	grouped := make(map[string][]ReviewIssue)
	for _, issue := range issues {
		grouped[issue.Category] = append(grouped[issue.Category], issue)
	}
	return grouped
}

var CategoryDisplayOrder = []string{"bug", "security", "performance", "style", "suggestion"}

func GetReviewSystemPrompt(cfg *Config) string {
	langMap := map[string]string{
		"en":    "English",
		"zh-TW": "Traditional Chinese",
		"zh-CN": "Simplified Chinese",
		"ja":    "Japanese",
	}

	lang := "English"
	if l, ok := langMap[cfg.Language]; ok {
		lang = l
	}

	return fmt.Sprintf(`You are an expert code reviewer. Analyze the given git diff and provide a structured code review.

Respond ONLY with valid JSON in this exact format (no markdown code blocks, no extra text):
{
  "summary": "Brief overall assessment in 1-2 sentences",
  "issues": [
    {
      "category": "bug|security|performance|style|suggestion",
      "severity": "high|medium|low",
      "file": "path/to/file",
      "line": <line_number_or_0_if_unknown>,
      "description": "Clear description of the issue and how to fix it"
    }
  ]
}

Categories:
- bug: Logic errors, potential crashes, nil pointer issues, race conditions
- security: Injection, data leaks, authentication flaws, hardcoded secrets
- performance: Inefficient algorithms, unnecessary allocations, N+1 queries
- style: Naming, readability, code organization, missing error handling
- suggestion: Refactoring opportunities, better patterns, improvements

Rules:
- Only report genuine issues, not nitpicks
- If code looks good, return an empty issues array
- Be specific: reference exact file names and line numbers from the diff
- Write descriptions in %s`, lang)
}

func BuildReviewPrompt(diff string) string {
	return fmt.Sprintf("Review the following code changes:\n\n%s", diff)
}
