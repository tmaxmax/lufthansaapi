package lufthansa

import (
	"encoding/xml"
	"fmt"
	"strings"
)

const (
	referenceAPI string = fetchAPI + "/mds-references"
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
	var ret string
	if p.Code != "" {
		ret += p.Code
	}
	if p.Lang != "" {
		ret += "?lang=" + p.Lang
	}
	if lo := processLimitOffset(p.Limit, p.Offset); lo != "" && strings.Contains(ret, "?") {
		ret += "&" + lo
	} else if lo != "" {
		ret += "?" + lo
	}
	return ret
}

// FetchCountries requests from the countries reference. Pass parameters as mentioned in
// the API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Countries.
//
// The function returns a pointer to the decoded cities response struct. If this is nil, the function
// will return either an APIError pointer, a GatewayError pointer or an error. If there is an APIError, then
// there is no GatewayError and vice-versa. Check first for errors.
func (a *API) FetchCountries(p RefParams) (*CountriesResponse, *APIError, *GatewayError, error) {
	url := fmt.Sprintf("%s/countries/%s", referenceAPI, p.ToURL())
	res, err := a.fetch(url)
	if err != nil {
		return nil, nil, nil, err
	}

	switch res.StatusCode {
	case 200:
		ret := &CountriesResponse{}
		err = xml.NewDecoder(res.Body).Decode(ret)
		if err != nil {
			return nil, nil, nil, err
		}
		res.Body.Close()
		return ret, nil, nil, nil
	default:
		apiError, gatewayError, err := decodeErrors(res)
		return nil, apiError, gatewayError, err
	}
}

// FetchCities requests from the cities reference. Pass parameters as mentioned in
// the API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Cities.
//
// The function returns a pointer to the decoded cities response struct. If this is nil, the function
// will return either an APIError pointer, a GatewayError pointer or an error. If there is an APIError, then
// there is no GatewayError and vice-versa. Check first for errors.
func (a *API) FetchCities(p RefParams) (*CitiesResponse, *APIError, *GatewayError, error) {
	url := fmt.Sprintf("%s/cities/%s", referenceAPI, p.ToURL())
	res, err := a.fetch(url)
	if err != nil {
		return nil, nil, nil, err
	}

	switch res.StatusCode {
	case 200:
		ret := &CitiesResponse{}
		err = xml.NewDecoder(res.Body).Decode(ret)
		if err != nil {
			return nil, nil, nil, err
		}
		res.Body.Close()
		return ret, nil, nil, nil
	default:
		apiError, gatewayError, err := decodeErrors(res)
		return nil, apiError, gatewayError, err
	}
}
