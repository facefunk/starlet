package scarlet

import (
	"bytes"
	"strings"

	"github.com/aerogo/codetree"
)

// Compile compiles the given scarlet code to a CSS string.
func Compile(src string) (string, error) {
	tree, err := codetree.New(src)

	if err != nil {
		return "", err
	}

	var output bytes.Buffer

	for _, rule := range compileChildren(tree, nil) {
		if len(rule.Statements) == 0 {
			continue
		}

		if rule.Parent != nil {
			output.WriteString(rule.Parent.Selector)

			if !strings.HasPrefix(rule.Selector, ":") && !strings.HasPrefix(rule.Selector, "[") {
				output.WriteString(" ")
			}
		}

		output.WriteString(rule.Selector)
		output.WriteString(" {\n")

		for _, statement := range rule.Statements {
			output.WriteString("\t")
			output.WriteString(compileStatement(statement))
			output.WriteString("\n")
		}

		output.WriteString("}\n\n")
	}

	return strings.TrimRight(output.String(), "\n"), nil
}
