package easyphpserialize

import (
	"github.com/golyakov/easyphpserialize/plexer"
	"github.com/golyakov/easyphpserialize/pwriter"
)

// RawMessage is a raw piece of PHPSerialize (number, string, bool, object, array or
// null) that is extracted without parsing and output as is during marshaling.
type RawMessage []byte

// MarshalEasyPHPSerialize does PHPSerialize marshaling using easyphpserialize interface.
func (v *RawMessage) MarshalEasyPHPSerialize(w *pwriter.Writer) {
	if len(*v) == 0 {
		w.RawString("null")
	} else {
		w.Raw(*v, nil)
	}
}

// UnmarshalEasyPHPSerialize does PHPSerialize unmarshaling using easyphpserialize interface.
func (v *RawMessage) UnmarshalEasyPHPSerialize(l *plexer.Lexer) {
	*v = RawMessage(l.Raw())
}

// UnmarshalPHPSerialize implements easyphpserialize.Unmarshaler interface.
func (v *RawMessage) UnmarshalPHPSerialize(data []byte) error {
	*v = data
	return nil
}

var nullBytes = []byte("null")

// MarshalPHPSerialize implements easyphpserialize.Marshaler interface.
func (v RawMessage) MarshalPHPSerialize() ([]byte, error) {
	if len(v) == 0 {
		return nullBytes, nil
	}
	return v, nil
}

// IsDefined is required for integration with omitempty easyphpserialize logic.
func (v *RawMessage) IsDefined() bool {
	return len(*v) > 0
}
