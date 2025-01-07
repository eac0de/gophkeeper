package login

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

import (
	"fmt"
	"net"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eac0de/gophkeeper/client/internal/client"
	"github.com/eac0de/gophkeeper/client/internal/components"
)

var (
	errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

type loginModel struct {
	email        string
	emailCodeId  string
	invalidInput bool
	state        int
	input        textinput.Model
	errMsg       string
	client       *client.APIClient
	nextModel    tea.Model
}

func New(apiClient *client.APIClient, nextModel tea.Model, errMsg string) loginModel {
	// Здесь будет проверка на существование файла в tmp, в файле будет храниться информация о состоянии программы
	t := textinput.New()
	t.Focus()
	return loginModel{input: t, client: apiClient, nextModel: nextModel, errMsg: errMsg}
}

func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.errMsg = ""
			if m.state == 0 {
				if validateEmail(m.input.Value()) {
					var err error
					m.email = m.input.Value()
					m.emailCodeId, err = m.client.GenerateEmailCode(m.email)
					if err != nil {
						m.errMsg = err.Error()
						return m, nil
					}
					m.invalidInput = false
					m.input.Reset()
					m.state = 1
				} else {
					m.invalidInput = true
				}
			} else {
				length := utf8.RuneCountInString(m.input.Value())
				if length != 4 {
					m.invalidInput = true
					return m, nil
				}
				i, err := strconv.ParseInt(m.input.Value(), 10, 64)
				if err != nil {
					m.invalidInput = true
					return m, nil
				}
				statusCode, err := m.client.VerifyEmailCode(m.emailCodeId, int(i))
				if statusCode == 0 && err != nil {
					m.errMsg = err.Error()
					return m, nil
				}
				if statusCode != http.StatusOK && statusCode != http.StatusCreated {
					if statusCode == http.StatusPreconditionFailed {
						m.invalidInput = true
						return m, nil
					}
					if statusCode == http.StatusGone {
						m.state = 0
						m.errMsg = "Code expired"
						m.input.Reset()
						return m, nil
					}
					if err != nil {
						m.errMsg = err.Error()
						return m, nil
					}
				}
				if m.nextModel != nil {
					var cmd tea.Cmd
					msg := m.nextModel.Init()()
					if msg != nil {
						m.nextModel, cmd = m.nextModel.Update(msg)
					}
					return m.nextModel, cmd
				}
				return m, nil
			}
			return m, nil
		case "esc":
			m.state = 0
			m.input.SetValue(m.email)
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		}

	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.invalidInput = false
	return m, cmd
}

func (m loginModel) View() string {
	var b strings.Builder
	b.WriteString(components.GetLogo() + "\n")
	if m.state == 0 {
		m.input.CharLimit = 64
		b.WriteString("Enter your email:\n")
	} else {
		m.input.CharLimit = 4
		b.WriteString(fmt.Sprintf("Enter code from email(%s):\n", m.email))
	}
	if m.invalidInput {
		b.WriteString(errStyle.Render(m.input.View()))
	} else {
		b.WriteString(m.input.View())
	}
	if m.errMsg != "" {
		b.WriteString("\n\n" + errStyle.Render(m.errMsg))
	}
	return b.String()
}

func validateEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]
	_, err = net.LookupMX(domain)
	return err == nil
}
