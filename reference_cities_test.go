package lufthansa_test

import (
	"testing"

	lufthansa "github.com/tmaxmax/lufthansaapi"
	"golang.org/x/text/language"
)

func TestAPI_FetchCities(t *testing.T) {
	cr := api.FetchCities(&lufthansa.RefParams{
		Limit: 10,
		Lang:  &language.English,
	})

	for i := 0; cr.Next(ctx) && i < 2; i++ {
		t.Logf("%s", cr)
	}
	if cr.Error() != nil {
		t.Fatal(cr.Error())
	}
	cr.Last(ctx)
	if cr.Error() != nil {
		t.Fatal(cr.Error())
	}
	cr.Next(ctx)
	for i := 0; cr.Previous(ctx) && i < 2; i++ {
		t.Logf("%s", cr)
	}
	if cr.Error() != nil {
		t.Fatal(cr.Error())
	}
}

func TestAPI_FetchCity(t *testing.T) {
	city, err := api.FetchCity(ctx, "BUH", &language.Russian)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", city)
}
