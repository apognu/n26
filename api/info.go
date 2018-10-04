package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/apognu/n26/cli"
)

func (cl *N26Client) GetPersonalInformation(meta *cli.Metadata) (*cli.PersonalInformation, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    "/api/me",
		Decoder: NewJSON(new(cli.PersonalInformation)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if info, ok := output.(*cli.PersonalInformation); ok {
		return info, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) GetAccount(meta *cli.Metadata) (*cli.Account, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    "/api/accounts",
		Decoder: NewJSON(new(cli.Account)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if balance, ok := output.(*cli.Account); ok {
		return balance, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) GetBalance(meta *cli.Metadata) (*cli.Balance, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    "/api/accounts",
		Decoder: NewJSON(new(cli.Balance)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if balance, ok := output.(*cli.Balance); ok {
		return balance, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) GetSpaces(meta *cli.Metadata) (*cli.Spaces, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    "/api/spaces",
		Decoder: NewJSON(new(cli.Spaces)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if spaces, ok := output.(*cli.Spaces); ok {
		return spaces, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) GetCategories() (map[string]string, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    fmt.Sprintf("/api/smrt/categories"),
		Decoder: NewJSON(new([]cli.Category)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if cats, ok := output.(*[]cli.Category); ok {
		categories := make(map[string]string)
		for _, cat := range *cats {
			categories[cat.ID] = cat.Name
		}

		return categories, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) GetStatistics(meta *cli.Metadata, from, to string) (*cli.Statistics, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	if from != "" {
		var serr, eerr error
		start, serr = time.Parse("2006-01-02", from)
		end, eerr = time.Parse("2006-01-02", to)
		if serr != nil || eerr != nil {
			return nil, fmt.Errorf("could not parse the provided dates")
		}
	}

	req := &N26Request{
		Method:  http.MethodGet,
		Path:    fmt.Sprintf("/api/smrt/statistics/categories/%d/%d", start.Unix()*1000, end.Unix()*1000),
		Decoder: NewJSON(new(cli.Statistics)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if stats, ok := output.(*cli.Statistics); ok {
		if meta.GetCategories() == nil {
			return nil, fmt.Errorf("could not get categories")
		}

		stats.Categories = meta.Categories

		return stats, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}
