package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"stash/internal/storage"
)

// State represents the current UI state
type State int

const (
	StateNoteInput State = iota
	StateTitleInput
	StateDone
	StateError
)

// Model represents the UI model
type Model struct {
	state       State
	textarea    textarea.Model
	titleInput  textinput.Model
	targetDir   string
	noteContent string
	title       string
	savedPath   string
	err         error
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

// NewModel creates a new UI model
func NewModel(targetDir string) Model {
	ta := textarea.New()
	ta.Placeholder = "Type or paste your note here..."
	ta.Focus()

	// Full screen
	ta.SetWidth(80)
	ta.SetHeight(20)

	// Title input
	ti := textinput.New()
	ti.Placeholder = "Enter title..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	return Model{
		state:      StateNoteInput,
		textarea:   ta,
		titleInput: ti,
		targetDir:  targetDir,
	}
}

// NewModelWithContent creates a new UI model with pre-filled content
func NewModelWithContent(targetDir, content string) Model {
	m := NewModel(targetDir)
	m.noteContent = content

	// Generate auto-title and set it as the initial value
	autoTitle := storage.GenerateAutoTitle(content)
	m.titleInput.SetValue(autoTitle)
	m.titleInput.CursorEnd()

	// Skip to title input state
	m.state = StateTitleInput

	return m
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch m.state {
	case StateNoteInput:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC:
				return m, tea.Quit

			case tea.KeyCtrlS:
				// Save note and move to title input
				m.noteContent = m.textarea.Value()
				if strings.TrimSpace(m.noteContent) == "" {
					m.err = fmt.Errorf("note cannot be empty")
					m.state = StateError
					return m, nil
				}

				// Generate auto-title suggestion
				autoTitle := storage.GenerateAutoTitle(m.noteContent)
				m.titleInput.SetValue(autoTitle)
				m.titleInput.CursorEnd()

				m.state = StateTitleInput
				return m, textinput.Blink
			}
		}

		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)

	case StateTitleInput:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC:
				return m, tea.Quit

			case tea.KeyEnter:
				// Save the note
				m.title = m.titleInput.Value()
				if strings.TrimSpace(m.title) == "" {
					m.title = storage.GenerateAutoTitle(m.noteContent)
				}

				path, err := storage.SaveNote(m.targetDir, m.noteContent, m.title)
				if err != nil {
					m.err = err
					m.state = StateError
					return m, nil
				}

				m.savedPath = path
				m.state = StateDone
				return m, tea.Quit

			case tea.KeyEsc:
				// Go back to note input
				m.state = StateNoteInput
				return m, textarea.Blink
			}
		}

		m.titleInput, cmd = m.titleInput.Update(msg)
		cmds = append(cmds, cmd)

	case StateError:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeyCtrlC {
				return m, tea.Quit
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m Model) View() string {
	switch m.state {
	case StateNoteInput:
		return m.noteInputView()

	case StateTitleInput:
		return m.titleInputView()

	case StateDone:
		return fmt.Sprintf("✓ Note saved to: %s\n", m.savedPath)

	case StateError:
		return errorStyle.Render(fmt.Sprintf("Error: %v\n", m.err)) + "\nPress Enter to exit."

	default:
		return ""
	}
}

func (m Model) noteInputView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Stash Note"))
	b.WriteString("\n\n")
	b.WriteString(m.textarea.View())
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("ctrl+s: save  •  ctrl+c: quit"))

	return b.String()
}

func (m Model) titleInputView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Enter Title"))
	b.WriteString("\n\n")
	b.WriteString(m.titleInput.View())
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("enter: save  •  esc: back  •  ctrl+c: quit"))

	return b.String()
}

// WasSuccessful returns true if the note was saved successfully
func (m Model) WasSuccessful() bool {
	return m.state == StateDone
}

// SavedPath returns the path where the note was saved
func (m Model) SavedPath() string {
	return m.savedPath
}

// Error returns any error that occurred
func (m Model) Error() error {
	return m.err
}
