package api

import (
	"fmt"
	"net/http"

	"github.com/apognu/n26/cli"
)

func (cl *N26Client) GetCards(meta *cli.Metadata) (cli.CardList, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    "/api/v2/cards",
		Decoder: NewJSON(new(cli.CardList)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if cards, ok := output.(*cli.CardList); ok {
		return *cards, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) GetLimits(meta *cli.Metadata) (cli.LimitList, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    "/api/settings/account/limits",
		Decoder: NewJSON(new(cli.LimitList)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if limits, ok := output.(*cli.LimitList); ok {
		return *limits, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}
