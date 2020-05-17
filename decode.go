package lufthansa

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

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

// decodeErrors unmarshals the API response, according to the HTTP response error
// code sent. This functions should be used to decode errors, as they all have the same format.
func decodeErrors(res *http.Response) (interface{}, error) {
	defer res.Body.Close()

	switch res.StatusCode {
	case 400, 402, 404, 405:
		ret := &APIError{}
		err := xml.NewDecoder(res.Body).Decode(ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	case 401, 403:
		ret := &GatewayError{}
		err := json.NewDecoder(res.Body).Decode(ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	default:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return string(body), nil
	}
}
