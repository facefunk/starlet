package scarlet

import (
	"sort"
	"strconv"
	"strings"

	"github.com/OneOfOne/xxhash"
)

// Force interface implementation
var _ Renderable = (*CSSRule)(nil)

// CSSRule ...
type CSSRule struct {
	Selector
	Statements []*CSSStatement
	Duplicates []*CSSRule
	Parent     *CSSRule
}

// Render renders the CSS rule into the output stream.
func (rule *CSSRule) Render(output Builder, pretty bool) {
	if len(rule.Statements) == 0 {
		return
	}

	if mapper, ok := output.(*MappingBuilder); ok {
		last := rule.Selector[len(rule.Selector)-1]
		mapper.AddMapping(last.OriginalFile, last.OriginalLine, last.OriginalName)
	}

	_, _ = output.WriteString(rule.SelectorPath(pretty))

	if len(rule.Duplicates) > 0 {
		for _, duplicate := range rule.Duplicates {
			_ = output.WriteByte(',')

			if pretty {
				_ = output.WriteByte(' ')
			}

			_, _ = output.WriteString(duplicate.SelectorPath(pretty))
		}
	}

	if pretty {
		_ = output.WriteByte(' ')
	}

	_ = output.WriteByte('{')

	if pretty {
		_ = output.WriteByte('\n')
	}

	for index, statement := range rule.Statements {
		if pretty {
			_ = output.WriteByte('\t')
		}

		if mapper, ok := output.(*MappingBuilder); ok {
			mapper.AddMapping(statement.OriginalFile, statement.OriginalLine, statement.OriginalName)
		}

		_, _ = output.WriteString(statement.Property)
		_ = output.WriteByte(':')

		if pretty {
			_ = output.WriteByte(' ')
		}

		_, _ = output.WriteString(statement.Value)

		// Remove last semicolon
		if pretty || index != len(rule.Statements)-1 {
			_ = output.WriteByte(';')
		}

		if pretty {
			_ = output.WriteByte('\n')
		}
	}

	_ = output.WriteByte('}')

	if pretty {
		_, _ = output.WriteString("\n\n")
	}
}

// Root ...
func (rule *CSSRule) Root() *CSSRule {
	parent := rule

	for {
		nextParent := parent.Parent

		if nextParent == nil {
			return parent
		}

		parent = nextParent
	}
}

// Copy ...
func (rule *CSSRule) Copy() *CSSRule {
	return &CSSRule{
		Selector:   rule.Selector,
		Statements: rule.Statements,
		Parent:     rule.Parent,
	}
}

// SelectorPath returns the selector string for the rule (recursive, returns absolute path).
func (rule *CSSRule) SelectorPath(pretty bool) string {
	return strings.TrimSpace(rule.Selector.Render(pretty))
}

// StatementsHash returns a hash of all the statements which is used to find duplicate CSS rules.
// Also has the side-effect of sorting the rule's statements in place.
func (rule *CSSRule) StatementsHash() string {
	sort.Slice(rule.Statements, func(i, j int) bool {
		return rule.Statements[i].Property < rule.Statements[j].Property
	})

	hash := xxhash.NewS64(0)

	for _, statement := range rule.Statements {
		_, _ = hash.WriteString(statement.Property)
		_, _ = hash.WriteString(statement.Value)
	}

	return strconv.FormatUint(hash.Sum64(), 16)
}
