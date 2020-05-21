package lufthansa

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

// processLimitOffset transforms the general limit and offset parameters,
// which are available for most of the API requests, into an usable string
// for creating the API request URL.
func processLimitOffset(l, o int) string {
	var ret strings.Builder
	if l != 0 {
		ret.WriteString(fmt.Sprintf("limit=%d", l))
	}
	if o != 0 {
		if strings.Contains(ret.String(), "limit") {
			ret.WriteString(fmt.Sprintf("&offset=%d", o))
		} else {
			ret.WriteString(fmt.Sprintf("offset=%d", o))
		}
	}
	return ret.String()
}

// decodeErrors decodes the API response, according to the HTTP response error
// code sent. This functions should be used to decode errors, as they all have the same format.
func decodeErrors(res *http.Response) (interface{}, error) {
	defer res.Body.Close()

	switch res.StatusCode {
	case 400:
		ret := &BadRequestError{}
		err := xml.NewDecoder(res.Body).Decode(ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	case 402, 404, 405:
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
	case 500:
		ret := &jsonAPIError{}
		err := json.NewDecoder(res.Body).Decode(ret)
		if err != nil {
			return nil, err
		}
		return &APIError{
			RetryIndicator: ret.ProcessingErrors.ProcessingError.RetryIndicator,
			Type:           ret.ProcessingErrors.ProcessingError.Type,
			Code:           ret.ProcessingErrors.ProcessingError.Code,
			Description:    ret.ProcessingErrors.ProcessingError.Description,
			InfoURL:        ret.ProcessingErrors.ProcessingError.InfoURL,
		}, nil
	default:
		return nil, nil
	}
}
