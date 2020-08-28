package lufthansa

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type (
	city struct {
		CityCode    string                    `xml:"CityCode"`
		CountryCode string                    `xml:"CountryCode"`
		Names       []lufthansa.referenceName `xml:"Names>Name"`
	}
	// CityReference represents the decoded API response returned
	// by FetchCities method. It isn't analogous to the XML tags, to keep the structure simple.
	// Lufthansa API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Cities
	CityReference struct {
		City []city           `xml:"Cities>City"`
		Meta []lufthansa.meta `xml:"Meta"`
	}
)

func (cr *CityReference) decode(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	switch http.DetectContentType(data) {
	case "application/xml":
		return xml.Unmarshal(data, cr)
	}
	return lufthansa.ErrUnsupportedFormat
}

// FetchCities requests from the cities reference. Pass parameters as mentioned in
// the API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Cities.
func (a *lufthansa.API) FetchCities(ctx context.Context, p *lufthansa.RefParams) (*CityReference, error) {
	var err error

	fetched, err := a.fetch(ctx, fmt.Sprintf("%s/cities/%s", lufthansa.mdsReferenceAPI, p.ToURL()))
	if err != nil {
		return nil, err
	}
	defer lufthansa.checkClose(fetched, err)

	ret := &CityReference{}
	err = ret.decode(fetched)
	return ret, err
}
