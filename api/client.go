package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"golang.org/x/oauth2"
)

type N26Client http.Client

type N26Request struct {
	Method  string
	Path    string
	Params  map[string]string
	Body    interface{}
	Decoder *JSON
}

const (
	baseURL = "https://api.tech26.de"
	// baseURL = "http://127.0.0.1:10000"
)

func NewClient() (*N26Client, error) {
	c := oauth2.Config{
		Endpoint:     oauth2.Endpoint{TokenURL: fmt.Sprintf("%s/oauth/token", baseURL)},
		ClientID:     "android",
		ClientSecret: "secret",
	}

	var token *oauth2.Token
	if creds, err := LoadCredentials(); err == nil {
		token = &oauth2.Token{TokenType: creds.TokenType, AccessToken: creds.AccessToken, RefreshToken: creds.RefreshToken, Expiry: creds.Expiry}
		if !token.Valid() {
			token = &oauth2.Token{RefreshToken: creds.RefreshToken}
		}
	} else {
		fmt.Print("N26 email address: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		username := scanner.Text()

		fmt.Print("N26 Password: ")
		password, err := terminal.ReadPassword(syscall.Stdin)
		line()
		if err != nil {
			return nil, fmt.Errorf("could not read password")
		}

		token, err = c.PasswordCredentialsToken(oauth2.NoContext, username, string(password))
		if err != nil {
			return nil, err
		}

		token, _ = c.TokenSource(oauth2.NoContext, token).Token()

		SaveCredentials(token, time.Now().Add(50*time.Minute))
	}

	return (*N26Client)(c.Client(oauth2.NoContext, token)), nil
}

func (cl *N26Client) Request(r *N26Request, retry bool) (interface{}, error) {
	url := fmt.Sprintf("%s%s", baseURL, r.Path)
	if len(r.Params) > 0 {
		url = fmt.Sprintf("%s?%s", url, query(r.Params).Encode())
	}

	var body io.Reader
	if r.Body != nil {
		data, err := json.Marshal(r.Body)
		if err != nil {
			return nil, fmt.Errorf("could not marshal request")
		}
		body = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(r.Method, url, body)
	if err != nil {
		return nil, err
	}

	if r.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := (*http.Client)(cl).Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		ExpireCredentials()

		if !retry {
			cl, err := NewClient()
			if err == nil {
				return cl.Request(r, true)
			}
		}

		return nil, fmt.Errorf("credentials have expired, please try again")
	}

	defer resp.Body.Close()

	if os.Getenv("DEBUG") != "" {
		fmt.Printf("%s %s -> %d\n", r.Method, r.Path, resp.StatusCode)

		b, _ := ioutil.ReadAll(resp.Body)
		resp.Body = ioutil.NopCloser(bytes.NewReader(b))

		fmt.Println(string(b))
	}

	if resp.StatusCode > 399 {
		r.Decoder = NewJSON(new(Error))

		output, err := r.Decoder.Decode(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("an unknown error has occured")
		}
		return nil, fmt.Errorf(output.(*Error).Message)
	}

	if r.Decoder == nil {
		return nil, nil
	}

	output, err := r.Decoder.Decode(resp.Body)

	c := cl.Transport.(*oauth2.Transport)
	newToken, err := c.Source.Token()
	if err == nil {
		SaveCredentials(newToken, newToken.Expiry)
	}

	return output, err
}

func ConfigPath() string {
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("%s/.config/n26.auth", os.Getenv("HOME"))
	case "darwin":
		return fmt.Sprintf("%s/.n26.auth", os.Getenv("HOME"))
	default:
		Fatal(fmt.Errorf("platform '%s' unsupported", runtime.GOOS))
	}
	return ""
}

func SaveCredentials(token *oauth2.Token, exp time.Time) {
	creds := Credentials{
		TokenType:    token.TokenType,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       exp,
	}

	data, err := json.Marshal(creds)
	if err != nil {
		Fatal(fmt.Errorf("could not marshal credentials"))
	}

	err = ioutil.WriteFile(ConfigPath(), data, 0600)
	if err != nil {
		Fatal(fmt.Errorf("could not write credentials file to '%s'", ConfigPath()))
	}
}

func LoadCredentials() (*Credentials, error) {
	data, err := ioutil.ReadFile(ConfigPath())
	if err != nil {
		return nil, err
	}

	creds := new(Credentials)
	err = json.Unmarshal(data, creds)
	if err != nil {
		return nil, err
	}

	return creds, nil
}

func DeleteCredentials() {
	err := os.Remove(ConfigPath())
	if err != nil {
		Fatal(fmt.Errorf("could not delete credentials file at '%s'", ConfigPath()))
	}
}

func ExpireCredentials() {
	creds, err := LoadCredentials()
	if err != nil {
		DeleteCredentials()
	}

	SaveCredentials(&oauth2.Token{
		TokenType:    creds.TokenType,
		AccessToken:  creds.AccessToken,
		RefreshToken: creds.RefreshToken,
	}, time.Unix(0, 0))
}
