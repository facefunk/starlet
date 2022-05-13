package starlet

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/aerogo/codetree"
)

func TestRename(t *testing.T) {
	reader, _ := os.Open("testdata/classes.strlt")
	tree, err := codetree.FromReader(reader)
	defer reader.Close()
	if err != nil {
		t.Fatal(err)
		return
	}

	compiler := FromCodeTree(tree)
	renamingMap := compiler.RenameClasses()
	builder := &strings.Builder{}
	compiler.Render(builder, false)
	css := builder.String()

	wantedMap := &RenamingMap{
		Map: map[string]string{
			"class": "a",
			"foo":   "b",
			"bar":   "c",
		},
		Len: 3,
	}

	if !reflect.DeepEqual(renamingMap, wantedMap) {
		t.Errorf("RenamingMap doesn't match: %#v != %#v", renamingMap, wantedMap)
	}

	wantedCSS := ".a .a-b{font-weight:bold}.a .a-c{font-style:italic}"
	if css != wantedCSS {
		t.Errorf("CSS output doesn't match: %s != %s", css, wantedCSS)
	}

}

func TestAssign(t *testing.T) {
	m := NewRenamingMap()
	a := ""
	for i := 0; i < 1000; i++ {
		a = m.Assign(fmt.Sprintf("LongID%d", i))
	}
	if a != "aa1" {
		t.Errorf("Last assigned class name was %s", a)
	}
}
