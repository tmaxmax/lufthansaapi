package lufthansa

import (
	"testing"
)

type TestItem struct {
	f      CountriesParams
	Result string
}

func TestCountryURLs(t *testing.T) {
	tests := []TestItem{
		{CountriesParams{}, ""},
		{CountriesParams{Lang: "EN"}, "?lang=EN"},
		{CountriesParams{CountryCode: "DK"}, "DK"},
		{CountriesParams{Lang: "EN", CountryCode: "DK"}, "DK?lang=EN"},
		{CountriesParams{Limit: 20}, "?limit=20"},
		{CountriesParams{Limit: 20, Offset: 1}, "?limit=20&offset=1"},
		{CountriesParams{Lang: "EN", Limit: 20, Offset: 1}, "?lang=EN&limit=20&offset=1"},
	}

	for i := range tests {
		result := tests[i].f.ToURL()
		if result != tests[i].Result {
			t.Errorf("FAILED -- TEST %d:\n%s !=\n%s\n", i, result, tests[i].Result)
		} else {
			t.Logf("PASSED -- TEST %d", i)
		}
	}
}
