package lufthansa

import (
	"log"
	"os"
	"testing"
	"time"
)

const (
	testDelay          string = "200ms"
	formatAPIError     string = "API Error\nRetryIndicator: %t\nType: %s\nDescription: %s\nInfoURL: %s\n"
	formatGatewayError string = "Gateway Error: %s\n"
)

type TestCountryURLItem struct {
	RefP   RefParams
	Result string
}

func initializeAPI() *API {
	api, err := NewAPI(os.Getenv("LOA_ID"), os.Getenv("LOA_SECRET"))
	if err != nil {
		log.Fatalf("API initialization failed: %v\n", err)
	}
	return api
}

func apiFetchTestHasError(t *testing.T, err error, i int, apiError *APIError, gatewayError *GatewayError) bool {
	if err != nil {
		t.Errorf("Test %d failed, error: %+v", i, err)
		return true
	} else if apiError != nil {
		t.Logf(formatAPIError, apiError.RetryIndicator, apiError.Type, apiError.Description, apiError.InfoURL)
		return true
	} else if gatewayError != nil {
		t.Logf(formatGatewayError, gatewayError.Error)
		return true
	}
	return false
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
	sleepTime, _ := time.ParseDuration(testDelay)

	testParams := []RefParams{
		{},
		{Lang: "EN"},
		{Code: "DK"},
		{Limit: 5, Lang: "FR"},
		{Code: "ZZ"},
		{Limit: 100, Offset: 50},
		{Limit: 50, Offset: 30, Lang: "HR"},
	}

	for i, p := range testParams {
		t.Logf("\nTest %d...\n", i)

		fetched, apiError, gatewayError, err := api.FetchCountries(p)
		if apiFetchTestHasError(t, err, i, apiError, gatewayError) {
			continue
		}

		for _, c := range fetched.Countries {
			for _, n := range c.Names {
				if (n.LanguageCode == "EN" && p.Lang == "") || n.LanguageCode == p.Lang {
					t.Logf(n.Name)
				} else if p.Lang != "" {
					t.Logf("_")
				}
			}
		}

		time.Sleep(sleepTime)
	}
}

func TestFetchCities(t *testing.T) {
	t.Logf("\nTesting Fetch Cities...\n")

	api := initializeAPI()
	sleepTime, _ := time.ParseDuration(testDelay)

	testParams := []RefParams{
		{},
		{Lang: "EN"},
		{Code: "NYC"},
		{Limit: 5, Lang: "FR"},
		{Code: "ZZ"},
		{Limit: 1000, Offset: 4000},
		{Limit: 9876, Offset: 5432, Lang: "HR"},
	}

	for i, p := range testParams {
		t.Logf("\nTest %d...\n", i)

		fetched, apiError, gatewayError, err := api.FetchCities(p)
		if apiFetchTestHasError(t, err, i, apiError, gatewayError) {
			continue
		}

		for _, c := range fetched.Cities {
			for _, n := range c.Names {
				if (n.LanguageCode == "EN" && p.Lang == "") || n.LanguageCode == p.Lang {
					t.Logf(n.Name)
				} else if p.Lang != "" {
					t.Logf("_")
				}
			}
		}

		time.Sleep(sleepTime)
	}
}
