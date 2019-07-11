package scarlet

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/aerogo/codetree"
)

func TestRename(t *testing.T) {
	reader, _ := os.Open("testdata/classes.scarlet")
	tree, err := codetree.FromReader(reader)
	defer reader.Close()
	if err != nil {
		t.Fatal(err)
		return
	}

	compiler := FromCodeTree(tree)
	renamingMap := compiler.RenameClasses()
	css, err := compiler.Render(false)
	if err != nil {
		fmt.Printf("Error compiling:%s", err)
		return
	}
	mapJSON, err := json.Marshal(renamingMap.Map)
	if err != nil {
		fmt.Printf("Error marshaling:%s", err)
		return
	}
	fmt.Printf("renamingMap:%s\n", mapJSON)
	fmt.Print(css)
}

func TestAssign(t *testing.T) {
	m := NewRenamingMap()
	for i := 0; i < 1000; i++ {
		t.Log(m.Assign(fmt.Sprintf("LongID%d", i)))
	}
}
