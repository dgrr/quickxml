package xml

import "io"

// Writer is used to write the XML elements.
type Writer struct {
	w      io.Writer
	indent string
}

// NewWriter creates a new XML writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w, ""}
}

// Write writes the parsed element.
func (w *Writer) Write(e Element) error {
	return writeString(w.w, e.String())
}

// WriteIndent writes the parsed element indentating the elements.
func (w *Writer) WriteIndent(e Element) error {
	if _, ok := e.(*EndElement); ok {
		w.indent = w.indent[:len(w.indent)-2]
	}

	err := writeString(w.w, w.indent, e.String(), "\n")

	if e, ok := e.(*StartElement); ok && !e.hasEnd {
		w.indent += "  "
	}

	return err
}

func writeString(w io.Writer, strs ...string) (err error) {
	for _, str := range strs {
		_, err = io.WriteString(w, str)
		if err != nil {
			break
		}
	}

	return
}
