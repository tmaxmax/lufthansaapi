package lufthansa_test

import (
	"context"
	"log"
	"os"
	"testing"

	lufthansa "github.com/tmaxmax/lufthansaapi"
)

var (
	api *lufthansa.API
	ctx = context.Background()
)

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
