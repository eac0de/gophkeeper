package home

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eac0de/gophkeeper/client/internal/client"
	"github.com/eac0de/gophkeeper/client/internal/components"
	"github.com/eac0de/gophkeeper/client/internal/models/filepicker"
	"github.com/eac0de/gophkeeper/client/internal/models/item"
	"github.com/eac0de/gophkeeper/client/internal/models/login"
	"github.com/eac0de/gophkeeper/client/internal/schemes"
	"github.com/eac0de/gophkeeper/client/internal/utils"
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
	tabDocStyle       = lipgloss.NewStyle()
	listDocStyle      = lipgloss.NewStyle()

	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 0)
	activeTabStyle   = inactiveTabStyle.Border(activeTabBorder, true)
	errStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	helpStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	windowStyle      = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(0, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

type listItem struct {
	title, desc string
	itemID      uuid.UUID
	itemScheme  interface{}
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title }

type homeModel struct {
	Tabs       []string
	TabContent []list.Model
	TabHelps   []string
	activeTab  int
	client     *client.APIClient
	errMsg     string

	height int
}

func New(client *client.APIClient) tea.Model {
	listModels := []list.Model{}
	for i := 0; i < 4; i++ {
		lst := list.New(nil, list.NewDefaultDelegate(), 0, 10)
		lst.SetShowTitle(false)
		lst.SetShowPagination(false)
		lst.SetShowStatusBar(false)
		lst.SetShowHelp(false)
		lst.SetFilteringEnabled(false)
		lst.KeyMap = NewKeyMap()
		listModels = append(listModels, lst)
	}
	return homeModel{
		Tabs: components.Tabs,
		TabHelps: []string{
			"↑/↓/→ • enter select • ctrl+n new • ctrl+d delete • ctrl+c exit",
			"↑/↓/→/← • enter select • ctrl+n new • ctrl+d delete • enter+tab download • ctrl+c exit",
			"↑/↓/→/← • enter select • ctrl+n new • ctrl+d delete • ctrl+c exit",
			"↑/↓/← • enter select • ctrl+n new • ctrl+d delete • ctrl+c exit",
		},
		TabContent: listModels,
		client:     client,
		height:     10,
	}
}

func (m homeModel) Init() tea.Cmd {
	return func() tea.Msg { return m.client.GetUserTextDataList("0") }
}

func (m homeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case schemes.DeleteListItemMsg:
		if msg.Err != nil {
			if msg.StatusCode == http.StatusUnauthorized {
				return login.New(m.client, m, msg.Err.Error()), nil
			}
			m.errMsg = msg.Err.Error()
			return m, nil
		}
		m.TabContent[m.activeTab].RemoveItem(msg.ItemIndex)
		return m, nil
	case schemes.UpdateListItemMsg:
		var item listItem
		if m.activeTab <= 0 {
			data, ok := msg.Item.(schemes.UserTextData)
			if !ok {
				return m, nil
			}
			item.title = data.Name
			item.desc = fmt.Sprintf("Created: %s, Updated: %s",
				data.CreatedAt.Format("15:04:05 02.01.2006"),
				data.UpdatedAt.Format("15:04:05 02.01.2006"))
			item.itemID = data.ID
			item.itemScheme = data
		} else if m.activeTab == 1 {
			data, ok := msg.Item.(schemes.UserFileData)
			if !ok {
				return m, nil
			}
			item.title = data.Name
			item.desc = fmt.Sprintf("Created: %s, Updated: %s",
				data.CreatedAt.Format("15:04:05 02.01.2006"),
				data.UpdatedAt.Format("15:04:05 02.01.2006"))
			item.itemID = data.ID
			item.itemScheme = data
		} else if m.activeTab == 2 {
			data, ok := msg.Item.(schemes.UserBankCard)
			if !ok {
				return m, nil
			}
			item.title = data.Name
			item.desc = fmt.Sprintf("Created: %s, Updated: %s",
				data.CreatedAt.Format("15:04:05 02.01.2006"),
				data.UpdatedAt.Format("15:04:05 02.01.2006"))
			item.itemID = data.ID
			item.itemScheme = data
		} else {
			data, ok := msg.Item.(schemes.UserAuthInfo)
			if !ok {
				return m, nil
			}
			item.title = data.Name
			item.desc = fmt.Sprintf("Created: %s, Updated: %s",
				data.CreatedAt.Format("15:04:05 02.01.2006"),
				data.UpdatedAt.Format("15:04:05 02.01.2006"))
			item.itemID = data.ID
			item.itemScheme = data
		}
		if msg.ItemIndex < 0 {
			m.TabContent[m.activeTab].InsertItem(-1, item)
		} else {
			m.TabContent[m.activeTab].SetItem(msg.ItemIndex, item)
		}
		return m, nil

	case schemes.SaveTextDataMsg:
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			utils.SaveTokens(m.client.Tokens)
			return m, tea.Quit
		case "right":
			m.activeTab = min(m.activeTab+1, 3)
			itemsLen := len(m.TabContent[m.activeTab].Items())
			if itemsLen < 2 {
				offset := strconv.Itoa(itemsLen)
				if m.activeTab <= 0 {
					return m, func() tea.Msg { return m.client.GetUserTextDataList(offset) }
				}
				if m.activeTab == 1 {
					return m, func() tea.Msg { return m.client.GetUserFileDataList(offset) }
				}
				if m.activeTab == 2 {
					return m, func() tea.Msg { return m.client.GetUserBankCardList(offset) }
				}
				if m.activeTab >= 3 {
					return m, func() tea.Msg { return m.client.GetUserAuthInfoList(offset) }
				}
			}
			return m, nil
		case "left":
			m.activeTab = max(m.activeTab-1, 0)
			itemsLen := len(m.TabContent[m.activeTab].Items())
			if itemsLen < 2 {
				offset := strconv.Itoa(itemsLen)
				if m.activeTab <= 0 {
					return m, func() tea.Msg { return m.client.GetUserTextDataList(offset) }
				}
				if m.activeTab == 1 {
					return m, func() tea.Msg { return m.client.GetUserFileDataList(offset) }
				}
				if m.activeTab == 2 {
					return m, func() tea.Msg { return m.client.GetUserBankCardList(offset) }
				}
				if m.activeTab >= 3 {
					return m, func() tea.Msg { return m.client.GetUserAuthInfoList(offset) }
				}
			}
			return m, nil
		case "ctrl+o":
			return login.New(m.client, New(m.client), ""), nil
		case "ctrl+n":
			if m.activeTab == 1 {
				return filepicker.New(m, m.client, m.height), nil
			}
			return item.New(m.client, nil, m.activeTab, m, -1), nil
		case "down":
			itemsLen := len(m.TabContent[m.activeTab].Items())
			var cmds tea.BatchMsg
			if itemsLen <= 0 {
				if m.activeTab <= 0 {
					cmds = append(cmds, func() tea.Msg { return m.client.GetUserTextDataList("0") })
				}
				if m.activeTab == 1 {
					cmds = append(cmds, func() tea.Msg { return m.client.GetUserFileDataList("0") })
				}
				if m.activeTab == 2 {
					cmds = append(cmds, func() tea.Msg { return m.client.GetUserBankCardList("0") })
				}
				if m.activeTab >= 3 {
					cmds = append(cmds, func() tea.Msg { return m.client.GetUserAuthInfoList("0") })
				}
			} else if m.TabContent[m.activeTab].Index() == itemsLen-2 {
				offset := strconv.Itoa(itemsLen)
				if m.activeTab <= 0 {
					cmds = append(cmds, func() tea.Msg { return m.client.GetUserTextDataList(offset) })
				}
				if m.activeTab == 1 {
					cmds = append(cmds, func() tea.Msg { return m.client.GetUserFileDataList(offset) })
				}
				if m.activeTab == 2 {
					cmds = append(cmds, func() tea.Msg { return m.client.GetUserBankCardList(offset) })
				}
				if m.activeTab >= 3 {
					cmds = append(cmds, func() tea.Msg { return m.client.GetUserAuthInfoList(offset) })
				}
			}
			var cmd tea.Cmd
			m.TabContent[m.activeTab], cmd = m.TabContent[m.activeTab].Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case "enter":
			i := m.TabContent[m.activeTab].SelectedItem()
			if i == nil {
				return m, nil
			}
			return item.New(m.client, i.(listItem).itemScheme, m.activeTab, m, m.TabContent[m.activeTab].Index()), nil
		case "ctrl+d":
			i := m.TabContent[m.activeTab].SelectedItem()
			if i == nil {
				return m, nil
			}
			listItem := i.(listItem)
			return m, func() tea.Msg {
				return m.client.DeleteData(listItem.itemID, m.TabContent[m.activeTab].Index(), m.activeTab)
			}
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		h, v := listDocStyle.GetFrameSize()
		for i := 0; i < len(m.TabContent); i++ {
			m.TabContent[i].SetSize(msg.Width-h, msg.Height-v-10)
		}

	case schemes.GetUserTextDataListMsg:
		if msg.StatusCode == http.StatusUnauthorized {
			return login.New(m.client, New(m.client), msg.Err.Error()), nil
		} else if msg.Err != nil {
			m.errMsg = msg.Err.Error()
			return m, nil
		} else {
			var items []list.Item = m.TabContent[0].Items()
			for _, item := range msg.List {
				items = append(items, listItem{
					title: item.Name,
					desc: fmt.Sprintf("Created: %s, Updated: %s",
						item.CreatedAt.Format("15:04:05 02.01.2006"),
						item.UpdatedAt.Format("15:04:05 02.01.2006")),
					itemScheme: item,
					itemID:     item.ID})
			}
			m.TabContent[0].SetItems(items)
			return m, nil
		}
	case schemes.GetUserFileDataListMsg:
		if msg.StatusCode == http.StatusUnauthorized {
			return login.New(m.client, New(m.client), msg.Err.Error()), nil
		} else if msg.Err != nil {
			m.errMsg = msg.Err.Error()
			return m, nil
		} else {
			var items []list.Item = m.TabContent[1].Items()
			for _, item := range msg.List {
				items = append(items, listItem{
					title: item.Name,
					desc: fmt.Sprintf("Created: %s, Updated: %s",
						item.CreatedAt.Format("15:04:05 02.01.2006"),
						item.UpdatedAt.Format("15:04:05 02.01.2006")),
					itemScheme: item,
					itemID:     item.ID})
			}
			m.TabContent[1].SetItems(items)
			return m, nil
		}
	case schemes.GetUserBankCardListMsg:
		if msg.StatusCode == http.StatusUnauthorized {
			return login.New(m.client, New(m.client), msg.Err.Error()), nil
		} else if msg.Err != nil {
			m.errMsg = msg.Err.Error()
			return m, nil
		} else {
			var items []list.Item = m.TabContent[2].Items()
			for _, item := range msg.List {
				items = append(items, listItem{
					title: item.Name,
					desc: fmt.Sprintf("Created: %s, Updated: %s",
						item.CreatedAt.Format("15:04:05 02.01.2006"),
						item.UpdatedAt.Format("15:04:05 02.01.2006")),
					itemScheme: item,
					itemID:     item.ID,
				})
			}
			m.TabContent[2].SetItems(items)
			return m, nil
		}
	case schemes.GetUserAuthInfoListMsg:
		if msg.StatusCode == http.StatusUnauthorized {
			return login.New(m.client, New(m.client), msg.Err.Error()), nil
		} else if msg.Err != nil {
			m.errMsg = msg.Err.Error()
			return m, nil
		} else {
			var items []list.Item = m.TabContent[3].Items()

			for _, item := range msg.List {
				items = append(items, listItem{
					title: item.Name,
					desc: fmt.Sprintf("Created: %s, Updated: %s",
						item.CreatedAt.Format("15:04:05 02.01.2006"),
						item.UpdatedAt.Format("15:04:05 02.01.2006")),
					itemScheme: item,
					itemID:     item.ID,
				})
			}

			m.TabContent[3].SetItems(items)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.TabContent[m.activeTab], cmd = m.TabContent[m.activeTab].Update(msg)
	return m, cmd
}

func (m homeModel) View() string {

	sb := strings.Builder{}
	renderedTabs := m.renderTabs()
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	sb.WriteString(row)
	sb.WriteString("\n")
	sb.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render("\n" + m.TabContent[m.activeTab].View() + helpStyle.Render("\n"+m.TabHelps[m.activeTab]+"\nctrl+o logout")))
	if m.errMsg != "" {
		sb.WriteString("\n\n" + errStyle.Render(m.errMsg))
	}
	return tabDocStyle.Render(sb.String())
}

func (m homeModel) renderTabs() []string {
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
