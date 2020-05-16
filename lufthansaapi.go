package lufthansa

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	// APIOauth represents the URL directing to the Oauth token generation
	APIOauth string = "https://api.lufthansa.com/v1/oauth/token"
	// APIFetch is the main URL used to get data from the API
	APIFetch string = "https://api.lufthansa.com/v1"
)

// Token represents the object returned by the Lufthansa Oauth, containing the access token, token type and expiration time.
type Token struct {
	AccessToken    string `json:"access_token"`
	TokenType      string `json:"token_type"`
	ExpiresIn      int    `json:"expires_in"`
	generationTime time.Time
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s", strings.Title(t.TokenType), t.AccessToken)
}

// API represents the main object that you will use to interact with the Lufthansa API.
type API struct {
	clientID     string
	clientSecret string
	token        Token
}

func (a *API) getToken() error {
	payload := strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", a.clientID, a.clientSecret))
	req, err := http.NewRequest("POST", APIOauth, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	t := &Token{}
	err = json.NewDecoder(res.Body).Decode(t)
	if err != nil {
		return err
	}
	t.generationTime = time.Now()

	a.token = *t
	return nil
}

func (a *API) getNewToken() error {
	delta, err := time.ParseDuration(fmt.Sprintf("%ds", a.token.ExpiresIn))
	if err != nil {
		return err
	}
	if time.Now().After(a.token.generationTime.Add(delta)) {
		err = a.getToken()
		return err
	}
	return nil
}

// NewAPI constructs the API object, having as parametres the client's ID and client's secret.
func NewAPI(id, secret string) (*API, error) {
	ret := &API{}
	ret.clientID = id
	ret.clientSecret = secret
	err := ret.getToken()
	if err != nil {
		return &API{}, err
	}
	return ret, nil
}

// fetch function returns the API response from the provided URL as an io.Reader, making it
// easy to decode JSON afterwards
func (a *API) fetch(url string) (io.Reader, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", a.token.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(body)), nil
}
