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
		Operator string
	}
)

func parseSelector(selector string) Selector {

	var (
		name  []byte
		parts Selector
		t     int
		c     byte
		o     string
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
		o = string(c)
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
		o = ""
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
			o = ""
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
			element := false
			for {
				j := i + 1
				if j == l || selector[j] != ':' {
					break
				}
				i = j
				element = true
			}
			add()
			t = PseudoSelector
			if element {
				o = "::"
			}
		case '|':
			t = NamespaceSelector
			o = string(c)
			add()
			t = ElementSelector
			o = ""
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
			return nil, fmt.Errorf("can't prepend unsplit selector: %s", selector.Render(true))
		default:
			out = append(out, part)
		}
	}

	if !found {
		// Make sure to copy parent so as not to return duplicate pointers
		ret := append(Selector(nil), parent...)
		ret = append(ret, SelectorPart{
			Type:     CombinatorSelector,
			Operator: " ",
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

func (selector Selector) Render(pretty bool) string {
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
		case SeparatorSelector:
			if !pretty {
				out += string(part.Operator)
				continue
			}
			out += fmt.Sprintf("%s ", part.Operator)
		case CombinatorSelector:
			if part.Operator == " " || !pretty {
				out += string(part.Operator)
				continue
			}
			out += fmt.Sprintf(" %s", part.Operator)
		default:
			out += fmt.Sprintf("%s%s", part.Operator, part.Name)
		}
	}
	return out
}

func (selector Selector) String() string {
	return selector.Render(false)
}
