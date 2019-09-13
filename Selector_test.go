package scarlet

import (
	"testing"
)

func TestParseSelector(t *testing.T) {

	selector := "namespace|div div, |div>div, div+div~.class#id[attribute~=value]::after:active>"
	wantedSelector := "namespace|div div,|div>div,div+div~.class#id[attribute~=value]::after:active>"
	parts := parseSelector(selector)
	wanted := []SelectorPart{
		{Name: "namespace", Type: NamespaceSelector, Operator: "|"},
		{Name: "div", Type: ElementSelector, Operator: ""},
		{Name: "", Type: CombinatorSelector, Operator: " "},
		{Name: "div", Type: ElementSelector, Operator: ""},
		{Name: "", Type: SeparatorSelector, Operator: ","},
		{Name: "", Type: NamespaceSelector, Operator: "|"},
		{Name: "div", Type: ElementSelector, Operator: ""},
		{Name: "", Type: CombinatorSelector, Operator: ">"},
		{Name: "div", Type: ElementSelector, Operator: ""},
		{Name: "", Type: SeparatorSelector, Operator: ","},
		{Name: "div", Type: ElementSelector, Operator: ""},
		{Name: "", Type: CombinatorSelector, Operator: "+"},
		{Name: "div", Type: ElementSelector, Operator: ""},
		{Name: "", Type: CombinatorSelector, Operator: "~"},
		{Name: "class", Type: ClassSelector, Operator: "."},
		{Name: "id", Type: IDSelector, Operator: "#"},
		{Name: "attribute~=value", Type: AttributeSelector, Operator: "["},
		{Name: "after", Type: PseudoSelector, Operator: "::"},
		{Name: "active", Type: PseudoSelector, Operator: ":"},
		{Name: "", Type: CombinatorSelector, Operator: ">"},
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
	renderedSelector := parts.Render(false)
	if renderedSelector != wantedSelector {
		t.Errorf("Rendered selector %s != %s", renderedSelector, wantedSelector)
	}
}

func TestParentSelector(t *testing.T) {
	s1 := "div.class2"
	s2 := ".class1 & .class3"
	sel1 := parseSelector(s1)
	sel2 := parseSelector(s2)
	cs1, err := sel2.Prepend(sel1)
	if err != nil {
		t.Error(err)
	}
	r1 := cs1.Render(false)
	if r1 != ".class1 div.class2 .class3" {
		t.Errorf("Unexpected selector: %s", r1)
	}

	sel3 := parseSelector(s1)
	cs2, err := sel3.Prepend(sel1)
	if err != nil {
		t.Error(err)
	}
	r2 := cs2.Render(false)
	if r2 != "div.class2 div.class2" {
		t.Errorf("Unexpected selector: %s", r1)
	}
}
