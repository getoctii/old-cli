package apps

import (
	"github.com/urfave/cli/v2"
	"io/ioutil"
)

const template =
	`{
  "name": "",
  "dark": {
    "colors": {
      "primary": "",
      "secondary": "",
      "success": "",
      "info": "",
      "danger": "",
      "warning": "",
      "light": "",
      "dark": ""
    },
    "mention": {
      "me": "",
      "other": ""
    },
    "text": {
      "normal": "",
      "inverse": "",
      "primary": "",
      "danger": "",
      "warning": "",
      "secondary": ""
    },
    "backgrounds": {
      "primary": "",
      "secondary": ""
    },
    "sidebar": {
      "background": "",
      "seperator": "",
      "shadow": ""
    },
    "settings": {
      "background": "",
      "card": "",
      "input": ""
    },
    "channels": {
      "background": "",
      "seperator": ""
    },
    "chat": {
      "background": "",
      "hover": ""
    },
    "context": {
      "seperator": "",
      "background": ""
    },
    "status": {
      "selected": "",
      "online": "",
      "idle": "",
      "dnd": "",
      "offline": ""
    },
    "message": {
      "author": "",
      "date": "",
      "message": ""
    },
    "input": {
      "background": "",
      "text": ""
    },
    "modal": {
      "background": "",
      "foreground": ""
    },
    "emojis": {
      "background": "",
      "input": ""
    }
  },
  "light": {
    "colors": {
      "primary": "",
      "secondary": "",
      "success": "",
      "info": "",
      "danger": "",
      "warning": "",
      "light": "",
      "dark": ""
    },
    "text": {
      "normal": "",
      "inverse": "",
      "primary": "",
      "danger": "",
      "warning": "",
      "secondary": ""
    },
    "mention": {
      "me": "",
      "other": ""
    },
    "sidebar": {
      "background": "",
      "seperator": "",
      "shadow": ""
    },
    "settings": {
      "background": "",
      "card": "",
      "input": ""
    },
    "channels": {
      "background": "",
      "seperator": ""
    },
    "backgrounds": {
      "primary": "",
      "secondary": ""
    },
    "chat": {
      "background": "",
      "hover": ""
    },
    "context": {
      "seperator": "",
      "background": ""
    },
    "status": {
      "selected": "",
      "online": "",
      "idle": "",
      "dnd": "",
      "offline": ""
    },
    "message": {
      "author": "",
      "date": "",
      "message": ""
    },
    "input": {
      "background": "",
      "text": ""
    },
    "modal": {
      "background": "",
      "foreground": ""
    },
    "emojis": {
      "background": "",
      "input": ""
    }
  }
}`

func Theme(c *cli.Context) error {
	outputFile := c.Args().Get(0)

	if err := ioutil.WriteFile(outputFile, []byte(template), 0777); err != nil {
		return &errorString {
			s: errorStyle.Render("Could not write file"),
		}
	}

	return nil
}