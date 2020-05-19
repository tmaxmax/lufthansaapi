package lufthansa

type (
	// GatewayError struct is the object that JSON is decoded in when
	// the status response from the request is 401, which means the token
	// is invalid or missing.
	GatewayError struct {
		Error string `json:"Error"`
	}
	// APIError struct is the object all API errors XML are decoded in.
	// Fields are analogous to the XML sent. See https://developer.lufthansa.com/docs/read/api_basics/Error_Messages
	// for more information.
	APIError struct {
		RetryIndicator bool   `xml:"ProcessingError,RetryIndicator,attr"`
		Type           string `xml:"ProcessingError>Type"`
		Description    string `xml:"ProcessingError>Description"`
		InfoURL        string `xml:"ProcessingError>InfoURL"`
	}
)

type referenceName struct {
	LanguageCode string `xml:"LanguageCode,attr"`
	Name         string `xml:",chardata"`
}

// reference meta types
type (
	metaLink struct {
		Rel  string `xml:"Rel,attr"`
		Href string `xml:"Href,attr"`
	}
	meta struct {
		Version    string     `xml:"Version,attr"`
		Links      []metaLink `xml:"Link"`
		TotalCount int        `xml:"TotalCount"`
	}
)

type (
	country struct {
		CountryCode string          `xml:"CountryCode"`
		Names       []referenceName `xml:"Names>Name"`
	}
	// CountriesReference represents the decoded response returned
	// by FetchCountries API method. It isn't analogous to the XML
	// tags, to keep the structure simple.
	// Lufthansa API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Countries
	CountriesReference struct {
		Countries []country `xml:"Countries>Country"`
		Meta      meta      `xml:"Meta"`
	}
)

type (
	city struct {
		CityCode    string          `xml:"CityCode"`
		CountryCode string          `xml:"CountryCode"`
		Names       []referenceName `xml:"Names>Name"`
	}
	// CitiesReference represents the decoded API response returned
	// by FetchCities method. It isn't analogous to the XML tags, to keep the structure simple.
	// Lufthansa API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Cities
	CitiesReference struct {
		Cities []city `xml:"Cities>City"`
		Meta   []meta `xml:"Meta"`
	}
)

type (
	airportPosition struct {
		Latitude  float64 `xml:"Coordinate>Latitude"`
		Longitude float64 `xml:"Coordinate>Longitude"`
	}
	airport struct {
		AirportCode  string          `xml:"AirportCode"`
		Position     airportPosition `xml:"Position"`
		CityCode     string          `xml:"CityCode"`
		CountryCode  string          `xml:"CountryCode"`
		LocationType string          `xml:"LocationType"`
		Names        []referenceName `xml:"Names>Name"`
		UTCOffset    string          `xml:"UtcOffset"`
		TimeZoneID   string          `xml:"TimeZoneId"`
	}
	// AirportsReference represents the decoded API response returned
	// by FetchAirports method. It isn't analogous to the XML tags, to keep the structure simple.
	// Lufthansa API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Airports
	AirportsReference struct {
		Airports []airport `xml:"Airports>Airport"`
		Meta     []meta    `xml:"Meta"`
	}
)
