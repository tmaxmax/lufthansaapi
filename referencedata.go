package lufthansa

import (
	"encoding/xml"
	"fmt"
	"strings"
)

const (
	mdsReferenceAPI string = fetchAPI + "/mds-references"
	referenceAPI    string = fetchAPI + "/references"
)

// RefParams is a struct containing the parameters
// used to make requests to any of the Reference APIs
//
// Fields:
//  - Code will differ based on the references you request. See each Fetch function's mentions.
//  - Lang is 2 letter ISO 3166-1 language code, used to tell the API in what
//    language should the country names be sent. If it isn't given a value, the
//    API sends the country names in all available languages.
//  - Limit represents the number of records returned per request. Default is set
//    to 20, maximum is 100 (if a value bigger than 100 is given, 100 will be taken)
//  - Offset represents the number of records skipped (default is 0). For example,
//    if offset is 20 and limit is 100, the response will contain records from
//    country no. 20 to country no. 119 (100 countries).
type RefParams struct {
	Code   string
	Lang   string
	Limit  int
	Offset int
}

// ToURL transforms a CountriesParams struct into an URL usable format,
// so that it can be concatenated to te Request API URL.
func (p RefParams) ToURL() string {
	var ret strings.Builder
	if p.Code != "" {
		ret.WriteString(p.Code)
	}
	if p.Lang != "" {
		ret.WriteString(fmt.Sprintf("?lang=%s", p.Lang))
	}
	if lo := processLimitOffset(p.Limit, p.Offset); lo != "" && strings.Contains(ret.String(), "?") {
		ret.WriteString(fmt.Sprintf("&%s", lo))
	} else if lo != "" {
		ret.WriteString(fmt.Sprintf("?%s", lo))
	}
	return ret.String()
}

// FetchCountries requests from the countries reference. Pass parameters as mentioned in
// the API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Countries.
//
// The function returns a pointer to the decoded cities response struct. If this is nil, the function
// will return either an APIError pointer, a GatewayError pointer or an error. If there is an APIError, then
// there is no GatewayError and vice-versa. Check first for errors.
func (a *API) FetchCountries(p RefParams) (*CountriesReference, interface{}, error) {
	url := fmt.Sprintf("%s/countries/%s", mdsReferenceAPI, p.ToURL())
	res, err := a.fetch(url)
	if err != nil {
		return nil, nil, err
	}

	switch res.StatusCode {
	case 200:
		ret := &CountriesReference{}
		err = xml.NewDecoder(res.Body).Decode(ret)
		if err != nil {
			return nil, nil, err
		}
		res.Body.Close()
		return ret, nil, nil
	default:
		apiErr, err := decodeErrors(res)
		return nil, apiErr, err
	}
}

// FetchCities requests from the cities reference. Pass parameters as mentioned in
// the API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Cities.
//
// The function returns a pointer to the decoded cities response struct. If this is nil, the function
// will return either an APIError pointer, a GatewayError pointer or an error. If there is an APIError, then
// there is no GatewayError and vice-versa. Check first for errors.
func (a *API) FetchCities(p RefParams) (*CitiesReference, interface{}, error) {
	url := fmt.Sprintf("%s/cities/%s", mdsReferenceAPI, p.ToURL())
	res, err := a.fetch(url)
	if err != nil {
		return nil, nil, err
	}

	switch res.StatusCode {
	case 200:
		ret := &CitiesReference{}
		err = xml.NewDecoder(res.Body).Decode(ret)
		if err != nil {
			return nil, nil, err
		}
		res.Body.Close()
		return ret, nil, nil
	default:
		apiErr, err := decodeErrors(res)
		return nil, apiErr, err
	}
}

// FetchAirports requests from the airports reference. Pass parameters as mentioned in
// the API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Airports.
//
// The function returns a pointer to the decoded cities response struct. If this is nil, the function
// will return either an APIError pointer, a GatewayError pointer or an error. If there is an APIError, then
// there is no GatewayError and vice-versa. Check first for errors.
func (a *API) FetchAirports(p RefParams, LHOperated bool) (*AirportsReference, interface{}, error) {
	url := strings.Builder{}
	url.WriteString(fmt.Sprintf("%s/airports/%s", mdsReferenceAPI, p.ToURL()))
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
		apiErr, err := decodeErrors(res)
		return nil, apiErr, err
	}
}

// FetchNearestAirports requests from the airports reference. Pass parameters as mentioned in
// the API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Nearest_Airport.
//
// The function returns a pointer to the decoded cities response struct. If this is nil, the function
// will return either an APIError pointer, a GatewayError pointer or an error. If there is an APIError, then
// there is no GatewayError and vice-versa. Check first for errors.
func (a *API) FetchNearestAirports(lat, long float32, langCode string) (*NearestAirportsReference, interface{}, error) {
	url := fmt.Sprintf("%s/airports/nearest/%.3f,%.3f", referenceAPI, lat, long)
	if langCode != "" {
		url += fmt.Sprintf("?%s", langCode)
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
		apiErr, err := decodeErrors(res)
		return nil, apiErr, err
	}
}

// FetchAirlines requests from the airlines reference. Pass parameters as mentioned in
// the API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Nearest_Airport.
// Note that RefParam's field Lang is ignored!
//
// The function returns a pointer to the decoded cities response struct. If this is nil, the function
// will return either an APIError pointer, a GatewayError pointer or an error. If there is an APIError, then
// there is no GatewayError and vice-versa. Check first for errors.
func (a *API) FetchAirlines(p RefParams) (*AirlinesReference, interface{}, error) {
	if p.Lang != "" {
		p.Lang = ""
	}
	url := fmt.Sprintf("%s/airlines/%s", mdsReferenceAPI, p.ToURL())
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
		apiErr, err := decodeErrors(res)
		return nil, apiErr, err
	}
}

// FetchAircraft requests from the aircraft reference. Pass parameters as mentioned in
// the API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Aircraft.
// Note that RefParam's field Lang is ignored!
//
// The function returns a pointer to the decoded cities response struct. If this is nil, the function
// will return either an APIError pointer, a GatewayError pointer or an error. If there is an APIError, then
// there is no GatewayError and vice-versa. Check first for errors.
func (a *API) FetchAircraft(p RefParams) (*AircraftReference, interface{}, error) {
	if p.Lang != "" {
		p.Lang = ""
	}
	url := fmt.Sprintf("%s/aircraft/%s", mdsReferenceAPI, p.ToURL())
	res, err := a.fetch(url)
	if err != nil {
		return nil, nil, err
	}

	switch res.StatusCode {
	case 200:
		ret := &AircraftReference{}
		err = xml.NewDecoder(res.Body).Decode(ret)
		if err != nil {
			return nil, nil, err
		}
		res.Body.Close()
		return ret, nil, nil
	default:
		apiErr, err := decodeErrors(res)
		return nil, apiErr, err
	}
}
