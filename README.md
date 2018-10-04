# n26 - A command-line client for your N26 banking

This tool allows you to perform some of the actions you would on N26's website or mobile app.

There are and will be some limitations:

 * I am developing this from a French account, so some features will not be available because I don't have access to it.
 * There is, of course, **NO WARRANTY** for any consequences that may occur after the use of this tool (for instance, from experience, it is a **very bad idea** to block a card before activating it).

On a more positive note, this is what the tools can handle:

 * View personal and account information
 * View your main account's balance
 * View your spaces' balances and goals
 * View information, state and limits of your cards
 * Transfer money from one of your space to another
 * Transfer money to another N26 user through MoneyBeam
 * Display your past transactions
 * Display your expense and income statistics by category

As of now, those feature will be implemented soon:

 * Block and unblock a specific card
 * Change a card's capabilities and limits
 * Download your statements as PDFs

The following are the feature I am not yet interesting in implementing (mainly because they might be risky):
 * Activating a card (since I do not have spare cards to develop on)
 * Changing a card's PIN
 * Transfer money to another bank's account (IBAN transfer)

Any action that would result in money moving have to be reviewed and confirmed on the command-line. Money transfer to a third-party (MoneyBeam) have to be confirmed from your paired phone as well.

## Authentication

On first launch, your N26 email address and password to initiate a connection, those are not stored, either on your computer or anywhere else. Your credentials are used once to retrieve access and refresh tokens that are used in all requests. As long as the refresh token does not expire, the command-line client will keep on working.

The tokens **are stored** in your home directory (_~/.config/n26.auth_ on Linux, _~/.n26.auth_ on Mac OS) with 0600 permissions, so it is only readable by your user. It is to be noted that those tokens have control over your N26 account and should be protected appropriately.

It appears that only one access token is allowed at the same time for a specific user, so using this tool will end your session on your mobile, and vice-versa. The mobile app will log in again automatically with your fingerprint, and this tool will request another token automatically as well.

## Usage

```
./n26 --help
usage: n26 <command> [<args> ...]

N26 command-line client

Flags:
  -h, --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  info
    Display the account holder personal information

  account
    Display the account information

  balance
    Display the account current balance

  stats [<flags>]
    Get income and expense statistics

  spaces list
    List your spaces and their balances

  spaces transfer <source> <destination> <amount>
    Transfer money from one space to another

  cards list
    Display the cards linked to your account

  cards limits
    Displays the limits for your cards

  transactions list
    List your past transactions

  transactions beam [<flags>] <recipient> <amount>
    Create a Money Beam
```