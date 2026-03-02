package cmd

import (
	"errors"
	"fmt"
	"strings"

	"gcommit/internal"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	flagBranch string
	reviewLang string
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "AI-powered code review for your changes",
	Long:  "Review staged changes or branch diff using AI to find bugs, security issues, and improvements",
	RunE:  runReview,
}

func init() {
	reviewCmd.Flags().StringVarP(&flagBranch, "branch", "b", "", "Compare against specified branch (e.g., main)")
	reviewCmd.Flags().StringVar(&reviewLang, "lang", "", "Language for review output (en, zh-TW, zh-CN, ja)")
	rootCmd.AddCommand(reviewCmd)
}

func runReview(cmd *cobra.Command, args []string) error {
	cfg, err := internal.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if reviewLang != "" {
		cfg.Language = reviewLang
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgCyan).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()

	var diff string
	var files []string

	if flagBranch != "" {
		fmt.Printf("%s Comparing against %s...\n", cyan("🔍"), flagBranch)

		diff, err = internal.GetBranchDiff(flagBranch)
		if err != nil {
			return err
		}
		files, _ = internal.GetBranchFiles(flagBranch)
	} else {
		fmt.Println(cyan("🔍 Reviewing staged changes..."))

		diff, err = internal.GetStagedDiff()
		if err != nil {
			if errors.Is(err, internal.ErrNoStagedChanges) {
				return fmt.Errorf("no staged changes. Use 'git add' first or use --branch flag")
			}
			if errors.Is(err, internal.ErrNotGitRepo) {
				return fmt.Errorf("not a git repository")
			}
			return err
		}
		files, _ = internal.GetStagedFiles()
	}

	if len(files) > 0 {
		fmt.Printf("%s Files: %s\n", cyan("📁"), strings.Join(files, ", "))
	}

	fmt.Printf("%s Analyzing with %s...\n", cyan("🤖"), cfg.Model)

	result, err := internal.GenerateReview(cfg, diff)
	if err != nil {
		return err
	}

	fmt.Println()

	if result.Review == nil {
		fmt.Println(green("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		fmt.Println(result.RawResponse)
		fmt.Println(green("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	} else {
		review := result.Review

		fmt.Println(green("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		fmt.Printf("%s\n", green("📋 Summary"))
		fmt.Println(review.Summary)
		fmt.Println(green("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))

		if len(review.Issues) == 0 {
			fmt.Println()
			fmt.Println(green("✅ No issues found. Code looks good!"))
		} else {
			grouped := internal.GroupIssuesByCategory(review.Issues)

			for _, category := range internal.CategoryDisplayOrder {
				issues, ok := grouped[category]
				if !ok {
					continue
				}

				fmt.Println()

				var categoryColor func(a ...any) string
				var icon string
				switch category {
				case "bug", "security":
					categoryColor = red
					icon = "🔴"
				case "performance", "style":
					categoryColor = yellow
					icon = "🟡"
				default:
					categoryColor = blue
					icon = "🔵"
				}

				label := capitalize(category)
				fmt.Printf("%s %s (%d)\n", icon, categoryColor(label), len(issues))

				for _, issue := range issues {
					location := issue.File
					if issue.Line > 0 {
						location = fmt.Sprintf("%s:%d", issue.File, issue.Line)
					}
					fmt.Printf("  %s\n", gray(location))
					fmt.Printf("  %s\n", issue.Description)
					fmt.Println()
				}
			}
		}
	}

	fmt.Println(green("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Printf("%s\n", gray(fmt.Sprintf(
		"Tokens: %d (prompt: %d, completion: %d)",
		result.TotalTokens, result.PromptTokens, result.CompletionTokens,
	)))

	return nil
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
