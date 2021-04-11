package apps

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
	"octiiCli/utils"
	"time"
)

var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	keyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#43A9FF"))
	dataStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#874BFD"))
	validStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00AA00"))
)

func Account(_ *cli.Context) error {
	token, err := utils.Authorization()
	if err != nil {
		return &errorString{errorStyle.Render("You aren't authenticated! To login please use 'octii account login'")}
	}

	claims, err := utils.ParseAuthorization(token)

	if err != nil {
		return &errorString{errorStyle.Render("Couldn't parse token! Please reauthenticate using 'octii account login'")}
	}

	var valid string

	if claims.ExpiresAt.Unix() <= time.Now().Unix() {
		valid = errorStyle.Render("✗")
	} else {
		valid = validStyle.Render("✓")
	}

	println(utils.List(map[string]string {
		keyStyle.Render("ID"): dataStyle.Render(claims.Subject),
		keyStyle.Render("Expires"): dataStyle.Render(claims.ExpiresAt.String()),
		keyStyle.Render("Valid"): valid,
	}, 50))

	return nil
}