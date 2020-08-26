package lufthansa

import (
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

const (
	mdsReferenceAPI = fetchAPI + "/mds-references"
	referenceAPI    = fetchAPI + "/references"
)

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

func (rn *referenceNames) copy() referenceNames {
	r := make(referenceNames)
	for k, v := range *rn {
		r[k] = v
	}
	return r
}

// reference meta types
type (
	metaLinks         map[string]string
	metaLinkUnmarshal struct {
		Rel  string `xml:"Rel,attr" json:"@Rel"`
		Href string `xml:"Href,attr" json:"@Href"`
	}
	meta struct {
		Version    string
		Links      metaLinks
		TotalCount int
	}
	metaUnmarshal struct {
		Version    string              `xml:"Version,attr" json:"@Version"`
		Links      []metaLinkUnmarshal `xml:"Link" json:"Link"`
		TotalCount int                 `xml:"TotalCount" json:"TotalCount"`
	}
)

func (ml *metaLinks) make(arr []metaLinkUnmarshal) {
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

func (m *meta) make(mu *metaUnmarshal) {
	m.TotalCount = mu.TotalCount
	m.Version = mu.Version
	m.Links.make(mu.Links)
}

func (m *meta) copy() meta {
	return meta{
		Version:    m.Version,
		Links:      m.Links.copy(),
		TotalCount: m.TotalCount,
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

// ToURL transforms a CountriesParams struct into an URL usable format,
// so that it can be concatenated to te Request API URL.
func (p *RefParams) ToURL() string {
	if p == nil {
		return ""
	}
	var ret strings.Builder
	if p.code != "" {
		ret.WriteString(p.code)
	}
	if p.Lang != nil {
		b, _ := p.Lang.Base()
		ret.WriteString(fmt.Sprintf("?lang=%s", b.String()))
	}
	if lo := processLimitOffset(p.Limit, p.Offset); lo != "" && strings.Contains(ret.String(), "?") {
		ret.WriteString(fmt.Sprintf("&%s", lo))
	} else if lo != "" {
		ret.WriteString(fmt.Sprintf("?%s", lo))
	}
	return ret.String()
}
