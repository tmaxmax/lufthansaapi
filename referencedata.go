package lufthansa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/tmaxmax/lufthansaapi/types"
)

const referenceAPI string = FetchAPI + "/mds-references"

// processLimitOffset transforms the general limit and offset parameters,
// which are available for most of the API requests, into an usable string
// for creating the API request URL.
func processLimitOffset(l, o int) string {
	var ret string
	if l != 0 {
		ret += fmt.Sprintf("limit=%d", l)
	}
	if o != 0 {
		if strings.Contains(ret, "limit") {
			ret += fmt.Sprintf("&offset=%d", o)
		} else {
			ret += fmt.Sprintf("offset=%d", o)
		}
	}
	return ret
}

/*CountriesParams is a struct containing the Fetch Countries parameters
used to request from the API the Countries reference.

Fields:
 - CountryCode is a 2-letter ISO 3166-1 country code, used when you only
   want to request information about a single country, represented by the
   code.
 - Lang is 2 letter ISO 3166-1 language code, used to tell the API in what
   language should the country names be sent. If it isn't given a value, the
   API sends the country names in all available languages.
 - Limit represents the number of records returned per request. Default is set
   to 20, maximum is 100 (if a value bigger than 100 is given, 100 will be taken)
 - Offset represents the number of records skipped (default is 0). For example,
   if offset is 20 and limit is 100, the response will contain records from
   country no. 20 to country no. 119 (100 countries).
*/
type CountriesParams struct {
	CountryCode string
	Lang        string
	Limit       int
	Offset      int
}

// ToURL transforms a CountriesParams struct into an URL usable format,
// so that it can be concatenated to te Request API URL.
func (p CountriesParams) ToURL() string {
	var ret string
	if p.CountryCode != "" {
		ret += p.CountryCode
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

// FetchCountries makes an API request for the Countries reference, using the
// passed FetchCountriesParams struct fields. For more information about the
// request parameters, check the FetchCountriesParams struct documentation.
//
// The function returns a pointer to the fetched JSON and an error. The error is
// always nil when the pointer returned isn't null. You have to type assert the
// interface to use the result.
func (a *API) FetchCountries(p CountriesParams) (interface{}, error) {
	url := fmt.Sprintf("%s/countries/%s", referenceAPI, p.ToURL())
	res, err := a.fetch(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 200:
		ret := &types.CountriesResponse{}
		err = json.Unmarshal(body, ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	case 400, 402, 403, 404, 405:
		ret := &types.APIError{}
		err = json.Unmarshal(body, ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	case 401:
		ret := &types.TokenError{}
		err = json.Unmarshal(body, ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	default:
		return string(body) + fmt.Sprintf(" -- %d", res.StatusCode), err
	}
}
