package apps

import (
	"github.com/havoc-io/go-keytar"
	"github.com/urfave/cli/v2"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func Login(c *cli.Context) error {
	token := c.Args().First()

	keychain, err := keytar.GetKeychain()

	if err != nil {
		return err
	}

	if err := keytar.ReplacePassword(keychain, "octii", "authorization", token); err != nil {
		return err
	}

	return nil
}
