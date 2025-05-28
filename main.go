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
	bgColorIdx     int
	bgCategoryIdx  int
	fgColorIdx     int
	fgCategoryIdx  int
	bold           bool
	italic         bool
	faint          bool
	blink          bool
	strikethrough  bool
	underline      bool
	reverse        bool
}

func (m model) Init() tea.Cmd {
	// Initialize both background and foreground to the starting selection
	m.bgColorIdx = m.hue
	m.bgCategoryIdx = m.category
	m.fgColorIdx = m.hue
	m.fgCategoryIdx = m.category
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
		case "left", "shift+left":
			step := 1
			if msg.String() == "shift+left" {
				step = 2
			}
			m.hue -= step
			if m.hue < 0 {
				m.hue = 0
			}
		case "right", "shift+right":
			step := 1
			if msg.String() == "shift+right" {
				step = 2
			}
			m.hue += step
			if m.hue >= len(colorRanges[keys[m.category]]) {
				m.hue = len(colorRanges[keys[m.category]]) - 1
			}
		case "tab":
			// Persist the current selection to the appropriate color index/category
			if m.editForeground {
				m.fgColorIdx = m.hue
				m.fgCategoryIdx = m.category
			} else {
				m.bgColorIdx = m.hue
				m.bgCategoryIdx = m.category
			}
			m.editForeground = !m.editForeground
			// When switching, restore the previous selection for the new mode
			if m.editForeground {
				m.hue = m.fgColorIdx
				m.category = m.fgCategoryIdx
			} else {
				m.hue = m.bgColorIdx
				m.category = m.bgCategoryIdx
			}
		case "i":
			m.italic = !m.italic
		case "b":
			m.bold = !m.bold
		case "f":
			m.faint = !m.faint
		case "k":
			m.blink = !m.blink
		case "s":
			m.strikethrough = !m.strikethrough
		case "u":
			m.underline = !m.underline
		case "r":
			m.reverse = !m.reverse
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	keys := m.keys

	// Get the selected background and foreground color values
	bgCat := m.bgCategoryIdx
	bgIdx := m.bgColorIdx
	fgCat := m.fgCategoryIdx
	fgIdx := m.fgColorIdx

	// If currently editing, use the current selection for the active mode
	if m.editForeground {
		fgCat = m.category
		fgIdx = m.hue
	} else {
		bgCat = m.category
		bgIdx = m.hue
	}

	bgColor := colorRanges[keys[bgCat]][bgIdx]
	fgColor := colorRanges[keys[fgCat]][fgIdx]

	// Render the color preview bars and menu
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
		colors := colorRanges[cat]
		selectedIdx := 0
		if m.editForeground && i == fgCat {
			selectedIdx = fgIdx
		} else if !m.editForeground && i == bgCat {
			selectedIdx = bgIdx
		}
		maxWidth := 24
		start := 0
		end := len(colors)
		ellipsisLeft := false
		ellipsisRight := false
		if len(colors) > maxWidth {
			// Center the window on the selected color if possible
			start = selectedIdx - maxWidth/2
			if start < 0 {
				start = 0
			}
			end = start + maxWidth
			if end > len(colors) {
				end = len(colors)
				start = end - maxWidth
			}
			ellipsisLeft = start > 0
			ellipsisRight = end < len(colors)
		}
		if ellipsisLeft {
			colorPreview += "…"
		}
		for j := start; j < end; j++ {
			color := colors[j]
			marker := " "
			// Show marker for the currently selected color in the active mode
			if ((m.editForeground && i == fgCat && j == fgIdx) ||
				(!m.editForeground && i == bgCat && j == bgIdx)) &&
				m.keys[m.category] == cat {
				marker = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"}).
					Render("X")
			} else if i == fgCat && j == fgIdx && !m.editForeground {
				marker = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"}).
					Render("F")
			} else if i == bgCat && j == bgIdx && m.editForeground == true {
				marker = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"}).
					Render("B")
			}
			colorPreview += lipgloss.NewStyle().
				Background(lipgloss.Color(fmt.Sprintf("%d", color))).
				Render(marker)
		}
		if ellipsisRight {
			colorPreview += "…"
		}
		padding := longestCategory - len(cat) + 1
		menu += pointer + categoryStyle.Render(cat) + fmt.Sprintf("%*s", padding, "") + colorPreview + "\n"
	}

	selectedColorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(fmt.Sprintf("%d", bgColor))).
		Foreground(lipgloss.Color(fmt.Sprintf("%d", fgColor))).
		Padding(1, 4)
	if m.bold {
		selectedColorStyle = selectedColorStyle.Bold(true)
	}
	if m.italic {
		selectedColorStyle = selectedColorStyle.Italic(true)
	}
	if m.faint {
		selectedColorStyle = selectedColorStyle.Faint(true)
	}
	if m.blink {
		selectedColorStyle = selectedColorStyle.Blink(true)
	}
	if m.strikethrough {
		selectedColorStyle = selectedColorStyle.Strikethrough(true)
	}
	if m.underline {
		selectedColorStyle = selectedColorStyle.Underline(true)
	}
	if m.reverse {
		selectedColorStyle = selectedColorStyle.Reverse(true)
	}

	selectedColorExample := selectedColorStyle.Render(" Selected Color Example ")

	// Compose style indicators, always show all, faint if not enabled, highlight activation letter
	indicatorActive := lipgloss.NewStyle().Bold(true)
	indicatorInactive := lipgloss.NewStyle()
	indicatorFaint := indicatorInactive.Faint(true)
	letterActive := lipgloss.NewStyle().Foreground(lipgloss.Color("33")).Bold(true)
	letterInactive := lipgloss.NewStyle().Foreground(lipgloss.Color("33"))

	boldIndicator := indicatorInactive.Render(" [") + letterInactive.Render("B") + indicatorInactive.Render("old]")
	italicIndicator := indicatorInactive.Render(" [") + letterInactive.Render("I") + indicatorInactive.Render("talic]")
	faintIndicator := indicatorInactive.Render(" [") + letterInactive.Render("F") + indicatorInactive.Render("aint]")
	blinkIndicator := indicatorInactive.Render(" [blin") + letterInactive.Render("k") + indicatorInactive.Render("]")
	strikethroughIndicator := indicatorInactive.Render(" [") + letterInactive.Render("S") + indicatorInactive.Render("trikethrough]")
	underlineIndicator := indicatorInactive.Render(" [") + letterInactive.Render("U") + indicatorInactive.Render("nderline]")
	reverseIndicator := indicatorInactive.Render(" [") + letterInactive.Render("R") + indicatorInactive.Render("everse]")

	if m.bold {
		boldIndicator = indicatorActive.Render(" [") + letterActive.Render("B") + indicatorActive.Render("old]")
	} else {
		boldIndicator = indicatorFaint.Render(" [") + letterInactive.Render("B") + indicatorFaint.Render("old]")
	}
	if m.italic {
		italicIndicator = indicatorActive.Render(" [") + letterActive.Render("I") + indicatorActive.Render("talic]")
	} else {
		italicIndicator = indicatorFaint.Render(" [") + letterInactive.Render("I") + indicatorFaint.Render("talic]")
	}
	if m.faint {
		faintIndicator = indicatorActive.Render(" [") + letterActive.Render("F") + indicatorActive.Render("aint]")
	} else {
		faintIndicator = indicatorFaint.Render(" [") + letterInactive.Render("F") + indicatorFaint.Render("aint]")
	}
	if m.blink {
		blinkIndicator = indicatorActive.Render(" [blin") + letterActive.Render("k") + indicatorActive.Render("]")
	} else {
		blinkIndicator = indicatorFaint.Render(" [blin") + letterInactive.Render("k") + indicatorFaint.Render("]")
	}
	if m.strikethrough {
		strikethroughIndicator = indicatorActive.Render(" [") + letterActive.Render("S") + indicatorActive.Render("trikethrough]")
	} else {
		strikethroughIndicator = indicatorFaint.Render(" [") + letterInactive.Render("S") + indicatorFaint.Render("trikethrough]")
	}
	if m.underline {
		underlineIndicator = indicatorActive.Render(" [") + letterActive.Render("U") + indicatorActive.Render("nderline]")
	} else {
		underlineIndicator = indicatorFaint.Render(" [") + letterInactive.Render("U") + indicatorFaint.Render("nderline]")
	}
	if m.reverse {
		reverseIndicator = indicatorActive.Render(" [") + letterActive.Render("R") + indicatorActive.Render("everse]")
	} else {
		reverseIndicator = indicatorFaint.Render(" [") + letterInactive.Render("R") + indicatorFaint.Render("everse]")
	}

	ansiValue := fmt.Sprintf(" BG: %d | FG: %d", bgColor, fgColor)

	mode := "Editing: Background"
	if m.editForeground {
		mode = "Editing: Foreground"
	}
	// Stack preview info and style indicators vertically to the right of the preview
	styleIndicators0 := italicIndicator + boldIndicator + faintIndicator + blinkIndicator
	styleIndicators1 := strikethroughIndicator + underlineIndicator + reverseIndicator
	previewBlock := lipgloss.JoinHorizontal(
		lipgloss.Top,
		selectedColorExample,
		ansiValue+"\n"+styleIndicators0+"\n"+styleIndicators1,
	)
	return previewBlock + "\n" + mode + "\n\n" + menu + "\n\n" + highlightCommandsInstruction()
}

