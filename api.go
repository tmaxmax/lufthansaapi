package lufthansa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/tmaxmax/lufthansaapi/internal/util"

	"golang.org/x/time/rate"

	"github.com/tmaxmax/lufthansaapi/pkg/ratelimithttp"
)

const (
	fetchAPI string = "https://api.lufthansa.com/v1"
	oauthAPI        = fetchAPI + "/oauth/token"
)

type apiResponse interface {
	// decode decodes into the struct the data from the passed response body. If the response Content-Type
	// isn't supported by the implementation, it returns ErrUnsupportedFormat.
	// Every decode implementation shall close the reader1
	decode(io.ReadCloser) error
}

type expiresIn struct {
	time.Duration
}

func (e *expiresIn) UnmarshalJSON(data []byte) error {
	var t int
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	e.Duration = time.Second * time.Duration(t)
	return nil
}

// token represents the object returned by the Lufthansa Oauth,
// containing the access token, token type and expiration time.
// It also holds a generationTime timestamp, that is used for
// obtaining a new access token, when the old one expired.
type token struct {
	AccessToken    string    `json:"access_token"`
	TokenType      string    `json:"token_type"`
	ExpiresIn      expiresIn `json:"expires_in"`
	generationTime time.Time
}

func (t *token) decode(r io.ReadCloser) error {
	return util.Decode(r, t)
}

func (t *token) String() string {
	if t == nil {
		return ""
	}
	return fmt.Sprintf("%s %s", strings.Title(t.TokenType), t.AccessToken)
}

// API represents the main object that you will use to interact with the Lufthansa API. Initialize it with the
// NewAPI function. The struct can't be copied and shall be used multiple times. It is safe for concurrent use.
// Do not create a new API struct per HTTP request, if your application is a HTTP server, use the same struct globally!
// Not doing so will mess rate management and authentication, leading to undesired errors!
type API struct {
	client       *ratelimithttp.Client
	clientID     string
	clientSecret string
	token        *token
	tokenMu      sync.RWMutex
	addr         *API
}

//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func (a *API) copyCheck() {
	if a.addr == nil {
		a.addr = (*API)(noescape(unsafe.Pointer(a)))
	} else if a.addr != a {
		panic("lufthansaapi: API: illegal copy of API value")
	}
}

func (a *API) setToken(ctx context.Context) error {
	a.copyCheck()
	a.tokenMu.Lock()
	defer a.tokenMu.Unlock()

	var err error

	requestURL := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", url.QueryEscape(a.clientID), url.QueryEscape(a.clientSecret))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, oauthAPI, strings.NewReader(requestURL))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	a.token = &token{
		generationTime: time.Now(),
	}
	res, err := a.client.Do(req)
	if err != nil {
		return err
	}

	return a.token.decode(res.Body)
}

func (a *API) refreshToken(ctx context.Context) error {
	a.tokenMu.RLock()
	if time.Since(a.token.generationTime) >= a.token.ExpiresIn.Duration {
		a.tokenMu.RUnlock()
		return a.setToken(ctx)
	}
	a.tokenMu.RUnlock()
	return nil
}

// fetch function returns the API response from the provided URL as an io.ReadCloser. The caller goroutine shall close the reader.
func (a *API) fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	a.copyCheck()
	if err := a.refreshToken(ctx); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept", "application/xml")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Authorization", a.token.String())
	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	if err = decodeErrors(res); err != nil {
		return nil, err
	}
	return res.Body, nil
}

// NewAPI constructs the API object, having as parametres the client's ID and client's secret.
func NewAPI(ctx context.Context, id, secret string, reqPerSecond, reqPerHour int) (*API, error) {
	ret := &API{
		client: ratelimithttp.NewClient(
			&http.Client{
				Timeout: time.Second * 15,
			},
			rate.NewLimiter(rate.Every(time.Second), reqPerSecond),
			rate.NewLimiter(rate.Every(time.Hour), reqPerHour),
		),
		clientID:     id,
		clientSecret: secret,
	}
	if err := ret.setToken(ctx); err != nil {
		return nil, err
	}
	ret.addr = ret
	return ret, nil
}
