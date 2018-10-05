package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type js map[string]interface{}

func JSON(data interface{}) {
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		Fatal(err)
	}
	fmt.Println(string(j))
}

func (msg SimpleMessage) JSON(meta *Metadata) {
	logrus.Info(msg)
}

func (info PersonalInformation) JSON(meta *Metadata) {
	bd := time.Unix(info.BirthDate/1000, 0)

	JSON(js{
		"name":        fmt.Sprintf("%s %s", info.Firstname, info.Lastname),
		"email":       info.Email,
		"phone":       info.Phone,
		"birth_date":  bd.Format("02 Jan 2006"),
		"nationality": info.Nationality,
	})
}

func (account Account) JSON(meta *Metadata) {
	JSON(js{
		"bank": account.Bank,
		"iban": account.IBAN,
		"bic":  account.BIC,
	})
}

func (balance Balance) JSON(meta *Metadata) {
	JSON(js{
		"balance":        balance.AvailableBalance,
		"usable_balance": balance.UsageBalance,
	})
}

func (cards CardList) JSON(meta *Metadata) {
	data := make([]js, len(cards))

	for idx, card := range cards {
		exp := time.Unix(card.Expiration/1000, 0)
		model := card.ProductType
		if card.ProductType != card.Design {
			model = fmt.Sprintf("%s/%s", card.ProductType, card.Design)
		}

		data[idx] = js{
			"id":         card.ID,
			"holder":     card.Holder,
			"expiration": exp.Format("Jan 2006"),
			"type":       card.Type,
			"model":      model,
		}

		if s, ok := CardStatuses[card.Status]; ok {
			data[idx]["status"] = s.Text
		} else {
			data[idx]["status"] = card.Status
		}
	}

	JSON(data)
}

func (limits LimitList) JSON(meta *Metadata) {
}

func (transactions PastTransactionList) JSON(meta *Metadata) {
	data := make([]js, len(transactions))
	for idx, trx := range transactions {
		date := time.Unix(trx.Date/1000, 0)

		party := trx.MerchantName
		if trx.Partner != "" {
			party = trx.Partner
		}

		if trx.Scheme == "SPACES" {
			party = "N26 Spaces"
		}

		data[idx] = js{
			"date":        date.Format("02 Jan 2006 15:04"),
			"third_party": party,
			"amount":      trx.Amount,
			"category":    meta.GetCategory(trx.Category),
			"location":    trx.MerchantCity,
			"comment":     trx.Comment,
		}
	}

	JSON(data)
}

func (spaces Spaces) JSON(meta *Metadata) {
	data := make([]js, len(spaces.Spaces))

	for idx, space := range spaces.Spaces {
		data[idx] = js{
			"id":      space.ID,
			"name":    space.Name,
			"primary": true,
			"amount":  space.Balance.AvailableBalance,
		}

		if space.Goal.Amount > 0 {
			progress := space.Balance.AvailableBalance / space.Goal.Amount

			data[idx]["goal"] = space.Goal.Amount
			data[idx]["progress"] = progress
		}
	}

	JSON(data)
}

func (stats Statistics) JSON(meta *Metadata) {
	data := js{
		"global": js{
			"income":  stats.TotalIncome,
			"expense": stats.TotalExpense,
		},
		"income":  js{},
		"expense": js{},
	}

	for _, m := range stats.Movements {
		if m.Income > 0 {
			data["income"].(js)[stats.Categories[m.Category]] = m.Income
		}
		if m.Expense > 0 {
			data["expense"].(js)[stats.Categories[m.Category]] = m.Expense
		}
	}

	JSON(data)
}
