package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	hue            int
	category       int
	keys           []string
	editForeground bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keys := m.keys
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.category = (m.category - 1 + len(keys)) % len(keys)
		case "down":
			m.category = (m.category + 1) % len(keys)
		case "left":
			m.hue = (m.hue - 1 + len(colorRanges[keys[m.category]])) % len(colorRanges[keys[m.category]])
		case "right":
			m.hue = (m.hue + 1) % len(colorRanges[keys[m.category]])
		case "tab":
			m.editForeground = !m.editForeground
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	colorBar := ""
	keys := m.keys
	currentRange := colorRanges[keys[m.category]]
	for i, color := range currentRange {
		pointer := " "
		if i == m.hue {
			pointer = "O"
		}
		colorBar += lipgloss.NewStyle().
			Background(lipgloss.Color(fmt.Sprintf("%d", color))).
			Render(" ")
		if i == m.hue {
			colorBar = colorBar[:len(colorBar)-1] + pointer
		}
	}

	categoryStyle := lipgloss.NewStyle().Bold(true)
	menu := ""
	longestCategory := 0
	for _, cat := range m.keys {
		if len(cat) > longestCategory {
			longestCategory = len(cat)
		}
	}

	for i, cat := range m.keys {
		pointer := "  "
		if i == m.category {
			pointer = "→ "
		}
		colorPreview := ""
		for i, color := range colorRanges[cat] {
			pointer := " "
			if i == m.hue && m.keys[m.category] == cat {
				pointer = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"}).
					Render("X")
			}
			colorPreview += lipgloss.NewStyle().
				Background(lipgloss.Color(fmt.Sprintf("%d", color))).
				Render(pointer)
		}
		padding := longestCategory - len(cat) + 1
		menu += pointer + categoryStyle.Render(cat) + fmt.Sprintf("%*s", padding, "") + colorPreview + "\n"
	}

	selectedColorValue := lipgloss.Color(fmt.Sprintf("%d", currentRange[m.hue]))
	staticBackground := lipgloss.Color("0")
	staticForeground := lipgloss.Color("15")

	if !m.editForeground {
		staticBackground = lipgloss.Color(fmt.Sprintf("%d", currentRange[m.hue]))
	}

	selectedBackground := staticBackground
	selectedForeground := staticForeground
	if m.editForeground {
		selectedForeground = lipgloss.Color(fmt.Sprintf("%d", currentRange[m.hue]))
	}

	if m.editForeground {
		selectedForeground = selectedColorValue
	} else {
		selectedBackground = selectedColorValue
	}

	selectedColorStyle := lipgloss.NewStyle().
		Background(selectedBackground).
		Foreground(selectedForeground).
		Padding(1, 4)

	selectedColorExample := selectedColorStyle.Render(" Selected Color Example ")

	ansiValue := fmt.Sprintf("ANSI Value: %d", currentRange[m.hue])

	mode := "Editing: Background"
	if m.editForeground {
		mode = "Editing: Foreground"
	}
	return selectedColorExample + "  " + ansiValue + "\n" + mode + "\n\n" + menu + "\n\n↑/↓ to change category, ←/→ to change color, Tab to switch mode, q to quit"
}

var colorRanges = map[string][]int{
	"Grayscale":     {232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255},
	"Basic Colors":  {0, 1, 2, 3, 4, 5, 6, 7},
	"Bright Colors": {8, 9, 10, 11, 12, 13, 14, 15},
	//	"256 Colors":    {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255},
	"Red Scale":    {52, 88, 124, 160, 196, 202, 208, 214, 220, 226},
	"Blue Scale":   {17, 18, 19, 20, 21, 27, 33, 39, 45, 51},
	"Green Scale":  {22, 28, 34, 40, 46, 82, 118, 154, 190, 226},
	"Yellow Scale": {226, 220, 214, 208, 202, 196, 190, 184, 178, 172},
	"Purple Scale": {55, 56, 57, 93, 129, 165, 201, 200, 199, 198},
	"Orange Scale": {202, 208, 214, 220, 226, 220, 214, 208, 202, 196},
}

func main() {
	if err := tea.NewProgram(model{
		hue:      0,
		category: 0,
		keys: func() []string {
			keys := make([]string, 0, len(colorRanges))
			for k := range colorRanges {
				keys = append(keys, k)
			}
			return keys
		}(),
	}).Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
