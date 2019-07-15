package scarlet

import (
	"testing"
)

func TestParseSelector(t *testing.T) {

	selector := "namespace|div div, |div>div, div+div~.class#id[attribute~=value]::after:active>"
	wantedSelector := "namespace|div div,|div>div,div+div~.class#id[attribute~=value]:after:active>"
	parts := parseSelector(selector)
	wanted := []SelectorPart{
		{Name: "namespace", Type: NamespaceSelector, Operator: '|'},
		{Name: "div", Type: ElementSelector, Operator: 0},
		{Name: "", Type: CombinatorSelector, Operator: ' '},
		{Name: "div", Type: ElementSelector, Operator: 0},
		{Name: "", Type: SeparatorSelector, Operator: ','},
		{Name: "", Type: NamespaceSelector, Operator: '|'},
		{Name: "div", Type: ElementSelector, Operator: 0},
		{Name: "", Type: CombinatorSelector, Operator: '>'},
		{Name: "div", Type: ElementSelector, Operator: 0},
		{Name: "", Type: SeparatorSelector, Operator: ','},
		{Name: "div", Type: ElementSelector, Operator: 0},
		{Name: "", Type: CombinatorSelector, Operator: '+'},
		{Name: "div", Type: ElementSelector, Operator: 0},
		{Name: "", Type: CombinatorSelector, Operator: '~'},
		{Name: "class", Type: ClassSelector, Operator: '.'},
		{Name: "id", Type: IDSelector, Operator: '#'},
		{Name: "attribute~=value", Type: AttributeSelector, Operator: '['},
		{Name: "after", Type: PseudoSelector, Operator: ':'},
		{Name: "active", Type: PseudoSelector, Operator: ':'},
		{Name: "", Type: CombinatorSelector, Operator: '>'},
	}
	if len(parts) != len(wanted) {
		t.Error("Unexpected number of parts")
	}
	for i, part := range parts {
		want := wanted[i]
		if part.Name != want.Name || part.Type != want.Type || part.Operator != want.Operator {
			t.Errorf("Part %#v != %#v", part, want)
		}
	}
	renderedSelector := parts.Render()
	if renderedSelector != wantedSelector {
		t.Errorf("Rendered selector %s != %s", renderedSelector, wantedSelector)
	}
}
