package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apognu/n26/cli"
)

func (cl *N26Client) GetPastTransactions(meta *cli.Metadata, from, to string, limit int) (cli.PastTransactionList, error) {
	now := time.Now()
	defaultFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	defaultTo := defaultFrom.AddDate(0, 1, 0).Add(-time.Nanosecond)

	req := &N26Request{
		Method:  http.MethodGet,
		Path:    "/api/smrt/transactions",
		Decoder: NewJSON(new(cli.PastTransactionList)),
		Params: map[string]string{
			"from": fmt.Sprintf("%d", defaultFrom.Unix()*1000),
			"to":   fmt.Sprintf("%d", defaultTo.Unix()*1000),
		},
	}

	req.Params["limit"] = fmt.Sprintf("%d", limit)

	if from != "" {
		if to == "" {
			return nil, fmt.Errorf("both 'from' and 'to' must be provided")
		}

		f, ferr := time.Parse("2006-01-02", from)
		t, terr := time.Parse("2006-01-02", to)
		if ferr != nil || terr != nil {
			return nil, fmt.Errorf("could not parse provided dates")
		}

		req.Params["from"] = fmt.Sprint(f.Unix() * 1000)
		req.Params["to"] = fmt.Sprint(t.Unix() * 1000)
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if transactions, ok := output.(*cli.PastTransactionList); ok {
		return *transactions, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) GetContacts() (cli.ContactList, error) {
	req := &N26Request{
		Path:    "/api/smrt/contacts",
		Method:  http.MethodGet,
		Decoder: NewJSON(new(cli.ContactList)),
		Params:  map[string]string{},
	}

	output, err := cl.Request(req, false)
	if err != nil {
		return nil, err
	}

	if contacts, ok := output.(*cli.ContactList); ok {
		return *contacts, nil
	}

	return nil, fmt.Errorf("could not unmarshal upstream data")
}

func (cl *N26Client) CheckContact(id string) bool {
	req := &N26Request{
		Method:  http.MethodPost,
		Path:    "/api/contacts",
		Body:    []string{id},
		Decoder: NewJSON(new([]cli.ContactRequest)),
	}

	body, err := cl.Request(req, false)
	if err != nil {
		return false
	}

	if len(*body.(*[]cli.ContactRequest)) == 0 {
		return false
	}

	return true
}

func (cl *N26Client) CreateSpaceTransfer(meta *cli.Metadata, from, to string, amount float64) (cli.SimpleMessage, error) {
	spaces, err := cl.GetSpaces(meta)
	if err != nil {
		return "", fmt.Errorf("could not get your spaces")
	}

	fromSpace, toSpace := getSpaceFromID(spaces, from), getSpaceFromID(spaces, to)
	if fromSpace == nil || toSpace == nil {
		return "", fmt.Errorf("could not find the provided spaces")
	}

	cli.ConfirmSpaceTransfer(fromSpace, toSpace, amount)

	trx := cli.SpaceTransaction{
		FromSpaceID: fromSpace.ID,
		ToSpaceID:   toSpace.ID,
		Amount:      amount,
	}

	req := &N26Request{
		Method: http.MethodPost,
		Path:   "/api/spaces/transaction",
		Body:   trx,
	}

	_, err = cl.Request(req, false)
	if err != nil {
		return "", err
	}

	return (cli.SimpleMessage)(fmt.Sprintf("Your transfer of %s has been performed.", cli.Curr(amount, fromSpace.Balance.Currency))), nil
}

func (cl *N26Client) CreateMoneyBeam(meta *cli.Metadata, name, recipient string, amount float64, comment string) (cli.SimpleMessage, error) {
	if !cl.CheckContact(recipient) {
		return "", fmt.Errorf("the provided recipient ID is not associated with an N26 account")
	}

	details := cli.MoneyBeamDetails{Type: "FT", PartnerName: name, Amount: amount, Comment: comment}

	if name == "" {
		details.PartnerName = recipient
	}

	if strings.Contains(recipient, "@") {
		details.PartnerEmail = recipient
	} else if strings.HasPrefix(recipient, "+") {
		details.PartnerPhone = recipient
	} else {
		return "", fmt.Errorf("the recipient must be an email address or a phone number (starting with '+')")
	}

	balance, err := cl.GetBalance(meta)
	if err != nil {
		return "", fmt.Errorf("could not get current balance")
	}

	cli.ConfirmMoneyBeam(details, balance)

	pin, err := cli.ReadSecret("Enter your PIN:")
	if err != nil {
		return "", fmt.Errorf("could not read PIN")
	}

	trx := cli.MoneyBeam{
		PIN:         string(pin),
		Transaction: details,
	}

	req := &N26Request{
		Method: http.MethodPost,
		Path:   "/api/transactions",
		Body:   trx,
	}

	_, err = cl.Request(req, false)
	if err != nil {
		return "", err
	}

	return (cli.SimpleMessage)(fmt.Sprintf("Your transfer of %s has been requested, please confirm from your paired device.", cli.Curr(trx.Transaction.Amount, balance.Currency))), nil
}
