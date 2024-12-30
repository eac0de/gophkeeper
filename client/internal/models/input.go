package models

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eac0de/gophkeeper/client/internal/client"
	"github.com/eac0de/gophkeeper/client/internal/schemes"
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
	listDocStyle      = lipgloss.NewStyle().Margin(0, 0)

	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 0)
	activeTabStyle   = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle      = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(0, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

type listItem struct {
	title, desc string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title }

type inputModel struct {
	tokens     schemes.Tokens
	Tabs       []string
	TabContent []string
	activeTab  int
	client     *client.APIClient
	list       list.Model
}

func InitialInputModel(tokens schemes.Tokens, client *client.APIClient) inputModel {
	items := []list.Item{
		listItem{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		listItem{title: "Nutella", desc: "It's good on toast"},
		listItem{title: "Bitter melon", desc: "It cools you down"},
		listItem{title: "Nice socks", desc: "And by that I mean socks without holes"},
		listItem{title: "Eight hours of sleep", desc: "I had this once"},
		listItem{title: "Cats", desc: "Usually"},
		listItem{title: "Plantasia, the album", desc: "My plants love it too"},
		listItem{title: "Pour over coffee", desc: "It takes forever to make though"},
		listItem{title: "VR", desc: "Virtual reality...what is there to say?"},
		listItem{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
		listItem{title: "Linux", desc: "Pretty much the best OS"},
		listItem{title: "Business school", desc: "Just kidding"},
		listItem{title: "Pottery", desc: "Wet clay is a great feeling"},
		listItem{title: "Shampoo", desc: "Nothing like clean hair"},
		listItem{title: "Table tennis", desc: "It’s surprisingly exhausting"},
		listItem{title: "Milk crates", desc: "Great for packing in your extra stuff"},
		listItem{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
		listItem{title: "Stickers", desc: "The thicker the vinyl the better"},
		listItem{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
		listItem{title: "Warm light", desc: "Like around 2700 Kelvin"},
		listItem{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		listItem{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		listItem{title: "Terrycloth", desc: "In other words, towel fabric"},
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.SetShowTitle(false)
	list.SetShowPagination(false)
	list.SetShowStatusBar(false)
	list.KeyMap = NewKeyMap()
	return inputModel{
		Tabs:       []string{"      Texts      ", "      Files      ", "      BankCards      ", "      AuthInfos      "},
		TabContent: []string{"", "", "", ""},
		tokens:     tokens,
		client:     client,
		list:       list,
	}
}

func (m inputModel) Init() tea.Cmd {
	return nil
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			SaveTokens(m.tokens.AccessToken, m.tokens.RefreshToken)
			return m, tea.Quit
		case "right":
			m.activeTab = min(m.activeTab+1, 3)
			return m, nil
		case "left":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := listDocStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-10)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	m.TabContent[0] = m.list.View()
	sb := strings.Builder{}
	renderedTabs := m.renderTabs()
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	sb.WriteString(row)
	sb.WriteString("\n")
	sb.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TabContent[m.activeTab]))
	return tabDocStyle.Render(sb.String())
}

func (m inputModel) renderTabs() []string {
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

func NewKeyMap() list.KeyMap {
	return list.KeyMap{
		// Browsing.
		CursorUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		ClearFilter: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear-filters"),
		),

		// Filtering.
		CancelWhileFiltering: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "close"),
		),
		AcceptWhileFiltering: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "accept-filters"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "exit"),
		),
	}
}
