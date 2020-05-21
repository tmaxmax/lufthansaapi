package lufthansa

import (
	"log"
	"os"
	"testing"
	"time"
)

const (
	testDelay             string = "200ms"
	formatAPIError        string = "\nAPI Error\nRetryIndicator: %t\nType: %s\nDescription: %s\nInfoURL: %s\n"
	formatGatewayError    string = "\nGateway Error: %s\n"
	formatBadRequestError string = "\nBad Request Error\nCategory: %s\nWhat: %s\n"
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

func apiFetchTestHasError(t *testing.T, err error, i int, aerr interface{}) bool {
	if err != nil {
		t.Errorf("Test %d failed, error: %+v", i, err)
		return true
	}
	switch v := aerr.(type) {
	case *APIError:
		t.Logf(formatAPIError, v.RetryIndicator, v.Type, v.Description, v.InfoURL)
	case *GatewayError:
		t.Logf(formatGatewayError, v.Error)
	case *BadRequestError:
		t.Logf(formatBadRequestError, v.Category, v.Text)
	default:
		return false
	}
	return true
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

		fetched, aerr, err := api.FetchCountries(p)
		if apiFetchTestHasError(t, err, i, aerr) {
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

		fetched, aerr, err := api.FetchCities(p)
		if apiFetchTestHasError(t, err, i, aerr) {
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

		fetched, aerr, err := api.FetchAirports(p.Ref, p.LHop)
		if apiFetchTestHasError(t, err, i, aerr) {
			continue
		}

		for _, a := range fetched.Airports {
			var j int
			for ; j < len(a.Names) && a.Names[j].LanguageCode != p.Ref.Lang; j++ {
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

type TestFetchNearestAirportsItem struct {
	Lat, Long float32
	LangCode  ReferenceLangCode
}

func TestAPI_FetchNearestAirports(t *testing.T) {
	t.Logf("\nTesting Fetch Nearest Airports...\n")

	api := initializeAPI()
	sleepTime, _ := time.ParseDuration(testDelay)

	testParams := []TestFetchNearestAirportsItem{
		{Lat: 51.507351, Long: -0.127758, LangCode: "EN"},
		{Lat: 51.540, Long: 5.935, LangCode: "DE"},
		{Lat: 36.562, Long: -115.596},
		{Lat: 44.435581, Long: 26.102221},
	}

	for i, p := range testParams {
		t.Logf("\nTest %d...\n", i)

		fetched, aerr, err := api.FetchNearestAirports(p.Lat, p.Long, p.LangCode)
		if apiFetchTestHasError(t, err, i, aerr) {
			continue
		}

		for _, na := range fetched.Airports {
			var j int
			for ; j < len(na.Names) && na.Names[j].LanguageCode != p.LangCode; j++ {
			}
			if j == len(na.Names) {
				j = 0
			}
			t.Logf(`
Airport code: %s
Airport name: %s (%s)
Airport position: lat %f, long %f
Airport distance: %d%s

`, na.AirportCode, na.Names[j].Name, na.Names[j].LanguageCode, na.Position.Latitude, na.Position.Longitude, na.Distance.Value, na.Distance.UnitOfMeasure)
		}

		time.Sleep(sleepTime)
	}
}

func TestAPI_FetchAirlines(t *testing.T) {
	t.Logf("\nTesting Fetch Airlines...\n")

	api := initializeAPI()
	sleepTime, _ := time.ParseDuration(testDelay)

	testParams := []RefParams{
		{},
		{Code: "TXL"},
		{Limit: 50},
		{Limit: 30, Offset: 200},
		{Offset: 1150},
	}

	for i, p := range testParams {
		t.Logf("\nTest %d...\n", i)

		fetched, aerr, err := api.FetchAirlines(p)
		if apiFetchTestHasError(t, err, i, aerr) {
			continue
		}

		for _, a := range fetched.Airlines {
			t.Logf(`
Airline IDs: IATA %s, ICAO %s
Airline Name: %s (%s)

`, a.IATA, a.ICAO, a.Names[0].Name, a.Names[0].LanguageCode)
		}

		time.Sleep(sleepTime)
	}
}

func TestAPI_FetchAircraft(t *testing.T) {
	t.Logf("\nTesting Fetch Aircraft...\n")

	api := initializeAPI()
	sleepTime, _ := time.ParseDuration(testDelay)

	testParams := []RefParams{
		{},
		{Code: "70M"},
		{Code: "sarmale"},
		{Limit: 50},
		{Limit: 100, Offset: 200},
		{Offset: 370},
	}

	for i, p := range testParams {
		fetched, aerr, err := api.FetchAircraft(p)
		if apiFetchTestHasError(t, err, i, aerr) {
			continue
		}

		for _, a := range fetched.AircraftSummaries {
			t.Logf(`
Aircraft Code: %s
Aircraft Name: %s (%s)
Airline Equipment Code: %s

`, a.AircraftCode, a.Names[0].Name, a.Names[0].LanguageCode, a.AirlineEquipCode)
		}

		time.Sleep(sleepTime)
	}
}
