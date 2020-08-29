package util

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	tjson "github.com/tmaxmax/json"

	"github.com/tmaxmax/lufthansaapi/pkg/typestringer"
)

var (
	Stringer = typestringer.NewStringifier()
	// ErrUnsupportedFormat is returned when decoding a response in an unsupported format
	ErrUnsupportedFormat = errors.New("lufthansaapi: Decode: unsupported format")
)

// Decode is a helper for decoding API responses.
func Decode(r io.ReadCloser, v interface{}) error {
	data, err := ReadAll(r)
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

// ReadAll reads all data from r and closes it
func ReadAll(r io.ReadCloser) ([]byte, error) {
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
