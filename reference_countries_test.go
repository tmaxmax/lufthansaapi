package lufthansa_test

import (
	"context"
	"testing"

	lufthansa "github.com/tmaxmax/lufthansaapi"
	"golang.org/x/text/language"
)

func TestAPI_FetchCountries(t *testing.T) {
	cr := api.FetchCountries(&lufthansa.RefParams{Lang: &language.English, Limit: 10})
	for i := 0; cr.Next(ctx) && i < 2; i++ {
		t.Log(cr, "\n", cr.String())
	}
	if cr.Error() != nil {
		t.Fatal(cr.Error())
	}
	cr.Previous(ctx)
	if cr.Error() != nil {
		t.Fatal(cr.Error())
	}
	cr.Next(ctx)
	for i := 0; cr.Previous(ctx) && i < 2; i++ {
		t.Log(cr, "\n", cr.String())
	}
	if cr.Error() != nil {
		t.Fatal(cr.Error())
	}
}

func TestAPI_FetchCountry(t *testing.T) {
	c, err := api.FetchCountry(context.Background(), "RO", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", c)
}
