package lufthansa

// TokenError struct is the object that JSON is decoded in when
// the status response from the request is 401, which means the token
// is invalid or missing.
type TokenError struct {
	Error string `json:"Error"`
}

// APIError struct is the object all API errors XML are decoded in.
// Fields are analoguous to the XML sent. See https://developer.lufthansa.com/docs/read/api_basics/Error_Messages
// for more information.
type APIError struct {
	RetryIndicator bool   `xml:"ProcessingError,RetryIndicator,attr"`
	Type           string `xml:"ProcessingError>Type"`
	Description    string `xml:"ProcessingError>Description"`
	InfoURL        string `xml:"ProcessingError>InfoURL"`
}

type metaLink struct {
	Rel  string `xml:"Rel,attr"`
	Href string `xml:"Href,attr"`
}

type meta struct {
	Version    string     `xml:"Version,attr"`
	Links      []metaLink `xml:"Link"`
	TotalCount int        `xml:"TotalCount"`
}

type countryName struct {
	LanguageCode string `xml:"LanguageCode,attr"`
	Name         string `xml:",chardata"`
}

type country struct {
	CountryCode string        `xml:"CountryCode"`
	Names       []countryName `xml:"Names>Name"`
}

// CountriesResponse represents the decoded response returned
// by FetchCountries API method. It isn't analoguous to the XML
// tags, to keep the structure simple.
// Lufthansa API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Countries
type CountriesResponse struct {
	Countries []country `xml:"Countries>Country"`
	Meta      meta      `xml:"Meta"`
}

type cityName struct {
	Languagecode string `xml:"LanguageCode,attr"`
	Name         string `xml:",chardata"`
}

type city struct {
	CityCode    string     `xml:"CityCode"`
	CountryCode string     `xml:"CountryCode"`
	Names       []cityName `xml:"Names>Name"`
}

// CitiesResponse represents the decoded API response returned
// by FetchCities method. It isn't analoguous to the XML tags, to keep the structure simple.
// Lufthansa API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Cities
type CitiesResponse struct {
	Cities []city `xml:"Cities>City"`
	Meta   []meta `xml:"Meta"`
}
