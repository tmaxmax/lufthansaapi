package typestringer

import (
	"fmt"
	"strconv"
	"time"
)

func ExampleStringifier_Stringify() {
	s := NewStringifier(
		Indent(4),
		StructSeparator(" :: "),
		MapSeparator(": "),
		ChanPrefix("<--"),
		IndexGen(func(i int) string {
			return "Index " + strconv.Itoa(i) + ": "
		}),
	)

	x := 5
	fmt.Println(s.Stringify(x, ""))

	type user struct {
		Name string
		Age  int
	}

	type software struct {
		Name        string
		DateCreated time.Time
		Properties  map[interface{}]interface{}
		Users       []*user
		password    string
	}

	type key struct {
		Value string
		Count int
	}

	var mySoftware = &software{
		Name:        "Typestringer",
		DateCreated: time.Date(2020, 8, 25, 0, 0, 0, 0, time.UTC),
		Properties: map[interface{}]interface{}{
			"easyToUse":       true,
			"importCount":     0,
			key{"sarmale", 1}: "ghimbir",
		},
		Users: []*user{
			{
				Name: "John",
				Age:  25,
			},
			nil,
		},
	}

	stringify := s.Stringify(mySoftware, "")
	fmt.Println(stringify)

	fmt.Print("END")

	// Output: 5
	// Name :: Typestringer
	// DateCreated :: 2020-08-25 00:00:00 +0000 UTC
	// Properties ::
	//     easyToUse: true
	//     importCount: 0
	//     (typestringer.key): ghimbir
	// Users ::
	//     Index 0: Name :: John
	//              Age :: 25
	//     Index 1: nil
	// END
}
