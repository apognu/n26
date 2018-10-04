package api

import (
	"net/url"
	"strings"

	"github.com/apognu/n26/cli"
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

func getSpaceFromID(spaces *cli.Spaces, id string) *cli.Space {
	for _, sp := range spaces.Spaces {
		if sp.ID == id || sp.Name == id {
			return &sp
		}
	}
	return nil
}
