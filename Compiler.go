package starlet

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/aerogo/codetree"
)

// compileChildren returns the CSS rules for a given code tree.
// It iterates over the child nodes and finds the CSS rules.
func compileChildren(node *codetree.CodeTree, parent *CSSRule, state *State) (
	[]*CSSRule, []*MediaGroup, []*MediaQuery, []*Animation) {

	// Comments
	if node.Type == codetree.CommentType {
		return nil, nil, nil, nil
	}

	var rules []*CSSRule
	var mediaGroups []*MediaGroup
	var mediaQueries []*MediaQuery
	var animations []*Animation
	var selectorsOnPreviousLines []Selector

	// Iterate over child nodes
	for _, child := range node.Children {

		// Don't consider comments as valid children
		hasChildren := false
		if len(child.Children) > 0 {
			for _, grandchild := range child.Children {
				if grandchild.Type != codetree.CommentType {
					hasChildren = true
					break
				}
			}
		}

		// Nodes with no children are selectors on previous lines or statements of some sort
		if !hasChildren {
			// Comments
			if child.Type == codetree.CommentType {
				continue
			}

			// Selector on previous line
			if strings.HasSuffix(child.Line, ",") {
				name := child.Line[:len(child.Line)-1]
				selector := parseSelector(name)
				selector.SetMapping(child.GetFilename(), child.LineNumber, name)
				selectors := selector.Split()
				selectorsOnPreviousLines = append(selectorsOnPreviousLines, selectors...)
				continue
			}

			equal := strings.IndexByte(child.Line, '=')

			if equal != -1 {
				value := strings.TrimSpace(child.Line[equal+1:])
				value = insertVariableValues(value, state)
				value = optimizeColors(value)

				if strings.HasPrefix(child.Line, "const ") {
					// Constants
					name := strings.TrimSpace(child.Line[len("const "):equal])
					state.Constants[name] = value
				} else {
					// Variables
					name := strings.TrimSpace(child.Line[:equal])
					state.VariableNames = append(state.VariableNames, name)
					state.Variables[name] = value
				}
			} else if parent != nil && strings.IndexByte(child.Line, ' ') != -1 {
				// Statements
				statement := compileStatement(child.Line, state)
				statement.SetMapping(child.GetFilename(), child.LineNumber, child.Line)
				parent.Statements = append(parent.Statements, statement)
			} else {
				// Mixin calls
				mixin, exists := state.Mixins[child.Line]

				if exists && parent != nil {
					mixinRules := mixin.Apply(parent)
					prependRules(parent, mixinRules...)
					rules = append(rules, mixinRules...)
				} else {
					panic(fmt.Sprintf("invalid statement %s", child))
				}
			}
			continue
		}

		// Mixin
		if strings.HasPrefix(child.Line, "mixin ") {
			name := child.Line[len("mixin "):]

			mixin := &Mixin{
				Root:  &CSSRule{},
				Rules: []*CSSRule{},
			}

			childRules, _, _, _ := compileChildren(child, mixin.Root, state)
			mixin.Rules = append(mixin.Rules, childRules...)
			state.Mixins[name] = mixin
			continue
		}

		// Media query
		if strings.HasPrefix(child.Line, "@media") {
			selector := strings.TrimSpace(child.Line)

			media := &MediaQuery{
				Selector: selector,
			}

			media.Rules, _, _, _ = compileChildren(child, nil, state)
			mediaQueries = append(mediaQueries, media)
			continue
		}

		// Media query by size
		if strings.HasPrefix(child.Line, "< ") || strings.HasPrefix(child.Line, "> ") {
			// If this is a top level definition, it's a screen size query
			if child.Indent == 0 {
				media := &MediaGroup{}
				parts := strings.Split(child.Line, " ")

				media.Operator = parts[0]
				media.Size = parts[1]

				if len(parts) >= 3 {
					media.Property = parts[2]
				} else {
					media.Property = "width"
				}

				media.Rules, _, _, _ = compileChildren(child, nil, state)
				mediaGroups = append(mediaGroups, media)
				continue
			}
		}

		// Animation
		if strings.HasPrefix(child.Line, "animation ") {
			anim := &Animation{
				Name: child.Line[len("animation "):],
			}

			anim.Keyframes, _, _, _ = compileChildren(child, nil, state)
			animations = append(animations, anim)
			continue
		}

		selector := parseSelector(child.Line)
		selector.SetMapping(child.GetFilename(), child.LineNumber, child.Line)
		selectors := selector.Split()

		// Append selectors from previous lines
		selectors = append(selectorsOnPreviousLines, selectors...)
		selectorsOnPreviousLines = []Selector(nil)

		for s := range selectors {
			// Child rule
			rule := &CSSRule{
				Selector: selectors[s],
				Parent:   parent,
			}
			prependRules(parent, rule)
			rules = append(rules, rule)
			childRules, childGroups, childQueries, childAnimations := compileChildren(child, rule, state)
			rules = append(rules, childRules...)
			mediaGroups = append(mediaGroups, childGroups...)
			mediaQueries = append(mediaQueries, childQueries...)
			animations = append(animations, childAnimations...)
		}
	}

	return rules, mediaGroups, mediaQueries, animations
}

func prependRules(parent *CSSRule, rules ...*CSSRule) {
	for _, rule := range rules {
		if rule.Selector != nil && parent != nil && parent.Selector != nil && rule != parent {
			var err error
			rule.Selector, err = rule.Selector.Prepend(parent.Selector)
			if err != nil {
				panic(err)
			}
		}
	}
}

// compileStatement compiles a Scarlet statement to CSS.
func compileStatement(statement string, state *State) *CSSStatement {
	space := strings.IndexByte(statement, ' ')

	if space == -1 {
		panic("Invalid statement: " + statement)
	}

	value := strings.TrimSpace(statement[space:])

	// Optimize color values
	value = insertVariableValues(value, state)
	value = optimizeColors(value)

	return &CSSStatement{
		Property: statement[:space],
		Value:    value,
	}
}

// insertVariableValues inserts the values of variables directly into the string.
func insertVariableValues(expression string, state *State) string {
	// EOF
	runes := append([]rune(expression), ' ')
	buffer := strings.Builder{}
	ignore := ignoreReader{}
	cursor := 0

	for index, char := range runes {
		if ignore.canIgnore(char) {
			buffer.WriteRune(char)
			cursor = index + 1
			continue
		}

		if char != '-' && (unicode.IsSpace(char) || unicode.IsPunct(char)) {
			if index != cursor {
				token := string(runes[cursor:index])

				// Check dynamic CSS variables
				_, exists := state.Variables[token]

				if exists {
					buffer.WriteString("var(--")
					buffer.WriteString(token)
					buffer.WriteString(")")
				} else {
					// Check constants
					value, exists := state.Constants[token]

					if exists {
						buffer.WriteString(value)
					} else {
						buffer.WriteString(token)
					}
				}
			}

			if index == len(runes)-1 {
				break
			}

			buffer.WriteRune(char)
			cursor = index + 1
		}
	}

	return buffer.String()
}
