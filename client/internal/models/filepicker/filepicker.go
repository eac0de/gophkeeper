package filepicker

import (
	"net/http"
	"os"
	"strings"
	"time"

	fper "github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eac0de/gophkeeper/client/internal/client"
	"github.com/eac0de/gophkeeper/client/internal/models/item"
	"github.com/eac0de/gophkeeper/client/internal/models/login"
	"github.com/eac0de/gophkeeper/client/internal/schemes"
	"github.com/google/uuid"
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	tabDocStyle       = lipgloss.NewStyle().Padding(0, 0, 0, 0)

	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 0)
	activeTabStyle   = inactiveTabStyle.Border(activeTabBorder, true)
	errStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	windowStyle      = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(0, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	blurredStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	helpStyle        = blurredStyle
)

type model struct {
	Tabs      []string
	activeTab int

	fp           fper.Model
	selectedFile string
	quitting     bool
	errMsg       string
	invalidFile  string
	fileNotPick  bool

	client    *client.APIClient
	nextModel tea.Model
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func New(nextModel tea.Model, client *client.APIClient) model {
	fp := fper.New()
	fp.AllowedTypes = []string{".mod", ".sum", ".go", ".txt", ".md", ".png", ".jpg", ".jpeg"}
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp, cmd := fp.Update(fp.Init()())
	fp.Update(cmd)
	fp.AutoHeight = false
	fp.Height = 10
	m := model{
		fp:        fp,
		client:    client,
		nextModel: nextModel,

		Tabs:      []string{"      Texts      ", "      Files      ", "      BankCards      ", "      AuthInfos      "},
		activeTab: 1,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.fp.Height = msg.Height - 12
	case schemes.SaveFileDataMsg:
		if msg.Err != nil {
			if msg.StatusCode == http.StatusUnauthorized {
				return login.New(m.client, m.nextModel, msg.Err.Error()), nil
			}
			m.errMsg = msg.Err.Error()
			return m, nil
		}
		return item.New(m.client, msg.FileData, 1, m.nextModel, -1), nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.nextModel, nil
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "ctrl+s":
			if m.selectedFile == "" {
				m.fileNotPick = true
				return m, clearErrorAfter(1 * time.Second)
			}
			return m, func() tea.Msg {
				return m.client.SaveUserFileData(uuid.Nil, "", m.selectedFile, nil)
			}
		}
	case clearErrorMsg:
		m.fileNotPick = false
		m.errMsg = ""
	}

	var cmd tea.Cmd
	m.fp, cmd = m.fp.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.fp.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.invalidFile = ""
		m.selectedFile = path
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.fp.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.invalidFile = path
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(1*time.Second))
	}

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.errMsg != "" {
		s.WriteString(m.fp.Styles.DisabledFile.Render(m.errMsg))
	} else if m.selectedFile == "" {
		if m.fileNotPick {
			s.WriteString("No file selected")
		} else {
			s.WriteString("Pick a file:")
		}

	} else {
		s.WriteString("Selected file: " + m.fp.Styles.Selected.Render(m.selectedFile))
	}
	if m.invalidFile != "" {
		s.WriteString("\n" + m.fp.Styles.DisabledFile.Render(m.invalidFile) + " is not a valid file.")
	}
	s.WriteString("\n\n" + m.fp.View())

	sb := strings.Builder{}
	renderedTabs := m.renderTabs()
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	sb.WriteString(row)
	sb.WriteString("\n")
	sb.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(s.String() + helpStyle.Render("\n\n"+"↑/↓ up/down •  • ctrl+s save • ctrl+d delete • ctrl+c exit")))

	if m.errMsg != "" {
		sb.WriteString("\n\n" + errStyle.Render(m.errMsg))
	}
	return tabDocStyle.Render(sb.String())
}

func (m model) renderTabs() []string {
	var renderedTabs []string
	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}
	return renderedTabs
}
