package lufthansa_test

import (
	"context"
	"log"
	"os"
	"testing"

	"golang.org/x/text/language"

	lufthansa "github.com/tmaxmax/lufthansaapi"
)

var api *lufthansa.API

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var err error

		api, err = lufthansa.NewAPI(context.Background(), os.Getenv("LOA_ID"), os.Getenv("LOA_SECRET"), 5, 1000)
		if err != nil {
			log.Println(err)
			return 1
		}

		return m.Run()
	}())
}

func TestAPI_FetchCountries(t *testing.T) {
	ctx := context.Background()
	cr := api.FetchCountries(&lufthansa.RefParams{
		Limit: 10,
	})
	t.Log("next")
	for cr.Next(ctx) {
		t.Log(cr, "\n", cr.String())
	}
	if err := cr.Error(); err != nil {
		t.Fatal(err)
	}
	t.Log("previous")
	for cr.Previous(ctx) {
		t.Log(cr, "\n", cr.String())
	}
	if err := cr.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestAPI_FetchCountry(t *testing.T) {
	c, err := api.FetchCountry(context.Background(), "RO", &language.English)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", c)
}
