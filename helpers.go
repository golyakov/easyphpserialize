// Package easyphpserialize contains marshaler/unmarshaler interfaces and helper functions.
package easyphpserialize

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"unsafe"

	"github.com/golyakov/easyphpserialize/plexer"
	"github.com/golyakov/easyphpserialize/pwriter"
)

// Marshaler is an easyphpserialize-compatible marshaler interface.
type Marshaler interface {
	MarshalEasyPHPSerialize(w *pwriter.Writer)
}

// Marshaler is an easyphpserialize-compatible unmarshaler interface.
type Unmarshaler interface {
	UnmarshalEasyPHPSerialize(w *plexer.Lexer)
}

// MarshalerUnmarshaler is an easyphpserialize-compatible marshaler/unmarshaler interface.
type MarshalerUnmarshaler interface {
	Marshaler
	Unmarshaler
}

// Optional defines an undefined-test method for a type to integrate with 'omitempty' logic.
type Optional interface {
	IsDefined() bool
}

// UnknownsUnmarshaler provides a method to unmarshal unknown struct fileds and save them as you want
type UnknownsUnmarshaler interface {
	UnmarshalUnknown(in *plexer.Lexer, key string)
}

// UnknownsMarshaler provides a method to write additional struct fields
type UnknownsMarshaler interface {
	MarshalUnknowns(w *pwriter.Writer, first bool)
}

func isNilInterface(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}

// Marshal returns data as a single byte slice. Method is suboptimal as the data is likely to be copied
// from a chain of smaller chunks.
func Marshal(v Marshaler) ([]byte, error) {
	if isNilInterface(v) {
		return nullBytes, nil
	}

	w := pwriter.Writer{}
	v.MarshalEasyPHPSerialize(&w)
	return w.BuildBytes()
}

// MarshalToWriter marshals the data to an io.Writer.
func MarshalToWriter(v Marshaler, w io.Writer) (written int, err error) {
	if isNilInterface(v) {
		return w.Write(nullBytes)
	}

	pw := pwriter.Writer{}
	v.MarshalEasyPHPSerialize(&pw)
	return pw.DumpTo(w)
}

// MarshalToHTTPResponseWriter sets Content-Length and Content-Type headers for the
// http.ResponseWriter, and send the data to the writer. started will be equal to
// false if an error occurred before any http.ResponseWriter methods were actually
// invoked (in this case a 500 reply is possible).
func MarshalToHTTPResponseWriter(v Marshaler, w http.ResponseWriter) (started bool, written int, err error) {
	if isNilInterface(v) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len(nullBytes)))
		written, err = w.Write(nullBytes)
		return true, written, err
	}

	jw := pwriter.Writer{}
	v.MarshalEasyPHPSerialize(&jw)
	if jw.Error != nil {
		return false, 0, jw.Error
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(jw.Size()))

	started = true
	written, err = jw.DumpTo(w)
	return
}

// Unmarshal decodes the PHPSerialize in data into the object.
func Unmarshal(data []byte, v Unmarshaler) error {
	l := plexer.Lexer{Data: data}
	v.UnmarshalEasyPHPSerialize(&l)
	return l.Error()
}

// UnmarshalFromReader reads all the data in the reader and decodes as PHPSerialize into the object.
func UnmarshalFromReader(r io.Reader, v Unmarshaler) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	l := plexer.Lexer{Data: data}
	v.UnmarshalEasyPHPSerialize(&l)
	return l.Error()
}
