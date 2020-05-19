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
		{},
		{RefP: RefParams{Lang: "EN"}, Result: "?lang=EN"},
		{RefP: RefParams{Code: "DK"}, Result: "DK"},
		{RefP: RefParams{Lang: "EN", Code: "DK"}, Result: "DK?lang=EN"},
		{RefP: RefParams{Limit: 20}, Result: "?limit=20"},
		{RefP: RefParams{Limit: 20, Offset: 1}, Result: "?limit=20&offset=1"},
		{RefP: RefParams{Lang: "EN", Limit: 20, Offset: 1}, Result: "?lang=EN&limit=20&offset=1"},
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

func TestAPI_FetchCountries(t *testing.T) {
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

func TestAPI_FetchCities(t *testing.T) {
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

type TestFetchAirplanesItem struct {
	Ref  RefParams
	LHop bool
}

func TestAPI_FetchAirports(t *testing.T) {
	t.Logf("\nTesting Fetch Airports...\n")

	api := initializeAPI()
	sleepTime, _ := time.ParseDuration(testDelay)

	testParams := []TestFetchAirplanesItem{
		{},
		{Ref: RefParams{Lang: "EN"}, LHop: true},
		{Ref: RefParams{Limit: 40, Offset: 100}},
		{Ref: RefParams{Limit: 11640}},
		{Ref: RefParams{Offset: 11700}},
		{Ref: RefParams{Lang: "DE", Limit: 200, Offset: 400}, LHop: true},
	}

	for i, p := range testParams {
		t.Logf("\nTest %d...\n", i)

		fetched, apiError, gatewayError, err := api.FetchAirports(p.Ref, p.LHop)
		if apiFetchTestHasError(t, err, i, apiError, gatewayError) {
			continue
		}

		for _, a := range fetched.Airports {
			var j int
			for ; j < len(a.Names) && a.Names[j].LanguageCode != p.Ref.Lang; j++ {
				continue
			}
			if j == len(a.Names) {
				for j = 0; j < len(a.Names) && a.Names[j].LanguageCode != "EN"; j++ {
				}
			}
			t.Logf("\nAirport code: %s\nAirport name: %s (%s)\nAirport position: lat %f, long %f\n\n", a.AirportCode, a.Names[j].Name, a.Names[j].LanguageCode, a.Position.Latitude, a.Position.Longitude)
		}

		time.Sleep(sleepTime)
	}
}
