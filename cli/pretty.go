package cli

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/pmylund/sortutil"
	"github.com/sirupsen/logrus"
)

func (msg SimpleMessage) Print(meta *Metadata) {
	logrus.Info(msg)
}

func (info PersonalInformation) Print(meta *Metadata) {
	bd := time.Unix(info.BirthDate/1000, 0)

	title("Card holder")

	attr("Name", fmt.Sprintf("%s %s", info.Firstname, info.Lastname))
	attr("E-mail address", info.Email)
	attr("Phone number", info.Phone)
	attr("Birth date", bd.Format("02 Jan 2006"))
	attr("Nationality", info.Nationality)
}

func (account Account) Print(meta *Metadata) {
	title("Account information")

	attr("Bank name", account.Bank)
	attr("IBAN", account.IBAN)
	attr("BIC", account.BIC)
}

func (balance Balance) Print(meta *Metadata) {
	title("Account balance")

	attr("Balance", fmt.Sprintf("%.2f %s", balance.AvailableBalance, balance.Currency))
	if balance.AvailableBalance != balance.UsageBalance {
		attr("Usable balance", fmt.Sprintf("%.2f %s", balance.UsageBalance, balance.Currency))
	}
}

func (cards CardList) Print(meta *Metadata) {
	for _, card := range cards {
		exp := time.Unix(card.Expiration/1000, 0)
		model := card.ProductType
		if card.ProductType != card.Design {
			model = fmt.Sprintf("%s/%s", card.ProductType, card.Design)
		}

		title(fmt.Sprintf("*-%s", card.Number[len(card.Number)-4:]))

		attr("ID", attrColor.Sprintf(card.ID))
		attr("Holder", card.Holder)
		attr("Expires on", exp.Format("Jan 2006"))
		attr("Type", card.Type)
		attr("Model", model)
		if s, ok := CardStatuses[card.Status]; ok {
			attr("Status", s.Color.Sprintf(s.Text))
		} else {
			attr("Status", card.Status)
		}

		line()
	}
}

func (limits LimitList) Print(meta *Metadata) {
	title("Card limits")

	for _, limit := range limits {
		if l, ok := LimitStatuses[limit.Limit]; ok {
			attr(l, fmt.Sprintf("%.2f", limit.Amount))
		} else {
			attr(limit.Limit, fmt.Sprintf("%.2f", limit.Amount))
		}
	}
}

func (transactions PastTransactionList) Print(meta *Metadata) {
	headers := []string{
		"Date",
		"Third-party",
		"Amount",
		"Category",
		"Location",
		"Comment",
	}

	data := make([][]string, len(transactions))
	for idx, trx := range transactions {
		date := time.Unix(trx.Date/1000, 0)

		party := trx.MerchantName
		if trx.Partner != "" {
			party = trx.Partner
		}

		amount := okColor.Sprintf("← %.2f %s", trx.Amount, trx.Currency)
		if trx.Amount < 0 {
			amount = errColor.Sprintf("→ %.2f %s", math.Abs(trx.Amount), trx.Currency)
		}

		if trx.Scheme == "SPACES" {
			party = "N26 Spaces"
		}

		data[idx] = []string{
			titleColor.Sprintf(date.Format("02 Jan 2006 15:04")),
			party,
			amount,
			meta.GetCategory(trx.Category),
			trx.MerchantCity,
			attrColor.Sprintf(trx.Comment),
		}
	}

	table := table()
	table.SetHeader(headers)
	table.AppendBulk(data)

	table.Render()
}

func (spaces Spaces) Print(meta *Metadata) {
	for _, space := range spaces.Spaces {
		if space.Primary {
			title(fmt.Sprintf("%s (PRIMARY)", space.Name))
		} else {
			title(space.Name)
		}

		attr("ID", attrColor.Sprintf(space.ID))
		attr("Amount", fmt.Sprintf("%.2f %s", space.Balance.AvailableBalance, space.Balance.Currency))
		if space.Goal.Amount > 0 {
			progress := space.Balance.AvailableBalance / space.Goal.Amount * 100

			attr("Goal", fmt.Sprintf("%.2f %s", space.Goal.Amount, space.Balance.Currency))
			attr("Progress", fmt.Sprintf("%.1f %%", progress))
		}

		line()
	}
}

func (stats Statistics) Print(meta *Metadata) {
	title("Global movements")
	attr("Period", fmt.Sprintf("%s - %s", time.Unix(stats.From/1000, 0).Format("02 Jan 2006"), time.Unix(stats.To/1000, 0).Format("02 Jan 2006")))
	attr("Income", fmt.Sprintf("%.2f", stats.TotalIncome))
	attr("Expense", fmt.Sprintf("%.2f", stats.TotalExpense))

	line()
	title("Income by category")
	line()

	progressLength := 40

	income := table()
	income.SetHeader([]string{"Category", "Income", "Income %"})
	sortutil.DescByField(stats.Movements, "Income")
	for _, m := range stats.Movements {
		pct := m.Income / stats.TotalIncome * 100
		prog := strings.Repeat("▪", int(pct)/int(100/progressLength))

		income.Append([]string{
			stats.Categories[m.Category],
			fmt.Sprintf("%.2f", m.Income),
			fmt.Sprintf("%.1f %%", pct),
			prog,
		})
	}
	income.Render()

	line()
	title("Expense by category")
	line()

	expense := table()
	expense.SetHeader([]string{"Category", "Expense", "Expense %"})
	sortutil.DescByField(stats.Movements, "Expense")
	for _, m := range stats.Movements {
		pct := m.Expense / stats.TotalExpense * 100
		prog := strings.Repeat("▪", int(pct)/int(100/progressLength))

		expense.Append([]string{
			stats.Categories[m.Category],
			fmt.Sprintf("%.2f", m.Expense),
			fmt.Sprintf("%.1f %%", m.Expense/stats.TotalExpense*100),
			prog,
		})
	}
	expense.Render()
}
