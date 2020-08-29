package lufthansa

import (
	"fmt"
	"io"
	"net/http"

	"github.com/tmaxmax/lufthansaapi/internal/util"
)

type (
	//// BadRequestError is the type of error returned on HTTP status response code 400. This error is not documented!
	//BadRequestError struct {
	//	Category string `xml:"category" json:"category"`
	//	Text     string `xml:"text" json:"text"`
	//}

	// GatewayError struct holds the data for access token errors.
	// See https://developer.lufthansa.com/docs/read/api_basics/Error_Messages
	GatewayError struct {
		What string `json:"Error"`
	}
	// APIError struct holds the data for any request processing error.
	// See https://developer.lufthansa.com/docs/read/api_basics/Error_Messages
	APIError struct {
		RetryIndicator bool   `xml:"ProcessingError,RetryIndicator,attr" json:"ProcessingErrors.ProcessingError.@RetryIndicator"`
		Type           string `xml:"ProcessingError>Type" json:"ProcessingErrors.ProcessingError.Type"`
		Code           string `xml:"ProcessingError>Code" json:"ProcessingErrors.ProcessingError.Code"`
		Description    string `xml:"ProcessingError>Description" json:"ProcessingErrors.ProcessingError.Description"`
		InfoURL        string `xml:"ProcessingError>InfoURL" json:"ProcessingErrors.ProcessingError.InfoURL"`
	}
	// unknownError is a placeholder for error types that might not be documented. The struct holds the raw API response.
	unknownError struct {
		response string
	}
)

//func (bre *BadRequestError) Error() string {
//	return fmt.Sprintf("BadRequestError: %s: %s", bre.Category, bre.Text)
//}
//
//func (bre *BadRequestError) String() string {
//	return stringer.Stringify(bre, "")
//}
//
//func (bre *BadRequestError) decode(r io.ReadCloser) error {
//	return util.Decode(r, bre)
//}

func (ge *GatewayError) Error() string {
	return fmt.Sprintf("GatewayError: %s", ge.What)
}

func (ge *GatewayError) String() string {
	return util.Stringer.Stringify(ge, "")
}

func (ge *GatewayError) decode(r io.ReadCloser) error {
	return util.Decode(r, ge)
}

func (ae *APIError) Error() string {
	return fmt.Sprintf("APIError: Code %s, Type %s, Retry %t: %s", ae.Code, ae.Type, ae.RetryIndicator, ae.Description)
}

func (ae *APIError) String() string {
	return util.Stringer.Stringify(ae, "")
}

func (ae *APIError) decode(r io.ReadCloser) error {
	return util.Decode(r, ae)
}

func (ue *unknownError) Error() string {
	return ue.response
}

func (ue *unknownError) decode(r io.ReadCloser) error {
	data, err := util.ReadAll(r)
	if err != nil {
		return err
	}
	ue.response = string(data)
	return nil
}

// decodeErrors decodes the API response, according to the HTTP status code. If the API responded with an error, the
// response body will be closed, no further reading being possible.
func decodeErrors(res *http.Response) error {
	var apiError error

	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		apiError = &GatewayError{}
	case http.StatusBadRequest, http.StatusNotFound, http.StatusMethodNotAllowed:
		apiError = &APIError{}
	case http.StatusInternalServerError:
		apiError = &unknownError{}
	}
	if err := apiError.(apiResponse).decode(res.Body); err != nil {
		return err
	}
	return apiError
}
