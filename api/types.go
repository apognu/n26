package api

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/pmylund/sortutil"

	"github.com/sirupsen/logrus"

	"github.com/fatih/color"
)

type Metadata struct {
	Categories map[string]string
}

func (meta *Metadata) GetCategories() map[string]string {
	if meta != nil {
		return meta.Categories
	}
	return nil
}

func (meta *Metadata) GetCategory(id string) string {
	if meta != nil {
		if title, ok := meta.Categories[id]; ok {
			return title
		}
	}
	return ""
}

type Printable interface {
	Print(meta *Metadata)
}

type SimpleMessage string

func (msg SimpleMessage) Print(meta *Metadata) {
	logrus.Info(msg)
}

type Error struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type CredentialsExpiry time.Time

type Credentials struct {
	TokenType    string    `json:"token_type"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

type PersonalInformation struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	Title           string `json:"title"`
	Firstname       string `json:"firstName"`
	Lastname        string `json:"lastName"`
	BirthDate       int64  `json:"birthDate"`
	Nationality     string `json:"nationality"`
	SignupCompleted bool   `json:"signupCompleted"`
	Phone           string `json:"mobilePhoneNumber"`
}

func (info PersonalInformation) Print(meta *Metadata) {
	bd := time.Unix(info.BirthDate/1000, 0)

	title("Card holder")

	attr("Name", fmt.Sprintf("%s %s", info.Firstname, info.Lastname))
	attr("E-mail address", info.Email)
	attr("Phone number", info.Phone)
	attr("Birth date", bd.Format("02 Jan 2006"))
	attr("Nationality", info.Nationality)
	if info.SignupCompleted {
		attr("Signup status", okColor.Sprintf("COMPLETE"))
	} else {
		attr("Signup status", errColor.Sprintf("INCOMPLETE"))
	}
}

type Account struct {
	Bank string `json:"bankName"`
	IBAN string `json:"iban"`
	BIC  string `json:"bic"`
}

func (account Account) Print(meta *Metadata) {
	title("Account information")

	attr("Bank name", account.Bank)
	attr("IBAN", account.IBAN)
	attr("BIC", account.BIC)
}

type Balance struct {
	AvailableBalance float64 `json:"availableBalance"`
	UsageBalance     float64 `json:"usableBalance"`
	Currency         string  `json:"currency"`
}

func (balance Balance) Print(meta *Metadata) {
	title("Account balance")

	attr("Balance", fmt.Sprintf("%.2f %s", balance.AvailableBalance, balance.Currency))
	if balance.AvailableBalance != balance.UsageBalance {
		attr("Usable balance", fmt.Sprintf("%.2f %s", balance.UsageBalance, balance.Currency))
	}
}

type CardList []Card

type Card struct {
	ID          string `json:"id"`
	Holder      string `json:"usernameOnCard"`
	Number      string `json:"maskedPan"`
	Expiration  int64  `json:"expirationDate"`
	Type        string `json:"cardType"`
	ProductType string `json:"cardProductType"`
	Design      string `json:"design"`
	Status      string `json:"status"`
}

type CardStatus struct {
	Color *color.Color
	Text  string
}

var (
	CardStatuses = map[string]CardStatus{
		"M_ACTIVE": {
			Color: okColor,
			Text:  "ACTIVE",
		},
		"M_LINKED": {
			Color: warnColor,
			Text:  "LINKED",
		},
		"M_DISABLED": {
			Color: errColor,
			Text:  "BLOCKED",
		},
		"M_PHYSICAL_UNCONFIRMED_DISABLED": {
			Color: errColor,
			Text:  "UNCONFIRMED",
		},
	}
)

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

type LimitList []Limit

type Limit struct {
	Limit  string  `json:"limit"`
	Amount float64 `json:"amount"`
}

var (
	LimitStatuses = map[string]string{
		"POS_DAILY_ACCOUNT": "Payment",
		"ATM_DAILY_ACCOUNT": "Withdrawal",
	}
)

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

type SpaceTransaction struct {
	Amount      float64 `json:"amount"`
	FromSpaceID string  `json:"fromSpaceId"`
	ToSpaceID   string  `json:"toSpaceId"`
}

type PastTransactionList []PastTransaction

type PastTransaction struct {
	ID              string  `json:"id"`
	Type            string  `json:"type"`
	Date            int64   `json:"visibleTS"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currencyCode"`
	Partner         string  `json:"partnerName"`
	Pending         bool    `json:"pending"`
	MerchantName    string  `json:"merchantName"`
	MerchantCity    string  `json:"merchantCity"`
	MerchantCountry string  `json:"merchantCountry"`
	Comment         string  `json:"referenceText"`
	Category        string  `json:"category"`
	Scheme          string  `json:"paymentScheme"`
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

		location := ""
		if trx.MerchantCity != "" {
			location = trx.MerchantCity
		}
		if trx.MerchantCountry != "" {
			location = fmt.Sprintf("%s (%s)", location, trx.MerchantCountry)
		}

		if trx.Scheme == "SPACES" {
			party = "N26 Spaces"
		}

		data[idx] = []string{
			titleColor.Sprintf(date.Format("02 Jan 2006 15:04")),
			party,
			amount,
			meta.GetCategory(trx.Category),
			location,
			attrColor.Sprintf(trx.Comment),
		}
	}

	table := table()
	table.SetHeader(headers)
	table.AppendBulk(data)

	table.Render()
}

type MoneyBeam struct {
	PIN         string           `json:"pin"`
	Transaction MoneyBeamDetails `json:"transaction"`
}

type MoneyBeamDetails struct {
	Type         string  `json:"type"`
	Amount       float64 `json:"amount"`
	PartnerName  string  `json:"partnerName"`
	PartnerEmail string  `json:"partnerEmail,omitempty"`
	PartnerPhone string  `json:"partnerPhone,omitempty"`
	Comment      string  `json:"referenceText,omitempty"`
}

type MoneyBeamPartner struct {
	Name  string
	Email string
	Phone string
}

type Spaces struct {
	Balance float64 `json:"totalBalance"`
	Spaces  []Space `json:"spaces"`
}

type Space struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Primary bool   `json:"isPrimary"`
	Balance struct {
		AvailableBalance float64 `json:"availableBalance"`
		Currency         string  `json:"currency"`
	} `json:"balance"`
	Goal struct {
		Amount float64 `json:"amount"`
	}
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

type ContactList []Contact

type ContactRequest struct {
	Phone string `json:"mobilePhoneNumber"`
	Email string `json:"email"`
}

type Contact struct {
	ID string `json:"id"`
}

func (contacts ContactList) Print(meta *Metadata) {
	fmt.Printf("%#v\n", contacts)
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Statistics struct {
	Categories   map[string]string `json:"-"`
	Currency     string            `json:"-"`
	From         int64             `json:"from"`
	To           int64             `json:"to"`
	TotalExpense float64           `json:"totalExpense"`
	TotalIncome  float64           `json:"totalIncome"`
	Movements    []struct {
		Category string  `json:"id"`
		Expense  float64 `json:"expense"`
		Income   float64 `json:"income"`
	} `json:"items"`
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
