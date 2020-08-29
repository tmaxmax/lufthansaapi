package lufthansa

import (
	"context"
	"errors"
	"io"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/text/language"
)

const (
	mdsReferenceAPI = fetchAPI + "/mds-references"
	referenceAPI    = fetchAPI + "/references"

	metaVersion = "1.0.0"
)

var (
	metaKeySelf     = metaKey{"self"}
	metaKeyFirst    = metaKey{"first"}
	metaKeyLast     = metaKey{"last"}
	metaKeyNext     = metaKey{"next"}
	metaKeyPrevious = metaKey{"previous"}
	metaKeyRelated  = metaKey{"related"}

	ErrInvalidMeta    = errors.New("lufthansa: reference: new meta version detected")
	errMissingMetaKey = errors.New("lufthansa: iterator: fetchFromMeta: missing meta key")
)

type referenceAPIResponse interface {
	apiResponse
	metadata() *meta
}

type metaKey struct {
	string
}

func (mk *metaKey) UnmarshalText(data []byte) error {
	mk.string = string(data)
	return nil
}

func (mk *metaKey) opposite() metaKey {
	switch *mk {
	case metaKeySelf:
		return metaKeySelf
	case metaKeyFirst:
		return metaKeyLast
	case metaKeyLast:
		return metaKeyFirst
	case metaKeyNext:
		return metaKeyPrevious
	case metaKeyPrevious:
		return metaKeyNext
	case metaKeyRelated:
		return metaKeyRelated
	}
	return metaKey{}
}

type (
	referenceNames         map[language.Tag]string
	referenceNameUnmarshal struct {
		LanguageCode language.Tag `xml:"LanguageCode,attr" json:"@LanguageCode"`
		Name         string       `xml:",chardata" json:"$"`
	}
)

func (rn *referenceNames) make(arr []referenceNameUnmarshal) {
	*rn = make(referenceNames)
	for i := range arr {
		(*rn)[arr[i].LanguageCode] = arr[i].Name
	}
}

func (rn *referenceNames) Copy() referenceNames {
	r := make(referenceNames)
	for k, v := range *rn {
		r[k] = v
	}
	return r
}

// reference meta types
type (
	metaLinks         map[metaKey]string
	MetaLinkUnmarshal struct {
		Rel  metaKey `xml:"Rel,attr" json:"@Rel"`
		Href string  `xml:"Href,attr" json:"@Href"`
	}
	meta struct {
		version    string
		links      metaLinks
		totalCount int
	}
	metaUnmarshal struct {
		Version    string              `xml:"Version,attr" json:"@Version"`
		Links      []MetaLinkUnmarshal `xml:"Link" json:"Link"`
		TotalCount int                 `xml:"TotalCount" json:"TotalCount"`
	}
)

func (ml *metaLinks) make(arr []MetaLinkUnmarshal) {
	*ml = make(metaLinks)
	for i := range arr {
		(*ml)[arr[i].Rel] = arr[i].Href
	}
}

func (ml *metaLinks) copy() metaLinks {
	r := make(metaLinks)
	for k, v := range *ml {
		r[k] = v
	}
	return r
}

// checkVersion panics if the meta version isn't the same as the version this API implementation is based on
func (m *meta) checkVersion() bool {
	return m.version == metaVersion
}

func (m *meta) make(mu *metaUnmarshal) {
	m.totalCount = mu.TotalCount
	m.version = mu.Version
	m.links.make(mu.Links)
}

func (m *meta) copy() meta {
	return meta{
		version:    m.version,
		links:      m.links.copy(),
		totalCount: m.totalCount,
	}
}

func (m *meta) hasSelf() bool {
	_, ok := m.links[metaKeySelf]
	return ok && m.checkVersion()
}

func (m *meta) fetch(a *API, ctx context.Context, rel metaKey) (bool, io.ReadCloser, error) {
	if !m.checkVersion() {
		return false, nil, ErrInvalidMeta
	}
	url, ok := m.links[rel]
	if !ok {
		return false, nil, nil
	}
	fetched, err := a.fetch(ctx, url)
	if err != nil {
		return false, nil, err
	}
	return true, fetched, nil
}

type iterator struct {
	ref       referenceAPIResponse
	api       *API
	err       error
	mu        sync.RWMutex
	firstNext sync.Once
}

func (i *iterator) fetchFromMeta(ctx context.Context, rel metaKey) error {
	ok, fetched, err := i.ref.metadata().fetch(i.api, ctx, rel)
	if !ok {
		if err != nil {
			return err
		}
		return errMissingMetaKey
	}
	return i.ref.decode(fetched)
}

