package scarlet

import "fmt"

const (
	ElementSelector = iota
	ClassSelector
	IDSelector
	AttributeSelector
	PseudoSelector
	CombinatorSelector
	SeparatorSelector
	NamespaceSelector
	ParentSelector
)

type (
	Selector     []SelectorPart
	SelectorPart struct {
		Name     string
		Type     int
		Operator byte
	}
)

func parseSelector(selector string) Selector {

	var (
		name  []byte
		parts Selector
		t     int
		c     byte
		o     byte
		i     int
	)

	l := len(selector)

	add := func() {
		if len(name) != 0 || t != ElementSelector {
			parts = append(parts, SelectorPart{
				Name:     string(name),
				Type:     t,
				Operator: o,
			})
		}
		name = name[:0]
		o = c
	}

	// Combinator operator, trailing spaces are meaningless
	combinator := func(nextT int) {
		// Consume trailing spaces
		for {
			j := i + 1
			if j == l || selector[j] != ' ' {
				break
			}
			i = j
		}
		// Add previous part
		add()
		t = nextT
		// Add combinator
		add()
		// Reset
		t = ElementSelector
		o = 0
	}

	for ; i < l; i++ {
		c = selector[i]
		switch c {
		case ',':
			combinator(SeparatorSelector)
		case ' ', '>', '+', '~':
			combinator(CombinatorSelector)
		case '&':
			// Not a combinator, trailing spaces are meaningful
			add()
			t = ParentSelector
			add()
			t = ElementSelector
			o = 0
		case '.':
			add()
			t = ClassSelector
		case '#':
			add()
			t = IDSelector
		case '[':
			add()
			t = AttributeSelector
			for {
				j := i + 1
				if j == l {
					break
				}
				i = j
				c = selector[i]
				if c == ']' {
					break
				}
				name = append(name, c)
			}
		case ':':
			for {
				j := i + 1
				if j == l || selector[j] != ':' {
					break
				}
				i = j
			}
			add()
			t = PseudoSelector
		case '|':
			t = NamespaceSelector
			o = c
			add()
			t = ElementSelector
			o = 0
		default:
			name = append(name, c)
		}
	}
	add()
	return parts
}

func (selector Selector) Prepend(parent Selector) (Selector, error) {

	if parent == nil || len(parent) == 0 {
		return selector, nil
	}

	out := make(Selector, 0, len(selector)+len(parent)+1)
	found := false
	for _, part := range selector {
		switch part.Type {
		case ParentSelector:
			out = append(out, parent...)
			found = true
		case SeparatorSelector:
			return nil, fmt.Errorf("can't prepend unsplit selector: %s", selector.Render())
		default:
			out = append(out, part)
		}
	}

	if !found {
		ret := append(parent, SelectorPart{
			Type:     CombinatorSelector,
			Operator: ' ',
		})
		ret = append(ret, out...)
		return ret, nil
	}

	return out, nil
}

func (selector Selector) Split() []Selector {
	selectors := make([]Selector, 0, len(selector))
	l := 0
	for p, part := range selector {
		if part.Type == SeparatorSelector {
			selectors = append(selectors, selector[l:p])
			l = p + 1
		}
	}
	selectors = append(selectors, selector[l:])
	return selectors
}

func (selector Selector) Render() string {
	out := ""
	if len(selector) == 0 {
		return out
	}
	if selector[0].Type == NamespaceSelector && selector[0].Name == "" {
		out = " "
	}
	for _, part := range selector {
		switch part.Type {
		case AttributeSelector:
			out += fmt.Sprintf("[%s]", part.Name)
		case NamespaceSelector:
			out += fmt.Sprintf("%s|", part.Name)
		case ElementSelector:
			out += part.Name
		default:
			out += fmt.Sprintf("%c%s", part.Operator, part.Name)
		}
	}
	return out
}
