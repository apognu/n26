package api

import (
	"fmt"
	"net/http"
	"time"
)

func (cl *N26Client) GetPersonalInformation(meta *Metadata) (*PersonalInformation, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    "/api/me",
		Decoder: NewJSON(new(PersonalInformation)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if info, ok := output.(*PersonalInformation); ok {
		return info, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) GetAccount(meta *Metadata) (*Account, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    "/api/accounts",
		Decoder: NewJSON(new(Account)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if balance, ok := output.(*Account); ok {
		return balance, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) GetBalance(meta *Metadata) (*Balance, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    "/api/accounts",
		Decoder: NewJSON(new(Balance)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if balance, ok := output.(*Balance); ok {
		return balance, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) GetSpaces(meta *Metadata) (*Spaces, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    "/api/spaces",
		Decoder: NewJSON(new(Spaces)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if spaces, ok := output.(*Spaces); ok {
		return spaces, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) GetCategories() (map[string]string, error) {
	req := &N26Request{
		Method:  http.MethodGet,
		Path:    fmt.Sprintf("/api/smrt/categories"),
		Decoder: NewJSON(new([]Category)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if cats, ok := output.(*[]Category); ok {
		categories := make(map[string]string)
		for _, cat := range *cats {
			categories[cat.ID] = cat.Name
		}

		return categories, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) GetStatistics(meta *Metadata, from, to string) (*Statistics, error) {
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
		Decoder: NewJSON(new(Statistics)),
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if stats, ok := output.(*Statistics); ok {
		if meta.GetCategories() == nil {
			return nil, fmt.Errorf("could not get categories")
		}

		stats.Categories = meta.Categories

		return stats, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}
