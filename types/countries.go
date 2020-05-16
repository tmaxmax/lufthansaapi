package types

import "encoding/json"

type countryName struct {
	LanguageCode string `json:"@LanguageCode"`
	Content      string `json:"$"`
}
type countryNames []countryName

func (cns *countryNames) UnmarshalJSON(data []byte) error {
	if data[0] == '[' {
		return json.Unmarshal(data, (*[]countryName)(cns))
	} else if data[0] == '"' {
		var cn countryName
		if err := json.Unmarshal(data, &cn); err != nil {
			return err
		}
		*cns = append(*cns, cn)
	}
	return nil
}

type country struct {
	CountryCode string `json:"CountryCode"`
	Names       struct {
		Name countryNames `json:"Name"`
	} `json:"Names"`
}
type countries []country

func (cs *countries) UnmarshalJSON(data []byte) error {
	if data[0] == '[' {
		return json.Unmarshal(data, (*[]country)(cs))
	} else if data[0] == '"' {
		var c country
		if err := json.Unmarshal(data, &c); err != nil {
			return err
		}
		*cs = append(*cs, c)
	}
	return nil
}

type link struct {
	Href string `json:"@Href"`
	Rel  string `json:"@Rel"`
}
type links []link

func (ls *links) UnmarshalJSON(data []byte) error {
	if data[0] == '[' {
		return json.Unmarshal(data, (*[]link)(ls))
	} else if data[0] == '"' {
		var l link
		if err := json.Unmarshal(data, &l); err != nil {
			return err
		}
		*ls = append(*ls, l)
	}
	return nil
}

type responseMeta struct {
	Version    string `json:"@Version"`
	Link       links  `json:"Link"`
	TotalCount int    `json:"TotalCount"`
}

// CountriesResponse represents the decoded API response returned
// by FetchCountries API method. It is analoguous to the official
// Lufthansa API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Countries
type CountriesResponse struct {
	CountryResources struct {
		Countries struct {
			Country countries `json:"Country"`
		} `json:"Countries"`
	} `json:"CountryResource"`
	Meta responseMeta `json:"Meta"`
}
