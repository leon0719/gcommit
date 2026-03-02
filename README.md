# gcommit

AI-powered git commit message generator and code reviewer using GPT.

## Installation

### From source

```bash
go install github.com/yourusername/gcommit@latest
```

### Manual build

```bash
git clone https://github.com/yourusername/gcommit.git
cd gcommit
go build -o gcommit .
sudo cp gcommit /usr/local/bin/
```

## Setup

```bash
# Set your OpenAI API key
gcommit config set api_key sk-xxxxxxxx

# Optional: change model (default: gpt-4o-mini)
gcommit config set model gpt-4o
```

## Usage

```bash
# Stage your changes
git add .

# Generate commit message
gcommit

# Options
gcommit -y           # Commit directly without confirmation
gcommit -e           # Open editor to modify message
gcommit -c           # Copy message to clipboard (no commit)
gcommit --gitmoji    # Add emoji to commit message
gcommit --why        # Include explanation of why changes were made
gcommit --lang zh-TW # Generate message in Traditional Chinese
```

### Code Review

```bash
# Review staged changes
gcommit review

# Review branch diff against main
gcommit review --branch main

# Review in Traditional Chinese
gcommit review --lang zh-TW
```

## Features

### Code Review

Reviews are categorized by severity:

```
🔍 Reviewing staged changes...
📁 Files: internal/handler.go, internal/db.go
🤖 Analyzing with gpt-4o-mini...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 Summary
Overall quality is good with 2 issues to address.
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔴 Bug (1)
  internal/handler.go:45
  nil pointer dereference - user may be nil

🔴 Security (1)
  internal/db.go:102
  SQL injection risk in raw query

🔵 Suggestion (1)
  cmd/root.go:23
  Consider extracting retry logic into a helper

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Tokens: 1,203 (prompt: 980, completion: 223)
```

Categories: Bug, Security, Performance, Style, Suggestion

### Token Usage Display

Each generation shows token consumption:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
feat(auth): add login endpoint
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Tokens: 245 (total: 245)
```

When using retry, total tokens accumulate:

```
Tokens: 238 (total: 483)
```

### Interactive Options

After generating a message, you can:
- `Y` - Commit with the generated message
- `n` - Cancel
- `r` - Retry (generate a new message)
- `e` - Edit the message before committing

## Configuration

Config file location: `~/.config/gcommit/config.yaml`

```bash
# View all settings
gcommit config list

# Set a value
gcommit config set <key> <value>

# Get a value
gcommit config get <key>
```

### Available options

| Key | Default | Description |
|-----|---------|-------------|
| `api_key` | - | OpenAI API key |
| `model` | `gpt-4o-mini` | GPT model to use |
| `conventional` | `true` | Use Conventional Commits format |
| `gitmoji` | `false` | Add gitmoji to messages |
| `why` | `false` | Include why explanation |
| `language` | `en` | Message language (en, zh-TW, zh-CN, ja) |
| `max_length` | `72` | Max length for commit title |

## License

MIT
