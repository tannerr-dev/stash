# Stash

![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-blue)
![License](https://img.shields.io/badge/license-MIT-green)

Stash lets you quickly save notes from your terminal. Whether you want to jot down a quick thought, pipe content from another command, or use the interactive editor, Stash makes it effortless to capture your thoughts.
Command:
```
stash
```
<img width="615" height="528" alt="1_tmp_screenshot" src="https://github.com/user-attachments/assets/42a4abe9-6bc9-4dda-a738-974f8afb44f9" />

## Features

- **Multiple Input Methods**: Type interactively, pass arguments, or pipe content from other commands
- **Smart Auto-Titles**: Automatically generates titles from your note's first line
- **Duplicate Protection**: Never overwrites existing notes (appends `-1`, `-2`, etc.)
- **Date-Prefixed Filenames**: All notes saved as `YYYY-MM-DD-sanitized-title.md`
- **Clean Markdown Output**: Every note is a `.md` file ready for your knowledge base
- **Full-Screen Editor**: Interactive mode uses a beautiful TUI with vim-style navigation
- **Simple Configuration**: One command to set your target directory

## Installation

### Prerequisites

- Go 1.18 or later

### Install with `go install`

```bash
go install github.com/tannerr-dev/stash@latest
```

Make sure `~/go/bin` is in your PATH:

```bash
export PATH="$PATH:$HOME/go/bin"
```

### Build from Source

```bash
git clone https://github.com/yourusername/stash.git
cd stash
go build -o stash .

# Optional: move to PATH
sudo mv stash /usr/local/bin/
```

## Quick Start

1. **Configure your notes directory:**

```bash
stash config --dir ~/notes
```

2. **Create your first note:**

```bash
stash "My first note using stash!"
```

## Usage

Stash supports three input methods:

### 1. Interactive Mode

Launch the full-screen editor:

```bash
stash
```

**Controls:**
- `Enter` or `Ctrl+S` - Save immediately with auto-generated title
- `Ctrl+R` - Customize title before saving (title editor opens with auto-generated title pre-filled)
- `Ctrl+C` - Quit

The auto-generated title is based on your note's first 32 characters.

### 2. Command Arguments

Pass your note directly:

```bash
stash "Meeting notes from today's standup"
stash "Idea: build a CLI tool in Go"
```

The note saves immediately with an auto-generated title.

### 3. Piped Input

Pipe content from other commands:

```bash
cat todo.txt | stash
echo "Quick thought" | stash
./my-script.sh | stash
```

## Configuration

Set your target directory where all notes will be saved:

```bash
stash config --dir ~/Documents/notes
```

Supports `~` for home directory:

```bash
stash config --dir ~/notes
```

Configuration is stored at `~/.config/stash/config.json`.

### View Current Config

```bash
cat ~/.config/stash/config.json
```

## File Naming

Notes are saved with the format:

```
YYYY-MM-DD-sanitized-title.md
```

**Examples:**
- `Meeting notes from today` → `2026-04-20-meeting-notes-from-today.md`
- `Idea: CLI tool!` → `2026-04-20-idea-cli-tool.md`

If a file with the same name exists, Stash appends a number:

```
2026-04-20-meeting-notes-from-today.md
2026-04-20-meeting-notes-from-today-1.md
2026-04-20-meeting-notes-from-today-2.md
```

## Auto-Title Generation

When using arguments or piped input, Stash automatically generates a title from the first line:

- Takes the first 32 characters
- Appends `...` if truncated
- Sanitizes for filesystem safety

**Example:**
```
"This is a very long note that exceeds thirty-two characters now"
→ "this-is-a-very-long-note-that-exceeds-thirt..."
```

## Examples

### Daily Journal Entry

```bash
stash "$(date '+%Y-%m-%d') Journal Entry

Today I learned about Go interfaces..."
```

### Save Code Snippet

```bash
grep -A 5 "func main" main.go | stash
```

### Meeting Notes

```bash
stash "Team Meeting - $(date)

Attendees: Alice, Bob, Charlie

Action items:
- [ ] Review PR #123
- [ ] Update documentation"
```

### From Clipboard

```bash
# macOS
pbpaste | stash

# Linux (with xclip)
xclip -o | stash

# Windows (with PowerShell)
Get-Clipboard | stash
```

## Interactive Mode Details

When running `stash` without arguments:

1. **Note Input**: Full-screen textarea for typing your note
2. **Quick Save**: Press `Enter` or `Ctrl+S` to save immediately with auto-generated title
3. **Customize Title** (Optional): Press `Ctrl+R` to edit the title before saving
4. **Confirmation**: Note saved, path displayed

**Navigation:**
- Standard text editing with cursor keys
- Paste with your terminal's paste command (Ctrl+Shift+V or Cmd+V)
- Vim-style keybindings supported

**Workflow:**
- Type your note → Press `Enter` → Saved with auto-generated title
- Type your note → Press `Ctrl+R` → Edit title → Press `Enter` → Saved with custom title

## Error Handling

Stash provides clear error messages:

```bash
# Missing config
$ stash
Stash is not configured yet.
Please run: stash config --dir <path>

# Missing directory
$ stash "test"
Error: Target directory does not exist: /path/to/notes
Please create the directory or update the config.

# Empty note
$ stash ""
Error: note cannot be empty
```

## Tips & Tricks

### Create Aliases

```bash
# Quick daily note
alias journal='stash "$(date +%Y-%m-%d)"'

# Note with timestamp
alias now='stash "$(date "+%Y-%m-%d %H:%M")"'
```

### Integrate with Git

```bash
# Commit notes automatically
cd ~/notes && git add . && git commit -m "Update notes $(date)"
```

### Search Notes

```bash
# Find all notes mentioning "meeting"
grep -r "meeting" ~/notes/

# List today's notes
ls ~/notes/$(date +%Y-%m-%d)*
```

## Development

### Project Structure

```
stash/
├── main.go              # Entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── storage/         # File operations
│   └── ui/              # Bubble Tea TUI
└── cmd/
    └── config.go        # Config subcommand
```

### Build

```bash
go build -o stash .
```

### Run Tests

```bash
go test ./...
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions welcome! Please feel free to submit a Pull Request.

## Acknowledgments

Built with ❤️ using the [Charm](https://charm.sh) ecosystem.
