package lufthansa

import (
	"context"
	"io"

	"github.com/tmaxmax/lufthansaapi/internal/util"

	"golang.org/x/text/language"
)

type (
	Country struct {
		CountryCode string
		Names       referenceNames
	}
	// Countries represents the decoded API response of the countries reference endpoint.
	// Lufthansa API documentation: https://developer.lufthansa.com/docs/read/api_details/reference_data/Countries
	Countries struct {
		Countries []Country
		meta      meta
		iterator
	}

	countryUnmarshal struct {
		CountryCode string                   `xml:"CountryCode" json:"CountryCode"`
		Names       []referenceNameUnmarshal `xml:"Names>Name" json:"Names.Name"`
	}
	countriesUnmarshal struct {
		Countries []countryUnmarshal `xml:"Countries>Country" json:"CountryResource.Countries.Country"`
		Meta      *metaUnmarshal     `xml:"Meta" json:"CountryResource.Meta"`
	}
)

func (cs *Countries) decode(r io.ReadCloser) error {
	cu := &countriesUnmarshal{}
	if err := util.Decode(r, cu); err != nil {
		return err
	}
	cs.Countries = make([]Country, len(cu.Countries))
	for i := range cs.Countries {
		cs.Countries[i].make(&cu.Countries[i])
	}
	cs.meta.make(cu.Meta)
	return nil
}

func (cs *Countries) metadata() *meta {
	return &cs.meta
}

// Copy creates a copy of the current set, optionally changing the underlying API. Errors that occurred while calling
// Countries.Next are discarded
func (cs *Countries) Copy(newAPI *API) *Countries {
	r := &Countries{
		Countries: make([]Country, 0, len(cs.Countries)),
		iterator:  cs.iterator.copy(newAPI),
	}
	for i := range cs.Countries {
		r.Countries = append(r.Countries, *cs.Countries[i].Copy())
	}
	return r
}

func (cs *Countries) String() string {
	return util.Stringer.Stringify(cs, "")
}

func (c *Country) decode(r io.ReadCloser) error {
	cu := &countriesUnmarshal{}
	if err := util.Decode(r, cu); err != nil {
		return err
	}
	c.make(&cu.Countries[0])
	return nil
}

func (c *Country) make(cu *countryUnmarshal) {
	c.CountryCode = cu.CountryCode
	c.Names.make(cu.Names)
}

func (c *Country) Copy() *Country {
	return &Country{
		CountryCode: c.CountryCode,
		Names:       c.Names.Copy(),
	}
}

func (c *Country) String() string {
	return util.Stringer.Stringify(c, "")
}

// FetchCountries requests from the countries reference. If you want to fetch a single country, use FetchCountry instead.
// The API request doesn't happen here, you must call the Next method before.
func (a *API) FetchCountries(p *RefParams) *Countries {
	c := &Countries{
		meta: meta{
			links: metaLinks{
				metaKeyNext: mdsReferenceAPI + "/countries/" + p.ToURL(),
			},
		},
		iterator: iterator{
			api: a,
		},
	}
	c.iterator.ref = c
	return c
}

// FetchCountry requests a single country, identified by its 2 letter ISO 3166-1 country code.
func (a *API) FetchCountry(ctx context.Context, countryCode string, lang *language.Tag) (*Country, error) {
	p := &RefParams{code: countryCode, Lang: lang}
	fetched, err := a.fetch(ctx, mdsReferenceAPI+"/countries/"+p.ToURL())
	if err != nil {
		return nil, err
	}

	c := &Country{}
	return c, c.decode(fetched)
}
