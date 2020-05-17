package lufthansa

import (
	"log"
	"os"
	"testing"
)

type TestCountryURLItem struct {
	f      RefParams
	Result string
}

// userAPI is the format an API access information is written in my JSON file
type userAPI struct {
	APIName      string `json:"name"`
	ClientID     string `json:"id"`
	ClientSecret string `json:"secret"`
}

// MyAPIs is the format my accessible userAPIs credentials are passed
type userAPIs struct {
	APIs []API `json:"apis"`
}

func initializeAPI() *API {
	api, err := NewAPI(os.Getenv("LOA_ID"), os.Getenv("LOA_SECRET"))
	if err != nil {
		log.Fatalf("API initialization failed: %v\n", err)
	}
	return api
}

func TestCountryURLs(t *testing.T) {
	tests := []TestCountryURLItem{
		{RefParams{}, ""},
		{RefParams{Lang: "EN"}, "?lang=EN"},
		{RefParams{Code: "DK"}, "DK"},
		{RefParams{Lang: "EN", Code: "DK"}, "DK?lang=EN"},
		{RefParams{Limit: 20}, "?limit=20"},
		{RefParams{Limit: 20, Offset: 1}, "?limit=20&offset=1"},
		{RefParams{Lang: "EN", Limit: 20, Offset: 1}, "?lang=EN&limit=20&offset=1"},
	}

	t.Logf("Testing Country URLs...\n")
	for i := range tests {
		result := tests[i].f.ToURL()
		if result != tests[i].Result {
			t.Errorf("FAILED -- TEST %d:\n%s !=\n%s\n", i, result, tests[i].Result)
		} else {
			t.Logf("PASSED -- TEST %d", i)
		}
	}
}

func TestFetchCountries(t *testing.T) {
	t.Logf("Testing Fetch Countries...\n")

	api := initializeAPI()

	fetched, err := api.FetchCountries(RefParams{Limit: 5, Offset: 0, Lang: "EN"})
	if err != nil {
		t.Errorf("Failed to fetch countries... error: %v", err)
		return
	}

	switch fetched.(type) {
	case *CountriesResponse:
		cr := fetched.(*CountriesResponse)
		t.Logf("%+v\n", cr)
	case *APIError:
		cr := fetched.(*APIError)
		t.Logf("%+v\n", cr)
	case *TokenError:
		cr := fetched.(*TokenError)
		t.Logf("%+v\n", cr)
	default:
		cr := fetched.(string)
		t.Errorf("%s\n", cr)
	}
}
