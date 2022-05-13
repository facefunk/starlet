package starlet

// Force interface implementation
var _ Renderable = (*MediaGroup)(nil)

// MediaGroup ...
type MediaGroup struct {
	Operator string
	Size     string
	Property string
	Rules    []*CSSRule
}

// Render renders the media group to the output stream.
func (media *MediaGroup) Render(output Builder, pretty bool) {
	if len(media.Rules) == 0 {
		return
	}

	_, _ = output.WriteString("@media all and (")

	switch media.Operator {
	case "<":
		_, _ = output.WriteString("max")
	case ">":
		_, _ = output.WriteString("min")
	default:
		panic("Invalid screen size operator in media query")
	}

	_ = output.WriteByte('-')
	_, _ = output.WriteString(media.Property)
	_ = output.WriteByte(':')

	if pretty {
		_ = output.WriteByte(' ')
	}

	_, _ = output.WriteString(media.Size)
	_ = output.WriteByte(')')

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
