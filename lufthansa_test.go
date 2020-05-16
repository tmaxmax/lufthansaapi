package lufthansa

import (
	"os"
	"testing"

	"github.com/tmaxmax/lufthansaapi/types"
)

type TestCountryURLItem struct {
	f      CountriesParams
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

func TestCountryURLs(t *testing.T) {
	tests := []TestCountryURLItem{
		{CountriesParams{}, ""},
		{CountriesParams{Lang: "EN"}, "?lang=EN"},
		{CountriesParams{CountryCode: "DK"}, "DK"},
		{CountriesParams{Lang: "EN", CountryCode: "DK"}, "DK?lang=EN"},
		{CountriesParams{Limit: 20}, "?limit=20"},
		{CountriesParams{Limit: 20, Offset: 1}, "?limit=20&offset=1"},
		{CountriesParams{Lang: "EN", Limit: 20, Offset: 1}, "?lang=EN&limit=20&offset=1"},
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

	api, err := NewAPI(os.Getenv("LOA_ID"), os.Getenv("LOA_SECRET"))
	if err != nil {
		t.Fatalf("API initialization failed: %v", err)
	}
	fetched, err := api.FetchCountries(CountriesParams{})
	switch fetched.(type) {
	case *types.CountriesResponse:
		cr := fetched.(*types.CountriesResponse)
		t.Logf("%+v", cr)
	case *types.APIError:
		cr := fetched.(*types.APIError)
		t.Logf("%+v", cr)
	case *types.TokenError:
		cr := fetched.(*types.TokenError)
		t.Logf("%+v", cr)
	default:
		cr := fetched.(string)
		t.Errorf("%s", cr)
	}
}