func (i *iterator) iterate(ctx context.Context, rel metaKey) bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.err = i.fetchFromMeta(ctx, rel)
	if i.err == errMissingMetaKey {
		if rel == metaKeyNext || rel == metaKeyPrevious {
			*i = iterator{
				ref: i.ref,
				api: i.api,
			}
			i.mu.Lock()
			m := i.ref.metadata()
			*m = meta{
				version: i.ref.metadata().version,
				links: metaLinks{
					metaKeyFirst:   m.links[metaKeyFirst],
					rel.opposite(): m.links[metaKeySelf],
					metaKeyLast:    m.links[metaKeyLast],
				},
				totalCount: m.totalCount,
			}
		}
		return false
	}
	if i.err != nil {
		return false
	}
	return true
}

// Next fetches the next set of countries, overwriting the current one. If an error occurs, the receiver is not overwritten.
// The method whether further iterations can be done.
// If you want to keep each set individually, call Countries.Copy after iterating to copy the newly fetched set.
func (i *iterator) Next(ctx context.Context) bool {
	iterated := false
	i.firstNext.Do(func() {
		i.mu.Lock()
		defer i.mu.Unlock()

		fetched, err := i.api.fetch(ctx, i.ref.metadata().links[metaKeyNext])
		if err != nil {
			i.err = err
			return
		}
		if err = i.ref.decode(fetched); err != nil {
			i.err = err
			return
		}
		iterated = true
	})
	if !iterated && i.err == nil {
		iterated = i.iterate(ctx, metaKeyNext)
	}
	return iterated
}

// Previous fetches the previous set of countries, overwriting the current one. If an error occurs, the receiver is not overwritten.
// The method returns the receiver, which is nil when there is no preceding set available, or returns nil if an error occurred.
// If you want to keep each set individually, call Countries.Copy after iterating to copy the newly fetched set.
func (i *iterator) Previous(ctx context.Context) bool {
	return i.iterate(ctx, metaKeyPrevious)
}

// HasSelf checks if the set can refetch itself from the API. This returns false only if the struct is in a state where
// only Countries.Next and Countries.Previous are valid operations (the struct holds no data). Use this to check if the
// data is available.
func (i *iterator) HasSelf() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.ref.metadata().hasSelf()
}

// Self refetches the last countries resource fetched, overwriting the current set.
// If you want a copy of the current Countries struct, use the Countries.Copy method instead.
func (i *iterator) Self(ctx context.Context) {
	i.iterate(ctx, metaKeySelf)
}

// First fetches the first set of countries, overwriting the current one, if available.
func (i *iterator) First(ctx context.Context) {
	i.iterate(ctx, metaKeySelf)
}

// Last fetches the last set of countries, overwriting the current one, if available.
func (i *iterator) Last(ctx context.Context) {
	i.iterate(ctx, metaKeyLast)
}

// Error returns any errors that occurred while calling Countries.Next. It is idempotent, multiple calls after a single
// iteration return the same result.
func (i *iterator) Error() error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.err
}

func (i *iterator) copy(newAPI *API) iterator {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if newAPI == nil {
		newAPI = i.api
	}
	return iterator{
		ref: i.ref,
		api: newAPI,
	}
}

// RefParams is a struct containing the parameters
// used to make requests to any of the Reference APIs
//
// Fields:
//  - Lang is a language.Tag pointer. If it's nil, the
//    API sends the names in all available languages.
//  - Limit represents the number of records returned per request. Default is set
//    to 20, maximum is 100 (if a value bigger than 100 is given, 100 will be taken)
//  - Offset represents the number of records skipped (default is 0). For example,
//    if offset is 20 and limit is 100, the response will contain records from
//    item no. 20 to item no. 119 (100 countries).
type RefParams struct {
	code   string
	Lang   *language.Tag
	Limit  int
	Offset int
}

func writeAmp(sb *strings.Builder, mustWrite *bool) {
	if *mustWrite {
		sb.WriteByte('&')
	} else {
		sb.WriteByte('?')
		*mustWrite = true
	}
}

// ToURL transforms a CountriesParams struct into an URL usable format,
// so that it can be concatenated to te Request API URL.
func (p *RefParams) ToURL() string {
	if p == nil {
		return ""
	}
	mustWriteAmp := false
	var sb strings.Builder
	if p.code != "" {
		sb.WriteString(p.code)
	}
	if p.Lang != nil {
		b, c := p.Lang.Base()
		if c == language.No {
			b, _ = language.English.Base()
		}
		sb.WriteString("?lang=")
		mustWriteAmp = true
		sb.WriteString(b.String())
	}
	if p.Limit != 0 {
		writeAmp(&sb, &mustWriteAmp)
		sb.WriteString("limit=")
		sb.WriteString(strconv.Itoa(p.Limit))
	}
	if p.Offset != 0 {
		writeAmp(&sb, &mustWriteAmp)
		sb.WriteString("offset=")
		sb.WriteString(strconv.Itoa(p.Offset))
	}
	return sb.String()
}
