package lufthansa

import (
	"context"
	"io"
	"strings"

	"github.com/tmaxmax/lufthansaapi/internal/util"
	"golang.org/x/text/language"
)

type (
	airportPosition struct {
		Latitude  float64 `xml:"Coordinate>Latitude" json:"Coordinate.Latitude"`
		Longitude float64 `xml:"Coordinate>Longitude" json:"Coordinate.Longitude"`
	}
	airportUnmarshal struct {
		AirportCode  string                   `xml:"AirportCode" json:"AirportCode"`
		Position     airportPosition          `xml:"Position" json:"Position"`
		CityCode     string                   `xml:"CityCode" json:"CityCode"`
		CountryCode  string                   `xml:"CountryCode" json:"CountryCode"`
		LocationType string                   `xml:"LocationType" json:"LocationType"`
		Names        []referenceNameUnmarshal `xml:"Names>Name" json:"Names.Name"`
		UTCOffset    string                   `xml:"UtcOffset" json:"UtcOffset"`
		TimeZoneID   string                   `xml:"TimeZoneId" json:"TimeZoneId"`
	}
	airportsUnmarshal struct {
		Airports []airportUnmarshal `xml:"Airports>Airport" json:"AirportResource.Airports.Airport"`
		Meta     metaUnmarshal      `xml:"Meta" json:"AirportResource.Meta"`
	}

	Airport struct {
		AirportCode  string
		Position     airportPosition
		CityCode     string
		CountryCode  string
		LocationType string
		Names        referenceNames
		UTCOffset    string
		TimeZoneID   string
	}
	Airports struct {
		Airports []Airport
		meta     meta
		iterator
	}
)

func (as *Airports) decode(r io.ReadCloser) error {
	au := &airportsUnmarshal{}
	if err := util.Decode(r, au); err != nil {
		return err
	}
	as.Airports = make([]Airport, len(au.Airports))
	for i := range as.Airports {
		as.Airports[i].make(&au.Airports[i])
	}
	as.meta.make(&au.Meta)
	return nil
}

func (as *Airports) metadata() *meta {
	return &as.meta
}

func (as *Airports) Copy(newAPI *API) *Airports {
	r := &Airports{
		Airports: make([]Airport, 0, len(as.Airports)),
		iterator: as.iterator.copy(newAPI),
	}
	for i := range as.Airports {
		r.Airports = append(r.Airports, *as.Airports[i].Copy())
	}
	return r
}

func (as *Airports) String() string {
	return util.Stringer.Stringify(as, "")
}

func (a *Airport) decode(r io.ReadCloser) error {
	au := &airportsUnmarshal{}
	if err := util.Decode(r, au); err != nil {
		return err
	}
	a.make(&au.Airports[0])
	return nil
}

func (a *Airport) make(au *airportUnmarshal) {
	a.AirportCode = au.AirportCode
	a.Position = au.Position
	a.CityCode = au.CityCode
	a.CountryCode = au.CountryCode
	a.LocationType = au.LocationType
	a.Names.make(au.Names)
	a.UTCOffset = au.UTCOffset
	a.TimeZoneID = au.TimeZoneID
}

func (a *Airport) Copy() *Airport {
	return &Airport{
		AirportCode:  a.AirportCode,
		Position:     a.Position,
		CityCode:     a.CityCode,
		CountryCode:  a.CountryCode,
		LocationType: a.LocationType,
		Names:        a.Names.Copy(),
		UTCOffset:    a.UTCOffset,
		TimeZoneID:   a.TimeZoneID,
	}
}

func (a *Airport) String() string {
	return util.Stringer.Stringify(a, "")
}

func (a *API) FetchAirports(p *RefParams, LHOperated bool) *Airports {
	url := mdsReferenceAPI + "/airports/" + p.ToURL()
	if LHOperated {
		if strings.Contains(url, "?") {
			url += "&LHoperated=1"
		} else {
			url += "?LHoperated=1"
		}
	}
	as := &Airports{
		meta: meta{
			links: metaLinks{
				metaKeyNext: url,
			},
		},
		iterator: iterator{
			api: a,
		},
	}
	as.iterator.ref = as
	return as
}

func (a *API) FetchAirport(ctx context.Context, airportCode string, lang *language.Tag) (*Airport, error) {
	p := &RefParams{code: airportCode, Lang: lang}
	fetched, err := a.fetch(ctx, mdsReferenceAPI+"/airports/"+p.ToURL())
	if err != nil {
		return nil, err
	}

	r := &Airport{}
	return r, r.decode(fetched)
}
