package scarlet

// Force interface implementation
var _ Renderable = (*Animation)(nil)

// Animation ...
type Animation struct {
	Name      string
	Keyframes []*CSSRule
}

// Render renders the animation to the output stream.
func (anim *Animation) Render(output Builder, pretty bool) {
	_, _ = output.WriteString("@keyframes ")
	_, _ = output.WriteString(anim.Name)

	if pretty {
		_ = output.WriteByte(' ')
	}

	_ = output.WriteByte('{')

	if pretty {
		_ = output.WriteByte('\n')
	}

	for _, keyframe := range anim.Keyframes {
		keyframe.Render(output, pretty)
	}

	_ = output.WriteByte('}')

	if pretty {
		_ = output.WriteByte('\n')
	}
}
