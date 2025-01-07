package item

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eac0de/gophkeeper/client/internal/client"
	"github.com/eac0de/gophkeeper/client/internal/components"
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
	tabDocStyle       = lipgloss.NewStyle().Padding(0, 0, 0, 0)

	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 0)
	activeTabStyle   = inactiveTabStyle.Border(activeTabBorder, true)
	errStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	windowStyle      = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(0, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	focusedStyle     = lipgloss.NewStyle().Foreground(highlightColor)
	blurredStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	successStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Align(lipgloss.Center)
	cursorStyle      = focusedStyle
	noStyle          = lipgloss.NewStyle()
	helpStyle        = blurredStyle
)

type itemModel struct {
	// Tabs
	Tabs      []string
	activeTab int

	// Input
	focusIndex    int
	inputs        []textinput.Model
	invalidInputs []bool

	// Custom
	itemID           uuid.UUID
	item             interface{}
	client           *client.APIClient
	errMsg           string
	nextModel        tea.Model
	help             string
	itemIndex        int
	fileIsDownloaded bool
	addInfo          string
	successMsg       string
}

func New(client *client.APIClient, item interface{}, activeTab int, nextModel tea.Model, itemIndex int) tea.Model {
	m := itemModel{
		client:        client,
		Tabs:          components.Tabs,
		invalidInputs: make([]bool, 5),
		nextModel:     nextModel,
		itemIndex:     itemIndex,
		item:          item,
		activeTab:     activeTab,
		addInfo:       "\n",
	}

	return m.initInputs()
}

type clearSuccessMsg struct{}

func clearSuccessMsgAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearSuccessMsg{}
	})
}

type clearInvalidInputsMsg struct{}

func clearInvalidInputsMsgAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearInvalidInputsMsg{}
	})
}

