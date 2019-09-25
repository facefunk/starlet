package scarlet

// Force interface implementation
var _ Renderable = (*MediaQuery)(nil)

// MediaQuery ...
type MediaQuery struct {
	Selector string
	Rules    []*CSSRule
}

// Render renders the media query to the output stream.
func (media *MediaQuery) Render(output Builder, pretty bool) {
	if len(media.Rules) == 0 {
		return
	}

	_, _ = output.WriteString(media.Selector)

	if pretty {
		_ = output.WriteByte(' ')
	}

	_ = output.WriteByte('{')

	if pretty {
		_ = output.WriteByte('\n')
	}

	for _, rule := range media.Rules {
		rule.Render(output, pretty)
	}

	_ = output.WriteByte('}')

	if pretty {
		_ = output.WriteByte('\n')
	}
}
