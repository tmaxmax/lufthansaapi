package lufthansa

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	tjson "github.com/tmaxmax/json"

	"github.com/tmaxmax/lufthansaapi/pkg/typestringer"
)

var stringifier = typestringer.NewStringifier()

// processLimitOffset transforms the general limit and offset parameters,
// which are available for most of the API requests, into an usable string
// for creating the API request URL.
func processLimitOffset(l, o int) string {
	var ret strings.Builder
	if l != 0 {
		ret.WriteString(fmt.Sprintf("limit=%d", l))
	}
	if o != 0 {
		if l == 0 {
			ret.WriteString(fmt.Sprintf("&offset=%d", o))
		} else {
			ret.WriteString(fmt.Sprintf("offset=%d", o))
		}
	}
	return ret.String()
}

// decode is a helper for decoding API responses.
func decode(r io.ReadCloser, v interface{}) error {
	data, err := readAllCloser(r)
	if err != nil {
		return err
	}
	switch mimeType(data) {
	case "text/xml", "application/xml":
		return xml.Unmarshal(data, v)
	case "application/json":
		return tjson.Unmarshal(data, v)
	}
	return ErrUnsupportedFormat
}

// readAllCloser reads all data from r and closes it
func readAllCloser(r io.ReadCloser) ([]byte, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if err = r.Close(); err != nil {
		return nil, err
	}
	return data, nil
}

// mimeType returns the mime type of the passed data, without charset
func mimeType(data []byte) string {
	return strings.Split(mimetype.Detect(data).String(), ";")[0]
}