func (m itemModel) initInputs() tea.Model {
	switch m.activeTab {
	case 0:
		var textData schemes.UserTextData
		if m.item != nil {
			var ok bool
			textData, ok = m.item.(schemes.UserTextData)
			if !ok {
				return m.nextModel
			}
			m.itemID = textData.ID
			m.addInfo = fmt.Sprintf("\n\nCreated:\n %s\n\nUpdated:\n %s\n", textData.CreatedAt.Format("15:04:05 02.01.2006"), textData.UpdatedAt.Format("15:04:05 02.01.2006"))
		}
		m.inputs = make([]textinput.Model, 2)
		var t textinput.Model
		for i := range m.inputs {
			t = textinput.New()
			t.Cursor.Style = cursorStyle
			t.CharLimit = 100
			switch i {
			case 0:
				t.SetValue(textData.Name)
				t.Focus()
				t.PromptStyle = focusedStyle
				t.TextStyle = focusedStyle
				t.Prompt = "Name:\n"
			case 1:
				t.SetValue(textData.Data)
				t.CharLimit = 1024
				t.Prompt = "Text:\n"
			}
			m.inputs[i] = t
		}
		m.help = "↑/↓ up/down • ctrl+s save • ctrl+d delete • ctrl+c exit"
	case 1:
		var fileData schemes.UserFileData
		if m.item == nil {
			return m.nextModel
		}
		var ok bool
		fileData, ok = m.item.(schemes.UserFileData)
		if !ok {
			return m.nextModel
		}
		m.itemID = fileData.ID
		m.inputs = make([]textinput.Model, 1)
		var t textinput.Model
		for i := range m.inputs {
			t = textinput.New()
			t.Cursor.Style = cursorStyle
			t.CharLimit = 100
			switch i {
			case 0:
				t.SetValue(fileData.Name)
				t.Focus()
				t.PromptStyle = focusedStyle
				t.TextStyle = focusedStyle
				t.Prompt = "Name:\n"
			}
			m.inputs[i] = t
		}
		m.help = "esc back • ↑/↓ up/down • ctrl+s save • ctrl+d delete • tab download • ctrl+c exit"
		m.addInfo = fmt.Sprintf("\n\nExtension:\n %s\n\nCreated:\n %s\n\nUpdated:\n %s\n", fileData.Ext, fileData.CreatedAt.Format("15:04:05 02.01.2006"), fileData.UpdatedAt.Format("15:04:05 02.01.2006"))

	case 2:
		var bankCard schemes.UserBankCard
		if m.item != nil {
			var ok bool
			bankCard, ok = m.item.(schemes.UserBankCard)
			if !ok {
				return m.nextModel
			}
			m.itemID = bankCard.ID
			m.addInfo = fmt.Sprintf("\n\nCreated:\n %s\n\nUpdated:\n %s\n", bankCard.CreatedAt.Format("15:04:05 02.01.2006"), bankCard.UpdatedAt.Format("15:04:05 02.01.2006"))
		}
		m.inputs = make([]textinput.Model, 5)
		var t textinput.Model
		for i := range m.inputs {
			t = textinput.New()
			t.Cursor.Style = cursorStyle
			t.CharLimit = 100
			switch i {
			case 0:
				t.SetValue(bankCard.Name)
				t.Focus()
				t.PromptStyle = focusedStyle
				t.TextStyle = focusedStyle
				t.Prompt = "Name:\n"
			case 1:
				t.SetValue(bankCard.Number)
				t.Prompt = "Number:\n"
			case 2:
				t.SetValue(bankCard.CardHolder)
				t.Prompt = "Card Holder:\n"
			case 3:
				t.SetValue(bankCard.ExpireDate)
				t.Prompt = "Expire Date:\n"
			case 4:
				t.SetValue(bankCard.CSC)
				t.Prompt = "CSC:\n"
			}
			m.inputs[i] = t
		}
		m.help = "esc back • ↑/↓ up/down • ctrl+s save • ctrl+d delete • ctrl+c exit"
	case 3:
		var authInfo schemes.UserAuthInfo
		if m.item != nil {
			var ok bool
			authInfo, ok = m.item.(schemes.UserAuthInfo)
			if !ok {
				return m.nextModel
			}
			m.itemID = authInfo.ID
			m.addInfo = fmt.Sprintf("\n\nCreated:\n %s\n\nUpdated:\n %s\n", authInfo.CreatedAt.Format("15:04:05 02.01.2006"), authInfo.UpdatedAt.Format("15:04:05 02.01.2006"))
		}
		m.inputs = make([]textinput.Model, 3)
		var t textinput.Model
		for i := range m.inputs {
			t = textinput.New()
			t.Cursor.Style = cursorStyle
			t.CharLimit = 100
			switch i {
			case 0:
				t.SetValue(authInfo.Name)
				t.Focus()
				t.PromptStyle = focusedStyle
				t.TextStyle = focusedStyle
				t.Prompt = "Name:\n"
			case 1:
				t.SetValue(authInfo.Login)
				t.Prompt = "Login:\n"
			case 2:
				t.SetValue(authInfo.Password)
				t.Prompt = "Password:\n"
			}
			m.inputs[i] = t
		}
		m.help = "esc back • ↑/↓ up/down • ctrl+s save • ctrl+d delete • ctrl+c exit"
	}
	return m
}

func (m itemModel) Init() tea.Cmd {
	return nil
}

