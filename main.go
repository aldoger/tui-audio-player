package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aldoger/tui-audio-player/internal/logger"
	tea "github.com/charmbracelet/bubbletea"
)

func dirExist(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func listMusic(path string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Supported audio extensions
	audioExt := map[string]bool{
		".mp3":  true,
		".wav":  true,
		".flac": true,
		".aac":  true,
		".ogg":  true,
		".m4a":  true,
	}

	var musicFiles []os.DirEntry

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Name()))
		if audioExt[ext] {
			musicFiles = append(musicFiles, file)
		}
	}

	if len(musicFiles) < 1 {
		return nil, errors.New("empty directory, no audio files found")
	}

	return musicFiles, nil
}

type Model struct {
	choices  []os.DirEntry
	cursor   int
	selected map[int]struct{}
}

func initialModel(musicFiles []os.DirEntry) Model {
	return Model{
		choices: musicFiles,

		selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m Model) View() string {
	// The header
	s := "What music should we play?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Name())
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func main() {

	mylog := logger.Logger{}

	Home, err := os.UserHomeDir()
	if err != nil {
		mylog.Error(err.Error())
		os.Exit(1)
	}

	MusicDir := Home + "/Music"

	result, err := dirExist(MusicDir)
	if err != nil {
		mylog.Error(err.Error())
		os.Exit(1)
	}

	if !result {
		mylog.Error("directory %s does not exist", MusicDir)
		os.Exit(1)
	}

	musicFiles, err := listMusic(MusicDir)
	if err != nil {
		mylog.Error(err.Error())
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(musicFiles))
	if _, err := p.Run(); err != nil {
		mylog.Error("Error: %s", err.Error())
		os.Exit(1)
	}

}
