package lufthansa

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/gabriel-vasile/mimetype"
	tjson "github.com/tmaxmax/json"
)

type (
	aircraft struct {
		AircraftCode     string                    `xml:"AircraftCode" json:"AircraftCode"`
		Name             []lufthansa.referenceName `xml:"Names>Name" json:"Names.Name"`
		AirlineEquipCode string                    `xml:"AirlineEquipCode" json:"AirlineEquipCode"`
	}
	// AircraftReference represents the decoded API response of the Aircraft Reference endpoint
	// Lufthansa API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Aircraft
	AircraftReference struct {
		AircraftSummary []aircraft       `xml:"AircraftSummaries>AircraftSummary"`
		Meta            []lufthansa.meta `xml:"Meta"`
	}
)

func (ar *AircraftReference) decode(r io.ReadCloser) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if err = r.Close(); err != nil {
		return err
	}
	switch mimetype.Detect(data).String() {
	case "application/xml":
		return xml.Unmarshal(data, ar)
	case "application/json":
		return tjson.Unmarshal(data, ar)
	}
	return lufthansa.ErrUnsupportedFormat
}

func (ar *AircraftReference) String() string {
	return lufthansa.stringifier.Stringify(ar, "")
}

// FetchAircraft requests from the aircraft reference. Here RefParams.Code is the aircraft code.
// Pass parameters as mentioned in the API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Aircraft.
// Note that RefParam's field Lang is ignored!
func (a *lufthansa.API) FetchAircraft(ctx context.Context, p *lufthansa.RefParams) (*AircraftReference, error) {
	if p.Lang != "" {
		p.Lang = ""
	}
	res, err := a.fetch(ctx, fmt.Sprintf("%s/aircraft/%s", lufthansa.mdsReferenceAPI, p.ToURL()))
	if err != nil {
		return nil, err
	}

}
