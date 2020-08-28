package lufthansa

import (
	"context"
	"fmt"
	"io"
	"sync"

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
		api       *API
		err       error
		mu        sync.RWMutex
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
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cu := &countriesUnmarshal{}
	if err := decode(r, cu); err != nil {
		cs = nil
		return err
	}
	cs.Countries = make([]Country, len(cu.Countries))
	for i := range cs.Countries {
		cs.Countries[i].CountryCode = cu.Countries[i].CountryCode
		cs.Countries[i].Names.make(cu.Countries[i].Names)
	}
	cs.meta.make(cu.Meta)
	return nil
}

// fetchFromMeta fetches the link associated with the passed rel. Always lock the appropriate guard before calling this
// function!
func (cs *Countries) fetchFromMeta(ctx context.Context, rel string) (*Countries, error) {
	url, ok := cs.meta.Links[rel]
	if !ok {
		return nil, nil
	}
	fetched, err := cs.api.fetch(ctx, url)
	if err != nil {
		return nil, err
	}

	rcs := &Countries{api: cs.api}
	return rcs, rcs.decode(fetched)
}

// iterate iterates using the meta links provided with the passed rel. This function locks, don't lock in the caller function!
func (cs *Countries) iterate(ctx context.Context, rel string) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	c, err := cs.fetchFromMeta(ctx, rel)
	if err != nil {
		cs.err = err
		return false
	}
	if c == nil {
		if rel == "next" || rel == "previous" {
			if rel == "next" {
				rel = "previous"
			} else {
				rel = "next"
			}
			*cs = Countries{
				meta: meta{
					Links: metaLinks{
						"first": cs.meta.Links["first"],
						rel:     cs.meta.Links["self"],
						"last":  cs.meta.Links["last"],
					},
				},
				api: cs.api,
			}
			cs.mu.Lock()
		}
		return false
	}
	*cs = Countries{
		Countries: c.Countries,
		meta:      c.meta,
		api:       c.api,
	}
	cs.mu.Lock()
	return true
}

// Next fetches the next set of countries, overwriting the current one. If an error occurs, the receiver is not overwritten.
// The method whether further iterations can be done.
// If you want to keep each set individually, call Countries.Copy after iterating to copy the newly fetched set.
func (cs *Countries) Next(ctx context.Context) bool {
	return cs.iterate(ctx, "next")
}

// Previous fetches the previous set of countries, overwriting the current one. If an error occurs, the receiver is not overwritten.
// The method returns the receiver, which is nil when there is no preceding set available, or returns nil if an error occurred.
// If you want to keep each set individually, call Countries.Copy after iterating to copy the newly fetched set.
func (cs *Countries) Previous(ctx context.Context) bool {
	return cs.iterate(ctx, "previous")
}

// HasSelf checks if the set can refetch itself from the API. This returns false only if the struct is in a state where
// only Countries.Next and Countries.Previous are valid operations (the struct holds no data). Use this to check if the
// data is available.
func (cs *Countries) HasSelf() bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	_, ok := cs.meta.Links["self"]
	return ok
}

// Self refetches the last countries resource fetched, overwriting the current set.
// If you want a copy of the current Countries struct, use the Countries.Copy method instead.
func (cs *Countries) Self(ctx context.Context) {
	cs.iterate(ctx, "self")
}

// First fetches the first set of countries, overwriting the current one, if available.
func (cs *Countries) First(ctx context.Context) {
	cs.iterate(ctx, "first")
}

// Last fetches the last set of countries, overwriting the current one, if available.
func (cs *Countries) Last(ctx context.Context) {
	cs.iterate(ctx, "last")
}

// Error returns any errors that occurred while calling Countries.Next. It is idempotent, multiple calls after a single
// iteration return the same result.
func (cs *Countries) Error() error {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	return cs.err
}

// Copy creates a copy of the current set, optionally changing the underlying API. Errors that occurred while calling
// Countries.Next are discarded
func (cs *Countries) Copy(newAPI *API) *Countries {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	if newAPI == nil {
		newAPI = cs.api
	}
	r := &Countries{
		Countries: make([]Country, len(cs.Countries)),
		meta:      cs.meta.copy(),
		api:       newAPI,
		err:       nil,
	}
	for i := range cs.Countries {
		r.Countries[i].Names = cs.Countries[i].Names.copy()
		r.Countries[i].CountryCode = cs.Countries[i].CountryCode
	}
	return r
}

func (cs *Countries) String() string {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	return stringifier.Stringify(cs, "")
}

func (c *Country) decode(r io.ReadCloser) error {
	cu := &countriesUnmarshal{}
	if err := decode(r, cu); err != nil {
		c = nil
		return err
	}
	c.Names.make(cu.Countries[0].Names)
	c.CountryCode = cu.Countries[0].CountryCode
	return nil
}

func (c *Country) String() string {
	return stringifier.Stringify(c, "")
}

// FetchCountries requests from the countries reference. If you want to fetch a single country, use FetchCountry instead.
// The API request doesn't happen here, you must call the Next method before.
func (a *API) FetchCountries(p *RefParams) *Countries {
	return &Countries{
		api: a,
		meta: meta{
			Links: map[string]string{
				"next": fmt.Sprintf("%s/countries/%s", mdsReferenceAPI, p.ToURL()),
			},
		},
	}
}

// FetchCountry requests a single country, identified by its 2 letter ISO 3166-1 country code.
func (a *API) FetchCountry(ctx context.Context, countryCode string, lang *language.Tag) (*Country, error) {
	p := &RefParams{code: countryCode, Lang: lang}
	fetched, err := a.fetch(ctx, fmt.Sprintf("%s/countries/%s", mdsReferenceAPI, p.ToURL()))
	if err != nil {
		return nil, err
	}

	c := &Country{}
	return c, c.decode(fetched)
}
