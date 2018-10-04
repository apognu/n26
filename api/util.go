package api

import (
	"net/url"
	"strings"

	"github.com/fatih/color"
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
