package apps

import (
	"bytes"
	"encoding/json"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/havoc-io/go-keytar"
	"github.com/urfave/cli/v2"
	"net/http"
)

var (
	blurredPrompt = "> "
	focusedPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#43A9FF")).Render("> ")
	blurredButton = "Login"
	focusedButton = lipgloss.NewStyle().Foreground(lipgloss.Color("#43A9FF")).Bold(true).Render("Login")
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#43A9FF"))
	titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("#43A9FF")).Foreground(lipgloss.Color("#FFFFFF")).Align(lipgloss.Center).Width(24).PaddingTop(1).PaddingBottom(1).PaddingLeft(4).PaddingRight(4)
	successStyle = lipgloss.NewStyle().Background(lipgloss.Color("#00AA00")).Foreground(lipgloss.Color("#FFFFFF")).Align(lipgloss.Center).Width(24).PaddingTop(1).PaddingBottom(1).PaddingLeft(4).PaddingRight(4)
)

type loginModel struct {
	emailInput textinput.Model
	passwordInput textinput.Model
	index int
	loginButton string
	error error
	success bool
	loading bool
	spinner spinner.Model
}

func newLoginModel() loginModel {
	model := loginModel{}

	model.emailInput = textinput.NewModel()
	model.emailInput.Placeholder = "email"
	model.emailInput.Focus()
	model.emailInput.Prompt = focusedPrompt

	model.passwordInput = textinput.NewModel()
	model.passwordInput.Prompt = blurredPrompt
	model.passwordInput.Placeholder = "password"
	model.passwordInput.EchoMode = textinput.EchoPassword
	model.passwordInput.EchoCharacter = '•'

	model.loginButton = blurredButton

	model.spinner = spinner.NewModel()
	model.spinner.Spinner.Frames = []string{
		"⠋",
		"⠙",
		"⠹",
		"⠸",
		"⠼",
		"⠴",
		"⠦",
		"⠧",
		"⠇",
		"⠏",
	}

	model.success = false
	model.loading = false

	return model
}

type errorMsg struct {
	error error
}
type completedMsg struct {}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func loginWithServer(email string, password string) tea.Cmd {
	return func() tea.Msg {
		body := map[string]string{"email": email, "password": password}
		bodyJson, err := json.Marshal(body)

		if err != nil {
			return errorMsg{
				err,
			}
		}

		response, err := http.Post("https://gateway.octii.chat/users/login", "application/json", bytes.NewBuffer(bodyJson))

		if err != nil {
			return errorMsg{
				err,
			}
		}

		if response.StatusCode != 200 {
			return errorMsg{
				error: &errorString{
					s: "Could not auth!",
				},
			}
		}

		var res struct {
			Authorization string `json:"authorization"`
		}

		if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
			return errorMsg{
				err,
			}
		}

		keychain, err := keytar.GetKeychain()

		if err != nil {
			return errorMsg{
				err,
			}
		}

		if err := keytar.ReplacePassword(keychain, "octii", "authorization", res.Authorization); err != nil {
			return errorMsg{
				err,
			}
		}

		return completedMsg{}
	}
}

func (m loginModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, spinner.Tick)
}

func (m loginModel) View() string {
	if m.success {
		return successStyle.Render("Success!\n")
	}

	s := titleStyle.Render("Login with Octii")
	s += "\n\n"

	inputs := []string {
		m.emailInput.View(),
		m.passwordInput.View(),
	}

	for index, view := range inputs {
		s += view
		if index < len(inputs) - 1 {
			s += "\n"
		}
	}

	s += "\n\n" + m.loginButton

	if m.loading {
		s += " " + spinnerStyle.Render(m.spinner.View())
	} else if m.error != nil {
		s += " " + lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("✗")
	}

	return s
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			{
				inputs := []textinput.Model{
					m.emailInput,
					m.passwordInput,
				}

				s := msg.String()

				// Did the user press enter while the submit button was focused?
				// If so, exit.
				if s == "enter" && m.index == len(inputs) && !m.loading {
					m.loading = true
					return m, loginWithServer(m.emailInput.Value(), m.passwordInput.Value())
				}

				// Cycle indexes
				if s == "up" || s == "shift+tab" {
					m.index--
				} else {
					m.index++
				}

				if m.index > len(inputs) {
					m.index = 0
				} else if m.index < 0 {
					m.index = len(inputs)
				}

				for i := 0; i <= len(inputs)-1; i++ {
					if i == m.index {
						// Set focused state
						inputs[i].Focus()
						inputs[i].Prompt = focusedPrompt
						continue
					}
					// Remove focused state
					inputs[i].Blur()
					inputs[i].Prompt = blurredPrompt
					inputs[i].TextColor = ""
				}

				m.emailInput = inputs[0]
				m.passwordInput = inputs[1]

				if m.index == len(inputs) {
					m.loginButton = focusedButton
				} else {
					m.loginButton = blurredButton
				}

				return m, nil
			}
		}
	case errorMsg: {
		m.loading = false
		m.error = msg.error
		return m, nil
	}

	case completedMsg: {
		m.loading = false
		m.success = true
		return m, tea.Quit
	}
	}

	var (
		cmd tea.Cmd
		commands []tea.Cmd
	)

	m.emailInput, cmd = m.emailInput.Update(msg)
	commands = append(commands, cmd)
	m.passwordInput, cmd = m.passwordInput.Update(msg)
	commands = append(commands, cmd)
	m.spinner, cmd = m.spinner.Update(msg)
	commands = append(commands, cmd)

	return m, tea.Batch(commands...)
}

func Login(_ *cli.Context) error {
	loginModel := newLoginModel()

	p := tea.NewProgram(loginModel)
	if err := p.Start(); err != nil {
		return &errorString{errorStyle.Render("An unexpected error happened in the login process")}
	}
	return nil
}
