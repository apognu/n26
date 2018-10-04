package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"syscall"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
)

var (
	titleColor = color.New(color.Bold, color.FgBlue)
	attrColor  = color.New(color.Faint)

	okColor   = color.New(color.FgGreen)
	warnColor = color.New(color.FgYellow)
	errColor  = color.New(color.FgRed)
)

func Fatal(err error) {
	if e1, ok := err.(*url.Error); ok {
		if e2, ok := e1.Err.(*oauth2.RetrieveError); ok {
			if e2.Response.StatusCode == 401 {
				// DeleteCredentials()
			}

			msg := make(map[string]string)
			json.Unmarshal(e2.Body, &msg)

			logrus.Fatal(msg["error_description"])
			return
		}
	}

	logrus.Fatal(err)
}

func ReadLine(prompt string) string {
	fmt.Print(fmt.Sprintf("%s ", prompt))

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}

func ReadSecret(prompt string) (string, error) {
	fmt.Print(prompt)
	secret, err := terminal.ReadPassword(syscall.Stdin)

	line()

	if err != nil {
		return "", fmt.Errorf("could not read PIN")
	}

	return string(secret), nil
}

func line() {
	fmt.Println()
}

func title(title string) {
	titleColor.Printf("%s:\n", title)
}

func table() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetColumnSeparator("  ")
	table.SetCenterSeparator("")
	table.SetRowSeparator("")
	table.SetAutoWrapText(false)

	return table
}

func attr(key, value string) {
	fmt.Print("  ")
	fmt.Printf("%s %s\n", attrColor.Sprintf("%s:", key), value)
}

func Curr(amount float64, currency string) string {
	return fmt.Sprintf("%.2f %s", amount, currency)
}

func ConfirmSpaceTransfer(from, to *Space, amount float64) {
	title("Please confirm you want to perform the following transfer")
	line()

	data := make([][]string, 3)

	data[0] = []string{errColor.Sprintf(from.Name), "→", Curr(amount, from.Balance.Currency), "→", okColor.Sprintf(to.Name)}
	data[1] = []string{attrColor.Sprintf(from.ID), "", "", "", attrColor.Sprintf(to.ID)}
	data[2] = []string{
		Curr(from.Balance.AvailableBalance, from.Balance.Currency),
		"", "", "",
		Curr(to.Balance.AvailableBalance, to.Balance.Currency),
	}

	table := table()
	table.AppendBulk(data)
	table.Render()

	line()

	if ReadLine("Are you sure you want to perform the transfer? (y/N) ") != "y" {
		Fatal(fmt.Errorf("the transfer was not performed"))
	}
}

func ConfirmMoneyBeam(trx MoneyBeamDetails, balance *Balance) {
	title("Please confirm you want to perform the following transfer")
	fmt.Println("You will be asked for your PIN and will have to confirm the transfer from your paired device.")
	line()

	partnerID := trx.PartnerEmail
	if trx.PartnerPhone != "" {
		partnerID = trx.PartnerPhone
	}

	data := make([][]string, 3)
	data[0] = []string{errColor.Sprintf("Main Account"), "→", Curr(trx.Amount, balance.Currency), "→", okColor.Sprintf(trx.PartnerName)}
	data[1] = []string{Curr(balance.AvailableBalance, balance.Currency), "", trx.Comment, "", attrColor.Sprintf(partnerID)}

	table := table()
	table.AppendBulk(data)
	table.Render()

	line()

	if ReadLine("Are you sure you want to perform the transfer? (y/N) ") != "y" {
		Fatal(fmt.Errorf("the transfer was not performed"))
	}
}
