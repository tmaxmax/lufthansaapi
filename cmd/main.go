package main

import (
	"context"
	"fmt"
	"log"
	"os"

	lufthansa "github.com/tmaxmax/lufthansaapi"
	"golang.org/x/text/language"
)

func main() {
	ctx := context.Background()
	api, err := lufthansa.NewAPI(ctx, os.Getenv("LOA_ID"), os.Getenv("LOA_SECRET"), 5, 1000)
	if err != nil {
		log.Fatalln(err)
	}
	countries := api.FetchCountries(&lufthansa.RefParams{
		Lang:   &language.English,
		Limit:  20,
		Offset: 0,
	})
	countriesSlice := []*lufthansa.Countries{}
	for countries.Next(ctx) {
		countriesSlice = append(countriesSlice, countries.Copy(nil))
	}
	if countries.Error() != nil {
		log.Fatalln(err)
	}
	for _, c := range countriesSlice {
		fmt.Println(c.String())
	}
}
