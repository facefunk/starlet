package starlet

import (
	"strings"

	"github.com/neelance/sourcemap"
)

// MappingBuilder implements Builder and builds a sourcemap.Map as it goes.
type MappingBuilder struct {
	strings.Builder
	sourcemap.Map
	col              int
	line             int
	originalFilename string
	originalLine     int
	originalName     string
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
func (b *MappingBuilder) Write(p []byte) (int, error) {
	b.count(p...)
	return b.Builder.Write(p)
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b *MappingBuilder) WriteByte(c byte) error {
	b.count(c)
	return b.Builder.WriteByte(c)
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to b's buffer.
// It returns the length of r and a nil error.
func (b *MappingBuilder) WriteRune(r rune) (int, error) {
	l, _ := b.Builder.WriteRune(r)
	if l > 1 {
		b.col += l
	} else {
		b.count(byte(r))
	}
	return l, nil
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b *MappingBuilder) WriteString(s string) (int, error) {
	b.count([]byte(s)...)
	return b.Builder.WriteString(s)
}

// AddMapping adds the mapping recorded thus far to the embedded sourcemap and starts a new mapping with the provided
// arguments.
func (b *MappingBuilder) AddMapping(filename string, line int, name string) {
	mapping := &sourcemap.Mapping{
		GeneratedLine:   b.line,
		GeneratedColumn: b.col,
		OriginalFile:    b.originalFilename,
		OriginalLine:    b.originalLine,
		OriginalColumn:  0,
		OriginalName:    b.originalName,
	}
	b.Map.AddMapping(mapping)
	b.originalFilename = filename
	b.originalLine = line
	b.originalName = name
}

// count keeps track of the current column and line
func (b *MappingBuilder) count(p ...byte) {
	for _, by := range p {
		if by == '\n' {
			b.col = 0
			b.line++
			continue
		}
		b.col++
	}
}
