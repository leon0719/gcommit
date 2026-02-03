# gcommit

AI-powered git commit message generator using GPT.

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
gcommit --gitmoji    # Add emoji to commit message
gcommit --lang zh-TW # Generate message in Traditional Chinese
```

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
| `language` | `en` | Message language (en, zh-TW, zh-CN, ja) |
| `max_length` | `72` | Max length for commit title |

## License

MIT
