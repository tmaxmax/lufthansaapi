package lufthansa

import (
	"encoding/xml"
	"fmt"
)

// FetchNearestAirports requests from the airports reference. Pass parameters as mentioned in
// the API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Nearest_Airport.
//
// The function returns a pointer to the decoded cities response struct. If this is nil, the function
// will return either an APIError pointer, a GatewayError pointer or an error. If there is an APIError, then
// there is no GatewayError and vice-versa. Check first for errors.
func (a *lufthansa.API) FetchNearestAirports(lat, long float32, langCode LangCode) (*NearestAirportsReference, interface{}, error) {
	url := fmt.Sprintf("%s/airports/nearest/%.3f,%.3f", lufthansa.referenceAPI, lat, long)
	if langCode != "" {
		url += fmt.Sprintf("?lang=%s", langCode)
	}

	res, err := a.fetch(url)
	if err != nil {
		return nil, nil, err
	}

	switch res.StatusCode {
	case 200:
		ret := &NearestAirportsReference{}
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

type (

	// AirportsReference represents the decoded API response returned
	// by FetchAirports method. It isn't analogous to the XML tags, to keep the structure simple.
	// Lufthansa API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Airports

	nearestAirportDistance struct {
		Value         int    `xml:"Value"`
		UnitOfMeasure string `xml:"UOM"`
	}
	nearestAirport struct {
		AirportCode  string                    `xml:"AirportCode"`
		Position     whatever.airportPosition  `xml:"Position"`
		CityCode     string                    `xml:"CityCode"`
		CountryCode  string                    `xml:"CountryCode"`
		LocationType string                    `xml:"LocationType"`
		Names        []lufthansa.referenceName `xml:"Names>Name"`
		Distance     nearestAirportDistance    `xml:"Distance"`
	}
	// NearestAirportsReference represents the decoded API response returned
	// by FetchNearestAirports method. It isn't analogous to the XML tags, to keep the structure simple.
	// Lufthansa API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Nearest_Airport
	NearestAirportsReference struct {
		Airports []nearestAirport `xml:"Airports>Airport"`
		Meta     []lufthansa.meta `xml:"Meta"`
	}
)
