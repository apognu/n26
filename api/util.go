package api

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	titleColor = color.New(color.Bold, color.FgBlue)
	attrColor  = color.New(color.Faint)

	okColor   = color.New(color.FgGreen)
	warnColor = color.New(color.FgYellow)
	errColor  = color.New(color.FgRed)
)

func query(params map[string]string) url.Values {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values
}

func idToLabel(id string) string {
	return strings.Title(strings.Replace(strings.ToLower(id), "_", " ", -1))
}

func getSpaceFromID(spaces *Spaces, id string) *Space {
	for _, sp := range spaces.Spaces {
		if sp.ID == id || sp.Name == id {
			return &sp
		}
	}
	return nil
}

func readLine(prompt string) string {
	fmt.Print(fmt.Sprintf("%s ", prompt))

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}

func readSecret(prompt string) (string, error) {
	fmt.Print(prompt)
	secret, err := terminal.ReadPassword(syscall.Stdin)

	line()

	if err != nil {
		return "", fmt.Errorf("could not read PIN")
	}

	return string(secret), nil
}
