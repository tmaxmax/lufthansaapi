package lufthansa

import (
	"context"
	"fmt"
	"io"

	"github.com/tmaxmax/lufthansaapi/internal/util"

	"golang.org/x/text/language"
)

type (
	CityUnmarshal struct {
		CityCode    string                   `xml:"CityCode" json:"CityCode"`
		CountryCode string                   `xml:"CountryCode" json:"CountryCode"`
		Names       []referenceNameUnmarshal `xml:"Names>Name" json:"Names.Name"`
	}
	CitiesUnmarshal struct {
		Cities []CityUnmarshal `xml:"Cities>City" json:"CityResource.Cities.City"`
		Meta   *metaUnmarshal  `xml:"Meta" json:"CityResource.Meta"`
	}

	City struct {
		CityCode    string
		CountryCode string
		Names       referenceNames
	}
	Cities struct {
		Cities []City
		meta   meta
		iterator
	}
)

func (cs *Cities) decode(r io.ReadCloser) error {
	cu := &CitiesUnmarshal{}
	if err := util.Decode(r, cu); err != nil {
		return err
	}
	cs.Cities = make([]City, len(cu.Cities))
	for i := range cs.Cities {
		cs.Cities[i].make(&cu.Cities[i])
	}
	cs.meta.make(cu.Meta)
	return nil
}

func (cs *Cities) metadata() *meta {
	return &cs.meta
}

func (cs *Cities) Copy(newAPI *API) *Cities {
	r := &Cities{
		Cities:   make([]City, 0, len(cs.Cities)),
		iterator: cs.iterator.copy(newAPI),
	}
	for i := range cs.Cities {
		r.Cities = append(r.Cities, *cs.Cities[i].Copy())
	}
	return r
}

func (cs *Cities) String() string {
	return util.Stringer.Stringify(cs, "")
}

func (c *City) make(cu *CityUnmarshal) {
	c.CityCode = cu.CityCode
	c.CountryCode = cu.CountryCode
	c.Names.make(cu.Names)
}

func (c *City) decode(r io.ReadCloser) error {
	cu := &CitiesUnmarshal{}
	if err := util.Decode(r, cu); err != nil {
		return err
	}
	c.make(&cu.Cities[0])
	return nil
}

func (c *City) Copy() *City {
	return &City{
		CityCode:    c.CityCode,
		CountryCode: c.CountryCode,
		Names:       c.Names.Copy(),
	}
}

func (c *City) String() string {
	return util.Stringer.Stringify(c, "")
}

func (a *API) FetchCities(p *RefParams) *Cities {
	c := &Cities{
		meta: meta{
			links: metaLinks{
				metaKeyNext: fmt.Sprintf("%s/cities/%s", mdsReferenceAPI, p.ToURL()),
			},
		},
		iterator: iterator{
			api: a,
		},
	}
	c.iterator.ref = c
	return c
}

func (a *API) FetchCity(ctx context.Context, cityCode string, lang *language.Tag) (*City, error) {
	p := &RefParams{code: cityCode, Lang: lang}
	fetched, err := a.fetch(ctx, fmt.Sprintf("%s/cities/%s", mdsReferenceAPI, p.ToURL()))
	if err != nil {
		return nil, err
	}

	c := &City{}
	return c, c.decode(fetched)
}
