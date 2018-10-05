package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/apognu/n26/api"
	"github.com/apognu/n26/cli"
)

func main() {
	kp := kingpin.New("n26", "N26 command-line client")
	kp.HelpFlag.Short('h')
	kp.UsageTemplate(kingpin.DefaultUsageTemplate)

	kpFormat := kp.Flag("format", "how to display data").Short('o').Default("pretty").Enum("pretty", "json")

	kpInfo := kp.Command("info", "Display the account holder personal information")
	kpAccount := kp.Command("account", "Display the account information")
	kpBalance := kp.Command("balance", "Display the account current balance")

	kpStats := kp.Command("stats", "Get income and expense statistics")
	kpStatsFrom := kpStats.Flag("from", "Date to start statistics from (e.g. 2018-01-01)").String()
	kpStatsTo := kpStats.Flag("to", "Date to end statistics at (e.g. 2018-01-31)").String()

	kpSpaces := kp.Command("spaces", "Manage your spaces")
	kpSpacesList := kpSpaces.Command("list", "List your spaces and their balances")

	kpSpacesTransfer := kpSpaces.Command("transfer", "Transfer money from one space to another")
	kpSpacesTransferFrom := kpSpacesTransfer.Arg("source", "ID or name of the source space").Required().String()
	kpSpacesTransferTo := kpSpacesTransfer.Arg("destination", "ID or name of the destination space").Required().String()
	kpSpacesTransferAmount := kpSpacesTransfer.Arg("amount", "amount of money to transfer").Required().Float64()

	kpCards := kp.Command("cards", "Display the cards linked to your account")
	kpCardsList := kpCards.Command("list", "Display the cards linked to your account")
	kpCardLimits := kpCards.Command("limits", "Displays the limits for your cards")

	kpTransactions := kp.Command("transactions", "Manage your transactions")
	kpTransactionsList := kpTransactions.Command("list", "List your past transactions")
	kpTransactionsFrom := kpTransactions.Flag("from", "date from which to list transactions").String()
	kpTransactionsTo := kpTransactions.Flag("to", "date to which to list transactions").String()
	kpTransactionsLimit := kpTransactions.Flag("limit", "number of transactions to display").Short('l').Default("50").Int()

	kpMoneyBeam := kpTransactions.Command("beam", "Create a Money Beam")
	kpMoneyBeamRecipient := kpMoneyBeam.Arg("recipient", "email or phone number of the recipient").Required().String()
	kpMoneyBeamName := kpMoneyBeam.Flag("name", "name of the recipient").Short('n').String()
	kpMoneyBeamAmount := kpMoneyBeam.Arg("amount", "amount to transfer").Required().Float64()
	kpMoneyBeamComment := kpMoneyBeam.Flag("comment", "comment to add to the transfer").Short('c').String()

	args := kingpin.MustParse(kp.Parse(os.Args[1:]))

	cl, err := api.NewClient()
	if err != nil {
		cli.Fatal(fmt.Errorf("could not authenticate to N26"))
	}

	categories, _ := cl.GetCategories()
	meta := &cli.Metadata{
		Categories: categories,
	}

	var cmd cli.Printable

	switch args {
	case kpInfo.FullCommand():
		cmd, err = cl.GetPersonalInformation(meta)
	case kpAccount.FullCommand():
		cmd, err = cl.GetAccount(meta)
	case kpBalance.FullCommand():
		cmd, err = cl.GetBalance(meta)
	case kpStats.FullCommand():
		cmd, err = cl.GetStatistics(meta, *kpStatsFrom, *kpStatsTo)
	case kpCardsList.FullCommand():
		cmd, err = cl.GetCards(meta)
	case kpCardLimits.FullCommand():
		cmd, err = cl.GetLimits(meta)
	case kpTransactionsList.FullCommand():
		cmd, err = cl.GetPastTransactions(meta, *kpTransactionsFrom, *kpTransactionsTo, *kpTransactionsLimit)
	case kpMoneyBeam.FullCommand():
		cmd, err = cl.CreateMoneyBeam(meta, *kpMoneyBeamName, *kpMoneyBeamRecipient, *kpMoneyBeamAmount, *kpMoneyBeamComment)
	case kpSpacesList.FullCommand():
		cmd, err = cl.GetSpaces(meta)
	case kpSpacesTransfer.FullCommand():
		cmd, err = cl.CreateSpaceTransfer(meta, *kpSpacesTransferFrom, *kpSpacesTransferTo, *kpSpacesTransferAmount)
	}

	if err != nil {
		cli.Fatal(err)
		return
	}

	if cmd != nil {
		switch *kpFormat {
		case "pretty":
			cmd.Print(meta)
		case "json":
			cmd.JSON(meta)
		}
	}
}
