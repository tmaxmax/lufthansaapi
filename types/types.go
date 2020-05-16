// Package types contains all the types that the
// API response JSON's are paresed in.
package types

// TokenError struct is the object that JSON is decoded in when
// the status response from the request is 401, which means the token
// is invalid or missing.
type TokenError struct {
	Error string `json:"Error"`
}

// APIError struct is the object all API errors JSON are decoded in.
// Fields are analoguous the the JSON sent. See https://developer.lufthansa.com/docs/read/api_basics/Error_Messages
// for more information.
type APIError struct {
	ProcessingErrors struct {
		ProcessingError struct {
			RetryIndicator bool   `json:"@RetryIndicator"`
			Type           string `json:"Type"`
			Code           string `json:"Code"`
			Description    string `json:"Description"`
			InfoURL        string `json:"InfoURL"`
		} `json:"ProcessingError"`
	} `json:"ProcessingErrors"`
}
