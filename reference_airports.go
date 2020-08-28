package lufthansa

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// FetchAirports requests from the airports reference. Pass parameters as mentioned in
// the API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Airports.
//
// The function returns a pointer to the decoded cities response struct. If this is nil, the function
// will return either an APIError pointer, a GatewayError pointer or an error. If there is an APIError, then
// there is no GatewayError and vice-versa. Check first for errors.
func (a *lufthansa.API) FetchAirports(p lufthansa.RefParams, LHOperated bool) (*AirportsReference, interface{}, error) {
	url := strings.Builder{}
	url.WriteString(fmt.Sprintf("%s/airports/%s", lufthansa.mdsReferenceAPI, p.ToURL()))
	if strings.Contains(url.String(), "?") {
		url.WriteString(fmt.Sprintf("&LHoperated=%t", LHOperated))
	} else {
		url.WriteString(fmt.Sprintf("?LHoperated=%t", LHOperated))
	}

	res, err := a.fetch(url.String())
	if err != nil {
		return nil, nil, err
	}

	switch res.StatusCode {
	case 200:
		ret := &AirportsReference{}
		err = xml.NewDecoder(res.Body).Decode(ret)
		if err != nil {
			return nil, nil, err
		}
		res.Body.Close()
		return ret, nil, nil
	default:
		apiErr, err := lufthansa.decodeErrors(res)
		return nil, apiErr, err
	}
}

type airportPosition struct {
	Latitude  float64 `xml:"Coordinate>Latitude"`
	Longitude float64 `xml:"Coordinate>Longitude"`
}

type airport struct {
	AirportCode  string                    `xml:"AirportCode"`
	Position     airportPosition           `xml:"Position"`
	CityCode     string                    `xml:"CityCode"`
	CountryCode  string                    `xml:"CountryCode"`
	LocationType string                    `xml:"LocationType"`
	Names        []lufthansa.referenceName `xml:"Names>Name"`
	UTCOffset    string                    `xml:"UtcOffset"`
	TimeZoneID   string                    `xml:"TimeZoneId"`
}

type AirportsReference struct {
	Airports []airport        `xml:"Airports>Airport"`
	Meta     []lufthansa.meta `xml:"Meta"`
}
