package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"stash/cmd"
	"stash/internal/storage"
	"stash/internal/ui"
)

func main() {
	// Check for subcommands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "config":
			if err := cmd.RunConfig(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return

		case "help", "--help", "-h":
			showHelp()
			return

		case "version", "--version", "-v":
			fmt.Println("stash v0.1.0")
			return
		}
	}

	// For stash command, we need config
	if err := cmd.CheckConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	targetDir, err := cmd.GetTargetDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Detect input source
	inputSource := detectInputSource()

	var m tea.Model

	switch inputSource {
	case "pipe":
		// Read from stdin
		content, err := readStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}

		if strings.TrimSpace(content) == "" {
			fmt.Fprintln(os.Stderr, "Error: note cannot be empty")
			os.Exit(1)
		}

		// Save directly without UI
		title := storage.GenerateAutoTitle(content)
		path, err := storage.SaveNote(targetDir, content, title)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving note: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Note saved to: %s\n", path)
		return

	case "args":
		// Read from command arguments
		content := strings.Join(os.Args[1:], " ")

		if strings.TrimSpace(content) == "" {
			fmt.Fprintln(os.Stderr, "Error: note cannot be empty")
			os.Exit(1)
		}

		// Save directly without UI (for non-interactive environments)
		title := storage.GenerateAutoTitle(content)
		path, err := storage.SaveNote(targetDir, content, title)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving note: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Note saved to: %s\n", path)
		return

	case "interactive":
		// Start UI for interactive input
		m = ui.NewModel(targetDir)
	}

	// Run the UI
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check if successful
	if um, ok := finalModel.(ui.Model); ok {
		if err := um.Error(); err != nil {
			// Error was already displayed in the UI
			os.Exit(1)
		}
	}
}

// detectInputSource determines where the input is coming from
func detectInputSource() string {
	// Check if stdin is a pipe
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return "pipe"
	}

	// Check if arguments were provided
	if len(os.Args) > 1 {
		return "args"
	}

	return "interactive"
}

// readStdin reads all content from stdin
func readStdin() (string, error) {
	var content strings.Builder
	buf := make([]byte, 1024)

	for {
		n, err := os.Stdin.Read(buf)
		if n > 0 {
			content.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	return content.String(), nil
}

func showHelp() {
	fmt.Println("stash - A simple CLI note-taking tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  stash                    Launch interactive mode")
	fmt.Println("  stash <note>             Save note from arguments")
	fmt.Println("  stash config --dir <path>  Configure target directory")
	fmt.Println("  echo <note> | stash      Save note from stdin")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  stash \"Meeting notes from today\"")
	fmt.Println("  cat notes.txt | stash")
	fmt.Println("  stash config --dir ~/Documents/notes")
	fmt.Println()
	fmt.Println("Interactive mode:")
	fmt.Println("  • Type your note in the full-screen editor")
	fmt.Println("  • Press Ctrl+S to save and enter a title")
	fmt.Println("  • Press Ctrl+C to quit")
}
