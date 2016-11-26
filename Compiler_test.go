package scarlet

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestCompiler(t *testing.T) {
	styl, _ := ioutil.ReadFile("test.styl")
	css, _ := Compile(string(styl))

	fmt.Println(css)
}
