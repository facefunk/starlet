package scarlet

// Renderable represents anything that can be rendered into final output.
type Renderable interface {
	Render(Builder, bool)
}

// Builder is extracted from strings.Builder. Renderables can write to Builders.
type Builder interface {
	Write(p []byte) (int, error)
	WriteByte(c byte) error
	WriteRune(r rune) (int, error)
	WriteString(s string) (int, error)
}
