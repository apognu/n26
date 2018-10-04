package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func Fatal(err error) {
	if e1, ok := err.(*url.Error); ok {
		if e2, ok := e1.Err.(*oauth2.RetrieveError); ok {
			if e2.Response.StatusCode == 401 {
				DeleteCredentials()
			}

			msg := make(map[string]string)
			json.Unmarshal(e2.Body, &msg)

			logrus.Fatal(msg["error_description"])
			return
		}
	}

	logrus.Fatal(err)
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

func curr(amount float64, currency string) string {
	return fmt.Sprintf("%.2f %s", amount, currency)
}

func confirmSpaceTransfer(from, to *Space, amount float64) {
	title("Please confirm you want to perform the following transfer")
	line()

	data := make([][]string, 3)

	data[0] = []string{errColor.Sprintf(from.Name), "→", curr(amount, from.Balance.Currency), "→", okColor.Sprintf(to.Name)}
	data[1] = []string{attrColor.Sprintf(from.ID), "", "", "", attrColor.Sprintf(to.ID)}
	data[2] = []string{
		curr(from.Balance.AvailableBalance, from.Balance.Currency),
		"", "", "",
		curr(to.Balance.AvailableBalance, to.Balance.Currency),
	}

	table := table()
	table.AppendBulk(data)
	table.Render()

	line()

	if readLine("Are you sure you want to perform the transfer? (y/N) ") != "y" {
		Fatal(fmt.Errorf("the transfer was not performed"))
	}
}

func confirmMoneyBeam(trx MoneyBeamDetails, balance *Balance) {
	title("Please confirm you want to perform the following transfer")
	fmt.Println("You will be asked for your PIN and will have to confirm the transfer from your paired device.")
	line()

	partnerID := trx.PartnerEmail
	if trx.PartnerPhone != "" {
		partnerID = trx.PartnerPhone
	}

	data := make([][]string, 3)
	data[0] = []string{errColor.Sprintf("Main Account"), "→", curr(trx.Amount, balance.Currency), "→", okColor.Sprintf(trx.PartnerName)}
	data[1] = []string{curr(balance.AvailableBalance, balance.Currency), "", trx.Comment, "", attrColor.Sprintf(partnerID)}

	table := table()
	table.AppendBulk(data)
	table.Render()

	line()

	if readLine("Are you sure you want to perform the transfer? (y/N) ") != "y" {
		Fatal(fmt.Errorf("the transfer was not performed"))
	}
}