func highlightCommandsInstruction() string {
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("33")).Bold(true)
	plain := func(s string) string { return lipgloss.NewStyle().Render(s) }
	return plain("Commands: ") +
		cmdStyle.Render("↑/↓") + plain(" to change category, ") +
		cmdStyle.Render("←/→") + plain(" to change color, ") +
		cmdStyle.Render("Tab") + plain(" to switch mode, ") +
		cmdStyle.Render("q") + plain(" to quit")
}

var allColors = func() []int {
	colors := make([]int, 256)
	for i := 0; i < 256; i++ {
		colors[i] = i
	}
	return colors
}()

var colorRanges = map[string][]int{
	"Grayscale":    {232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255},
	"Basic Colors": {0, 1, 2, 3, 4, 5, 6, 7},
	"Dark Colors":  {8, 9, 10, 11, 12, 13, 14, 15},
	"Red Scale":    {52, 88, 124, 160, 196, 202, 208, 214, 220, 226},
	"Blue Scale":   {17, 18, 19, 20, 21, 27, 33, 39, 45, 51},
	"Green Scale":  {22, 28, 34, 40, 46, 82, 118, 154, 190, 226},
	"Yellow Scale": {226, 220, 214, 208, 202, 196, 190, 184, 178, 172},
	"Purple Scale": {55, 56, 57, 93, 129, 165, 201, 200, 199, 198},
	"All Colors":   allColors,
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
			// Sort keys alphabetically for consistent order
			for i := 0; i < len(keys)-1; i++ {
				for j := i + 1; j < len(keys); j++ {
					if keys[i] > keys[j] {
						keys[i], keys[j] = keys[j], keys[i]
					}
				}
			}
			return keys
		}(),
	}).Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
