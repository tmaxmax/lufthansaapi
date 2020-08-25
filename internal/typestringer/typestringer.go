package typestringer

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type IndexGenFunc func(int) string

type Stringifier struct {
	indentDepth     int
	structSeparator string
	mapSeparator    string
	chanPrefix      string
	indexGen        IndexGenFunc

	indent string
}

type Option func(*Stringifier)

func Indent(depth int) Option {
	return func(s *Stringifier) {
		if depth >= 0 {
			s.indentDepth = depth
		}
	}
}

func StructSeparator(sep string) Option {
	return func(s *Stringifier) {
		s.structSeparator = sep
	}
}

func MapSeparator(sep string) Option {
	return func(s *Stringifier) {
		s.mapSeparator = sep
	}
}

func ChanPrefix(prefix string) Option {
	return func(s *Stringifier) {
		s.chanPrefix = prefix
	}
}

func IndexGen(f IndexGenFunc) Option {
	return func(s *Stringifier) {
		if f != nil {
			s.indexGen = f
		}
	}
}

// NewStringifier creates a stringifier struct that can be used globally to stringify Go variables.
func NewStringifier(opts ...Option) *Stringifier {
	s := &Stringifier{
		indentDepth:     2,
		structSeparator: ": ",
		mapSeparator:    " - ",
		chanPrefix:      "<-",
		indexGen: func(i int) string {
			return "[" + strconv.Itoa(i) + "]"
		},
	}

	for _, o := range opts {
		o(s)
	}

	s.indent = genIndent(s.indentDepth)

	return s
}

// Stringify creates a string representation for the passed interface.
func (s *Stringifier) Stringify(v interface{}, prefix string) string {
	sb := &strings.Builder{}

	s.stringify(sb, reflect.ValueOf(v), prefix, true)
	return sb.String()
}

func (s *Stringifier) stringify(sb *strings.Builder, value reflect.Value, prefix string, writeFirstPrefix bool) {
	stringer, ok := value.Interface().(fmt.Stringer)
	if ok {
		if writeFirstPrefix {
			sb.WriteString(prefix)
		}
		sb.WriteString(stringer.String())
		return
	}
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		if value.IsNil() {
			if writeFirstPrefix {
				sb.WriteString(prefix)
			}
			sb.WriteString("nil")
			return
		}
		s.stringify(sb, value.Elem(), prefix, writeFirstPrefix)
	case reflect.Array, reflect.Slice:
		if value.IsZero() {
			if writeFirstPrefix {
				sb.WriteString(prefix)
			}
			sb.WriteString("(empty slice/array)")
		}
		l := value.Len()
		for i := 0; i < l; i++ {
			if writeFirstPrefix || i > 0 {
				sb.WriteString(prefix)
			}
			indexStr := s.indexGen(i)
			sb.WriteString(indexStr)
			v := value.Index(i)
			s.stringify(sb, v, prefix+genIndent(len(indexStr)), false)
			if i < l-1 {
				sb.WriteByte('\n')
			}
		}
	case reflect.Map:
		if value.IsZero() {
			if writeFirstPrefix {
				sb.WriteString(prefix)
			}
			sb.WriteString("(empty map)")
		}
		l := len(value.MapKeys())
		for i, j := value.MapRange(), 0; i.Next(); j++ {
			s.stringify(sb, i.Key(), prefix, true)
			sb.WriteString(s.mapSeparator)
			v := i.Value()
			if mustWriteNewlineBefore(v) {
				sb.WriteByte('\n')
			}
			s.stringify(sb, v, prefix+s.indent, false)
			if j < l-1 {
				sb.WriteByte('\n')
			}
		}
	case reflect.Struct:
		if value.IsZero() {
			if writeFirstPrefix {
				sb.WriteString(prefix)
			}
			sb.WriteString("(empty struct)")
		}
		var hasExportedFields bool
		numField := value.NumField()
		exportedFields := numField
		for i := 0; i < numField; i++ {
			field := value.Field(i)
			if !field.CanSet() {
				continue
			}
			hasExportedFields = true
			if writeFirstPrefix || i > 0 {
				sb.WriteString(prefix)
			}
			sb.WriteString(value.Type().Field(i).Name + s.structSeparator)
			wfp := false
			if mustWriteNewlineBefore(field) {
				sb.WriteByte('\n')
				wfp = true
			}
			s.stringify(sb, field, prefix+s.indent, wfp)
			if i < exportedFields-1 && value.Field(i+1).CanSet() {
				sb.WriteByte('\n')
			}
		}
		if !hasExportedFields {
			if writeFirstPrefix {
				sb.WriteString(prefix)
			}
			sb.WriteString("(" + value.Type().String() + ")")
		}
	//case reflect.Chan:
	//	if value.IsZero() {
	//		if writeFirstPrefix {
	//			sb.WriteString(prefix)
	//		}
	//		sb.WriteString("(empty channel)")
	//	}
	//	value.Close()
	//	l, i := value.Len(), 0
	//	for v, ok := value.Recv(); ok; v, ok = value.Recv() {
	//		sb.WriteString(prefix + "<-")
	//		s.stringify(sb, v, prefix+s.indent, false)
	//		if i < l-1 {
	//			sb.WriteByte('\n')
	//		}
	//		i++
	//	}
	default:
		if writeFirstPrefix {
			sb.WriteString(prefix)
		}
		sb.WriteString(fmt.Sprintf("%v", value.Interface()))
	}
}

func mustWriteNewlineBefore(v reflect.Value) bool {
	if _, ok := v.Interface().(fmt.Stringer); ok {
		return false
	}
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return mustWriteNewlineBefore(v.Elem())
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct /*, reflect.Chan*/ :
		return true
	}
	return false
}

func genIndent(count int) string {
	s := make([]byte, count)
	for i := range s {
		s[i] = ' '
	}
	return string(s)
}
