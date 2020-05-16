package types

type countryName struct {
	LanguageCode string `json:"@LanguageCode"`
	Content      string `json:"$"`
}

type country struct {
	CountryCode string `json:"CountryCode"`
	Names       struct {
		NameRaw []countryName `json:"Name"`
	} `json:"Names"`
}

type responseMeta struct {
	Version string `json:"@Version"`
	Link    []struct {
		Href string `json:"@Href"`
		Rel  string `json:"@Rel"`
	} `json:"Link"`
	TotalCount int `json:"TotalCount"`
}

type CountriesResponse struct {
	CountryResources struct {
		Countries struct {
			Country []country `json:"Country"`
		} `json:"Countries"`
	} `json:"CountryResource"`
	Meta responseMeta `json:"Meta"`
}
