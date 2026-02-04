package cmd

import (
	"fmt"

	"gcommit/internal"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE:  runConfigSet,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigGet,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	RunE:  runConfigList,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key, value := args[0], args[1]

	cfg, err := internal.LoadConfig()
	if err != nil {
		return err
	}

	switch key {
	case "api_key":
		cfg.APIKey = value
	case "model":
		cfg.Model = value
	case "language", "lang":
		cfg.Language = value
	case "conventional":
		cfg.Conventional = value == "true"
	case "gitmoji":
		cfg.GitMoji = value == "true"
	case "why":
		cfg.Why = value == "true"
	case "max_length":
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
			return fmt.Errorf("invalid number: %s", value)
		}
		cfg.MaxLength = n
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	if err := internal.SaveConfig(cfg); err != nil {
		return err
	}

	fmt.Printf("✅ Set %s = %s\n", key, maskIfSecret(key, value))
	return nil
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	cfg, err := internal.LoadConfig()
	if err != nil {
		return err
	}

	var value string
	switch key {
	case "api_key":
		value = maskIfSecret(key, cfg.APIKey)
	case "model":
		value = cfg.Model
	case "language", "lang":
		value = cfg.Language
	case "conventional":
		value = fmt.Sprintf("%v", cfg.Conventional)
	case "gitmoji":
		value = fmt.Sprintf("%v", cfg.GitMoji)
	case "why":
		value = fmt.Sprintf("%v", cfg.Why)
	case "max_length":
		value = fmt.Sprintf("%d", cfg.MaxLength)
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	fmt.Println(value)
	return nil
}

func runConfigList(cmd *cobra.Command, args []string) error {
	cfg, err := internal.LoadConfig()
	if err != nil {
		return err
	}

	display := struct {
		APIKey       string `yaml:"api_key"`
		Model        string `yaml:"model"`
		Conventional bool   `yaml:"conventional"`
		GitMoji      bool   `yaml:"gitmoji"`
		Why          bool   `yaml:"why"`
		Language     string `yaml:"language"`
		MaxLength    int    `yaml:"max_length"`
	}{
		APIKey:       maskIfSecret("api_key", cfg.APIKey),
		Model:        cfg.Model,
		Conventional: cfg.Conventional,
		GitMoji:      cfg.GitMoji,
		Why:          cfg.Why,
		Language:     cfg.Language,
		MaxLength:    cfg.MaxLength,
	}

	data, err := yaml.Marshal(display)
	if err != nil {
		return err
	}

	fmt.Print(string(data))
	return nil
}

func maskIfSecret(key, value string) string {
	if key == "api_key" && len(value) > 8 {
		return value[:4] + "****" + value[len(value)-4:]
	}
	return value
}
