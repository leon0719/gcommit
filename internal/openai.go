package internal

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

const (
	maxDiffLength = 8000
	apiTimeout    = 60 * time.Second
)

type GenerateResult struct {
	Message          string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

func GenerateCommitMessage(cfg *Config, diff string) (*GenerateResult, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key not set. Run: gcommit config set api_key <your-key>")
	}

	if len(diff) > maxDiffLength {
		diff = diff[:maxDiffLength] + "\n... (truncated)"
	}

	client := openai.NewClient(cfg.APIKey)

	prompt := buildPrompt(cfg, diff)

	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: cfg.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: getSystemPrompt(cfg),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   500,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	return &GenerateResult{
		Message:          cleanMessage(resp.Choices[0].Message.Content),
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
		TotalTokens:      resp.Usage.TotalTokens,
	}, nil
}

func cleanMessage(msg string) string {
	// Remove markdown code blocks
	re := regexp.MustCompile("(?s)```[a-z]*\n?(.*?)```")
	if matches := re.FindStringSubmatch(msg); len(matches) > 1 {
		msg = matches[1]
	}

	// Remove leading/trailing backticks
	msg = strings.Trim(msg, "`")

	return strings.TrimSpace(msg)
}

func getSystemPrompt(cfg *Config) string {
	prompt := "You are a helpful assistant that generates git commit messages."

	if cfg.Conventional {
		prompt += `
Use Conventional Commits format with these types:
- feat: new feature
- fix: bug fix
- docs: documentation
- style: formatting (no code change)
- refactor: code restructuring
- test: adding tests
- chore: maintenance tasks`
	}

	if cfg.GitMoji {
		prompt += `
Add a relevant gitmoji at the start:
- ✨ feat (new feature)
- 🐛 fix (bug fix)
- 📝 docs (documentation)
- 🎨 style (formatting)
- ♻️ refactor (restructuring)
- ✅ test (adding tests)
- 🔧 chore (maintenance)`
	}

	langMap := map[string]string{
		"en":    "English",
		"zh-TW": "Traditional Chinese",
		"zh-CN": "Simplified Chinese",
		"ja":    "Japanese",
	}

	if lang, ok := langMap[cfg.Language]; ok && cfg.Language != "en" {
		prompt += fmt.Sprintf("\nWrite the commit message in %s.", lang)
	}

	prompt += fmt.Sprintf("\nKeep the first line under %d characters.", cfg.MaxLength)

	if cfg.Why {
		prompt += "\nAfter the commit message, add a blank line and then 'Why: ' followed by a brief explanation of why these changes were made."
	}

	prompt += "\nDo NOT wrap the message in markdown code blocks or backticks. Output plain text only."

	return prompt
}

func buildPrompt(cfg *Config, diff string) string {
	return fmt.Sprintf(`Generate a concise commit message for these changes:

%s

Rules:
1. First line: brief summary (imperative mood)
2. Optional body: explain what and why (not how)
3. Be specific but concise`, diff)
}
