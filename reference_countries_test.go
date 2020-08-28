package lufthansa_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"golang.org/x/text/language"

	lufthansa "github.com/tmaxmax/lufthansaapi"
)

var api *lufthansa.API
var ctx = context.Background()

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

func ExampleAPI_FetchCountries() {
	// Fetch 50 countries starting from country 11, with their names in all available languages. Defaults are all languages, 20, 0.
	ref := api.FetchCountries(&lufthansa.RefParams{
		Lang:   nil,
		Limit:  50,
		Offset: 10,
	})
	// Fetching the API is delayed to the first call to the Next method. If Next returns false, an error occurred or there is no data.
	if !ref.Next(ctx) && ref.Error() != nil {
		// Error is idempotent, the value returned changes only after iterations (after calling Next or Previous), no need to save it in a variable.
		log.Fatalln(ref.Error())
	}
	// The Countries struct has a useful string representation, simply output it to console.
	fmt.Printf("%s\n", ref)
	// Save all fetched sets for later. If you want to fetch the rest of the data using a different API client, pass a pointer to it to the Copy method.
	references := []*lufthansa.Countries{ref.Copy(nil)}
	// Loop through the reference.
	for ref.Next(ctx) {
		references = append(references, ref.Copy(nil))
	}
	// Check if the loop ended because of an error
	if ref.Error() != nil {
		log.Fatalln(ref.Error())
	}
	// Output the country code and English name of every country. Use the golang.org/x/text/language package to find the preferred language.
	for _, r := range references {
		for _, c := range r.Countries {
			fmt.Println(c.CountryCode, c.Names[language.English])
		}
	}
	// You can also loop in the opposite direction. After the iterator reaches the end, it keeps a link to the previous, first and last data sets
	// to be able to start iterating again. If you already copied each set of data, use that instead, to prevent making more requests to the API.
	for ref.Previous(ctx) {
		for _, c := range ref.Countries {
			name, ok := c.Names[language.German]
			if !ok {
				continue
			}
			fmt.Println(name)
		}
	}
	if ref.Error() != nil {
		log.Fatalln(ref.Error())
	}
	// Get the last data set. This is not the same with the set retrieved by the last iteration of Next. This set has the last limit countries.
	ref.Last(ctx)
	if ref.Error() != nil {
		log.Fatalln(ref.Error())
	}
	// Get the first data set. This is not the same with the first retrieved set (the initial request). This set has the first limit countries from offset 0.
	ref.First(ctx)
	if ref.Error() != nil {
		log.Fatalln(ref.Error())
	}
	// Calling Previous will invalidate this data set.
	ref.Previous(ctx)
	// Use the HasSelf method to check if the struct holds a valid data set.
	if ref.HasSelf() {
		fmt.Println("valid dataset")
	} else {
		fmt.Println("invalid dataset") // this will be outputted.
	}
}

func ExampleAPI_FetchCountry() {
	country, err := api.FetchCountry(ctx, "US", nil)
	if err != nil {
		log.Fatalln(err)
	}
	// Not passing a language parameter returns the country name in all available languages.
	for lang, name := range country.Names {
		fmt.Println(lang.String(), "-", name)
	}
	country, err = api.FetchCountry(ctx, "US", &language.English)
	if err != nil {
		log.Fatalln(err)
	}
	_, ok := country.Names[language.English]
	fmt.Println(ok) // true
	_, ok = country.Names[language.Afrikaans]
	fmt.Println(ok) // false
	// The struct has a useful string representation, simply output it to console.
	fmt.Printf("%s\n", country)
}
