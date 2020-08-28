package lufthansa

import (
	"encoding/xml"
	"fmt"
)

// FetchAirlines requests from the airlines reference. Pass parameters as mentioned in
// the API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Nearest_Airport.
// Note that RefParam's field Lang is ignored!
//
// The function returns a pointer to the decoded cities response struct. If this is nil, the function
// will return either an APIError pointer, a GatewayError pointer or an error. If there is an APIError, then
// there is no GatewayError and vice-versa. Check first for errors.
func (a *lufthansa.API) FetchAirlines(p lufthansa.RefParams) (*AirlinesReference, interface{}, error) {
	if p.Lang != "" {
		p.Lang = ""
	}
	url := fmt.Sprintf("%s/airlines/%s", lufthansa.mdsReferenceAPI, p.ToURL())
	res, err := a.fetch(url)
	if err != nil {
		return nil, nil, err
	}

	switch res.StatusCode {
	case 200:
		ret := &AirlinesReference{}
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
	airline struct {
		IATA  string                    `xml:"AirlineID"`
		ICAO  string                    `xml:"AirlineID_ICAO"`
		Names []lufthansa.referenceName `xml:"Names>Name"`
	}
	// NearestAirportsReference represents the decoded API response returned
	// by FetchAirlines method. It isn't analogous to the XML tags, to keep the structure simple.
	// Lufthansa API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Airlines
	AirlinesReference struct {
		Airlines []airline        `xml:"Airlines>Airline"`
		Meta     []lufthansa.meta `xml:"Meta"`
	}
)
