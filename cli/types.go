package cli

import (
	"fmt"

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

type Account struct {
	Bank string `json:"bankName"`
	IBAN string `json:"iban"`
	BIC  string `json:"bic"`
}

type Balance struct {
	AvailableBalance float64 `json:"availableBalance"`
	UsageBalance     float64 `json:"usableBalance"`
	Currency         string  `json:"currency"`
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
