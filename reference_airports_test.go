package lufthansa_test

import (
	"testing"

	lufthansa "github.com/tmaxmax/lufthansaapi"
	"golang.org/x/text/language"
)

func TestAPI_FetchAirports(t *testing.T) {
	ar := api.FetchAirports(&lufthansa.RefParams{Lang: &language.English, Limit: 10}, true)
	for i := 0; ar.Next(ctx) && i < 2; i++ {
		t.Logf("%s", ar)
	}
	if ar.Error() != nil {
		t.Fatal(ar.Error())
	}
	ar.Last(ctx)
	if ar.Error() != nil {
		t.Fatal(ar.Error())
	}
	ar.Next(ctx)
	for i := 0; ar.Previous(ctx) && i < 2; i++ {
		t.Logf("%s", ar)
	}
	if ar.Error() != nil {
		t.Fatal(ar.Error())
	}
}

func TestAPI_FetchAirport(t *testing.T) {
	airport, err := api.FetchAirport(ctx, "TXL", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", airport)
}
