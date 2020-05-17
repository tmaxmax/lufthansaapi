package lufthansa

import (
	"log"
	"os"
	"testing"
	"time"
)

const (
	testdelay          string = "200ms"
	formatAPIError     string = "API Error\nRetryIndicator: %t\nType: %s\nDescription: %s\nInfoURL: %s\n"
	formatGatewayError string = "Gateway Error: %s\n"
)

type TestCountryURLItem struct {
	RefP   RefParams
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
		result := tests[i].RefP.ToURL()
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
	sleeptime, _ := time.ParseDuration(testdelay)

	testparams := []RefParams{
		{},
		{Lang: "EN"},
		{Code: "DK"},
		{Limit: 5, Lang: "FR"},
		{Code: "ZZ"},
		{Limit: 100, Offset: 50},
		{Limit: 50, Offset: 30, Lang: "HR"},
	}

	for i, p := range testparams {
		t.Logf("\nTest %d...\n", i)

		fetched, err := api.FetchCountries(p)
		if err != nil {
			t.Errorf("Test %d failed, error: %+v", i, err)
		}

		switch val := fetched.(type) {
		case *CountriesResponse:
			for _, c := range val.Countries {
				for _, n := range c.Names {
					if (n.LanguageCode == "EN" && p.Lang == "") || n.LanguageCode == p.Lang {
						t.Logf(n.Name)
					} else if p.Lang != "" {
						t.Logf("_")
					}
				}
			}
		case *APIError:
			t.Logf(formatAPIError, val.RetryIndicator, val.Type, val.Description, val.InfoURL)
		case *GatewayError:
			t.Logf(formatGatewayError, val.Error)
		default:
			t.Errorf("%+v", val)
		}

		time.Sleep(sleeptime)
	}
}

func TestFetchCities(t *testing.T) {
	t.Logf("\nTesting Fetch Cities...\n")

	api := initializeAPI()
	sleeptime, _ := time.ParseDuration(testdelay)

	testparams := []RefParams{
		{},
		{Lang: "EN"},
		{Code: "NYC"},
		{Limit: 5, Lang: "FR"},
		{Code: "ZZ"},
		{Limit: 1000, Offset: 4000},
		{Limit: 9876, Offset: 5432, Lang: "HR"},
	}

	for i, p := range testparams {
		t.Logf("\nTest %d...\n", i)

		fetched, err := api.FetchCities(p)
		if err != nil {
			t.Errorf("Test %d failed, error: %+v", i, err)
		}

		switch val := fetched.(type) {
		case *CitiesResponse:
			for _, c := range val.Cities {
				for _, n := range c.Names {
					if (n.LanguageCode == "EN" && p.Lang == "") || n.LanguageCode == p.Lang {
						t.Logf(n.Name)
					} else if p.Lang != "" {
						t.Logf("_")
					}
				}
			}
		case *APIError:
			t.Logf(formatAPIError, val.RetryIndicator, val.Type, val.Description, val.InfoURL)
		case *GatewayError:
			t.Logf(formatGatewayError, val.Error)
		default:
			t.Errorf("%+v", val)
		}

		time.Sleep(sleeptime)
	}
}