func (m itemModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case clearSuccessMsg:
		m.successMsg = ""
		return m, nil
	case clearInvalidInputsMsg:
		for i := range m.invalidInputs {
			m.invalidInputs[i] = false
		}
		cmds := make([]tea.Cmd, len(m.inputs))
		for i := range m.inputs {
			if i == m.focusIndex {
				// Set focused state
				cmds[i] = m.inputs[i].Focus()
				m.inputs[i].PromptStyle = focusedStyle
				m.inputs[i].TextStyle = focusedStyle
				continue
			}
			// Remove focused state
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].TextStyle = noStyle
		}
		return m, tea.Batch(cmds...)
	case schemes.DownloadFileMsg:
		if msg.Err != nil {
			if msg.StatusCode == http.StatusUnauthorized {
				return login.New(m.client, m, msg.Err.Error()), nil
			}
			m.errMsg = msg.Err.Error()
			return m, nil
		}
		m.fileIsDownloaded = true
		m.successMsg = "File downloaded successfully"
		return m, clearSuccessMsgAfter(3 * time.Second)
	case schemes.SaveTextDataMsg:
		if msg.Err != nil {
			if msg.StatusCode == http.StatusUnauthorized {
				return login.New(m.client, m, msg.Err.Error()), nil
			}
			m.errMsg = msg.Err.Error()
			return m, nil
		}
		m.item = msg.TextData
		m.successMsg = "Text data saved successfully"
		return m.initInputs(), clearSuccessMsgAfter(3 * time.Second)
	case schemes.SaveFileDataMsg:
		if msg.Err != nil {
			if msg.StatusCode == http.StatusUnauthorized {
				return login.New(m.client, m, msg.Err.Error()), nil
			}
			m.errMsg = msg.Err.Error()
			return m, nil
		}
		m.item = msg.FileData
		m.successMsg = "Filename saved successfully"
		return m.initInputs(), clearSuccessMsgAfter(3 * time.Second)
	case schemes.SaveBankCardMsg:
		if msg.Err != nil {
			if msg.StatusCode == http.StatusUnauthorized {
				return login.New(m.client, m, msg.Err.Error()), nil
			}
			m.errMsg = msg.Err.Error()
			return m, nil
		}
		m.item = msg.BankCard
		m.successMsg = "Bank card saved successfully"
		return m.initInputs(), clearSuccessMsgAfter(3 * time.Second)
	case schemes.SaveAuthInfoMsg:
		if msg.Err != nil {
			if msg.StatusCode == http.StatusUnauthorized {
				return login.New(m.client, m, msg.Err.Error()), nil
			}
			m.errMsg = msg.Err.Error()
			return m, nil
		}
		m.item = msg.AuthInfo
		m.successMsg = "Auth info saved successfully"
		return m.initInputs(), clearSuccessMsgAfter(3 * time.Second)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.nextModel, func() tea.Msg {
				return schemes.UpdateListItemMsg{ItemIndex: m.itemIndex, ActiveTab: m.activeTab, Item: m.item}
			}
		case "ctrl+c":
			utils.SaveTokens(m.client.Tokens)
			return m, tea.Quit
		case "ctrl+s":
			switch m.activeTab {
			case 0:
				var isValid bool
				var hasErr bool
				for i := range m.inputs {
					switch i {
					case 0:
						isValid = utils.StringIsValid(m.inputs[i].Value())
						m.invalidInputs[i] = !isValid
						if !isValid {
							hasErr = true
						}
					case 1:
						isValid = utils.StringIsValid(m.inputs[i].Value())
						m.invalidInputs[i] = !isValid
						if !isValid {
							hasErr = true
						}
					}
				}
				if hasErr {
					return m, clearInvalidInputsMsgAfter(1 * time.Second)
				}
				return m, func() tea.Msg {
					return m.client.SaveUserTextData(m.itemID, m.inputs[0].Value(), m.inputs[1].Value(), nil)
				}
			case 1:
				var isValid bool
				var hasErr bool
				for i := range m.inputs {
					switch i {
					case 0:
						isValid = utils.StringIsValid(m.inputs[i].Value())
						m.invalidInputs[i] = !isValid
						if !isValid {
							hasErr = true
						}
					}
				}
				if hasErr {
					return m, clearInvalidInputsMsgAfter(1 * time.Second)
				}
				return m, func() tea.Msg {
					return m.client.SaveUserFileData(m.itemID, m.inputs[0].Value(), "", nil)
				}
			case 2:
				var isValid bool
				var hasErr bool
				for i := range m.inputs {
					switch i {
					case 0:
						isValid = utils.StringIsValid(m.inputs[i].Value())
						m.invalidInputs[i] = !isValid
						if !isValid {
							hasErr = true
						}
					case 1:
						isValid = utils.CardNumberIsValid(m.inputs[i].Value())
						m.invalidInputs[i] = !isValid
						if !isValid {
							hasErr = true
						}
					case 2:
						isValid = utils.StringIsValid(m.inputs[i].Value())
						m.invalidInputs[i] = !isValid
						if !isValid {
							hasErr = true
						}
					case 3:
						isValid = utils.ExpireDateIsValid(m.inputs[i].Value())
						m.invalidInputs[i] = !isValid
						if !isValid {
							hasErr = true
						}
					case 4:
						isValid = utils.CSCIsValid(m.inputs[i].Value())
						m.invalidInputs[i] = !isValid
						if !isValid {
							hasErr = true
						}
					}
				}
				if hasErr {
					return m, clearInvalidInputsMsgAfter(1 * time.Second)
				}
				return m, func() tea.Msg {
					return m.client.SaveUserBankCard(m.itemID, m.inputs[0].Value(), m.inputs[1].Value(), m.inputs[2].Value(), m.inputs[3].Value(), m.inputs[4].Value(), nil)
				}
			case 3:
				var isValid bool
				var hasErr bool
				for i := range m.inputs {
					switch i {
					case 0:
						isValid = utils.StringIsValid(m.inputs[i].Value())
						m.invalidInputs[i] = !isValid
						if !isValid {
							hasErr = true
						}
					case 1:
						isValid = utils.StringIsValid(m.inputs[i].Value())
						m.invalidInputs[i] = !isValid
						if !isValid {
							hasErr = true
						}
					case 2:
						isValid = utils.StringIsValid(m.inputs[i].Value())
						m.invalidInputs[i] = !isValid
						if !isValid {
							hasErr = true
						}
					}
				}
				if hasErr {
					return m, clearInvalidInputsMsgAfter(1 * time.Second)
				}
				return m, func() tea.Msg {
					return m.client.SaveUserAuthInfo(m.itemID, m.inputs[0].Value(), m.inputs[1].Value(), m.inputs[2].Value(), nil)
				}
			}
		case "ctrl+d":
			if m.itemID == uuid.Nil {
				return m.nextModel, nil
			}
			return m.nextModel, func() tea.Msg { return m.client.DeleteData(m.itemID, m.itemIndex, m.activeTab) }
		case "tab":
			if m.activeTab != 1 {
				return m, nil
			}
			return m, func() tea.Msg { return m.client.DownloadFile(m.itemID) }
		case "up", "down":
			s := msg.String()
			if s == "up" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}
			if m.focusIndex > len(m.inputs)-1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *itemModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m itemModel) View() string {
	for i := range m.inputs {
		log.Println(m.invalidInputs)
		if m.invalidInputs[i] {
			m.inputs[i].PromptStyle = errStyle
			m.inputs[i].TextStyle = errStyle
		}
	}
	inputsBuilder := strings.Builder{}

	for i := range m.inputs {
		inputsBuilder.WriteRune('\n')
		inputsBuilder.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			inputsBuilder.WriteRune('\n')
		}
	}
	inputsBuilder.WriteString(m.addInfo)
	if m.successMsg != "" {
		inputsBuilder.WriteString(successStyle.Render("\n" + m.successMsg + "\n"))
	}
	sb := strings.Builder{}
	renderedTabs := m.renderTabs()
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	sb.WriteString(row)
	sb.WriteString("\n")
	sb.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(inputsBuilder.String() + helpStyle.Render("\n\n"+m.help)))

	if m.errMsg != "" {
		sb.WriteString("\n\n" + errStyle.Render(m.errMsg))
	}
	return tabDocStyle.Render(sb.String())
}

func (m itemModel) renderTabs() []string {
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
