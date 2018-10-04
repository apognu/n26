package api

import (
	"fmt"
	"net/http"
)

func (cl *N26Client) GetCards(meta *Metadata) (CardList, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    "/api/v2/cards",
		Decoder: NewJSON(new(CardList)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if cards, ok := output.(*CardList); ok {
		return *cards, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) GetLimits(meta *Metadata) (LimitList, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    "/api/settings/account/limits",
		Decoder: NewJSON(new(LimitList)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if limits, ok := output.(*LimitList); ok {
		return *limits, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}
