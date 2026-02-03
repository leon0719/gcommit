package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gcommit/internal"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	flagYes     bool
	flagEdit    bool
	flagGitMoji bool
	flagLang    string
)

var rootCmd = &cobra.Command{
	Use:   "gcommit",
	Short: "AI-powered git commit message generator",
	Long:  "Generate meaningful git commit messages using GPT",
	RunE:  runCommit,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "Commit directly without confirmation")
	rootCmd.Flags().BoolVarP(&flagEdit, "edit", "e", false, "Open editor after generating message")
	rootCmd.Flags().BoolVar(&flagGitMoji, "gitmoji", false, "Add gitmoji to commit message")
	rootCmd.Flags().StringVar(&flagLang, "lang", "", "Language for commit message (en, zh-TW, zh-CN, ja)")
}

func runCommit(cmd *cobra.Command, args []string) error {
	cfg, err := internal.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if flagGitMoji {
		cfg.GitMoji = true
	}
	if flagLang != "" {
		cfg.Language = flagLang
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Println(cyan("📝 Getting staged changes..."))

	diff, err := internal.GetStagedDiff()
	if err != nil {
		if err == internal.ErrNoStagedChanges {
			return fmt.Errorf("no staged changes. Use 'git add' first")
		}
		if err == internal.ErrNotGitRepo {
			return fmt.Errorf("not a git repository")
		}
		return err
	}

	files, _ := internal.GetStagedFiles()
	if len(files) > 0 {
		fmt.Printf("%s Staged files: %s\n", cyan("📁"), strings.Join(files, ", "))
	}

	fmt.Println(cyan("🤖 Generating commit message..."))

	message, err := internal.GenerateCommitMessage(cfg, diff)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(green("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(message)
	fmt.Println(green("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()

	if flagYes {
		return doCommit(message)
	}

	fmt.Printf("%s [Y]es / [n]o / [r]etry / [e]dit: ", yellow("Commit?"))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "", "y", "yes":
		return doCommit(message)
	case "r", "retry":
		return runCommit(cmd, args)
	case "e", "edit":
		return commitWithEditor(message)
	default:
		fmt.Println("Commit cancelled.")
		return nil
	}
}

func doCommit(message string) error {
	green := color.New(color.FgGreen).SprintFunc()

	if err := internal.Commit(message); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	fmt.Println(green("✅ Committed successfully!"))
	return nil
}

func commitWithEditor(message string) error {
	tmpfile, err := os.CreateTemp("", "gcommit-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(message); err != nil {
		return err
	}
	tmpfile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := newEditorCmd(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	edited, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return err
	}

	return doCommit(strings.TrimSpace(string(edited)))
}
