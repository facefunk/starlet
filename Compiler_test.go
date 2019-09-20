package scarlet_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aerogo/codetree"
	"github.com/aerogo/scarlet"
	"github.com/akyoto/color"
)

func testFile(t *testing.T, filePath string, result string) {
	src, _ := ioutil.ReadFile(filePath)
	code := string(src)

	start := time.Now()
	css, _ := scarlet.Compile(code, false)
	elapsed := time.Since(start)

	fmt.Println(css)

	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("Lines:", color.YellowString("%d", len(strings.Split(css, "\n"))))
	fmt.Println("Size: ", color.YellowString("%d", len(css)))
	fmt.Println("Time: ", color.GreenString("%v", elapsed))

	if css != result {
		t.Error("Unexpected output")
	}
}

func TestCompiler1(t *testing.T) {
	testFile(t, "testdata/test.scarlet", `a{color:black}a :hover{color:blue}a :hover div{color:red}p :hover{color:blue}p :hover div{color:red}@keyframes appear{0%{opacity:0}100%{opacity:1}}@media all and (min-width:900px){p{display:flex;animation-name:appear}}`)
}

func TestCompiler2(t *testing.T) {
	testFile(t, "testdata/test2.scarlet", `:root{--text-color:blue;--text-hover-color:var(--text-color);--gradient:linear-gradient(to bottom,0% var(--text-color),100% var(--text-color));}body{display:flex;flex-direction:row;color:var(--text-color);background-color:#202020}p{display:flex;flex-direction:row;color:blue;background-color:#202020}a{color:red}a :hover{color:var(--text-hover-color)}a :hover div{width:100%}a :hover div img{height:100%}a :active{color:blue}#content{color:red}#content :hover{color:red}#content >div{color:orange}#content img{border:none}#content [aria-class="button"]{color:green}div :hover{color:white}div span{display:none}div address{display:none}p :hover{color:white}p span{display:none}p address{display:none}h1{display:none}h2{display:none}@media all and (min-height: 320px){body{background-color:red}}`)
}

func BenchmarkCompiler(b *testing.B) {
	src, _ := ioutil.ReadFile("testdata/test.scarlet")
	code := string(src)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := scarlet.Compile(code, false)

			if err != nil {
				b.Fail()
			}
		}
	})
}

func TestCompilerFilterTags(t *testing.T) {
	tags := []string{"a", "b", "body", "br", "button", "div", "fieldset", "footer", "form", "h1", "h2", "head",
		"header", "html", "iframe", "input", "legend", "li", "link", "meta", "noscript", "ol", "p", "pre", "script",
		"span", "string", "style", "table", "textarea", "title", "ul"}
	reader, _ := os.Open("testdata/normalize.scarlet")
	tree, err := codetree.FromReader(reader)
	defer reader.Close()
	if err != nil {
		t.Fatal(err)
		return
	}
	css, err := scarlet.FromCodeTree(tree).FilterTags(tags).Render(false)
	if err != nil {
		t.Fatalf("Error compiling:%s", err)
		return
	}
	expected := `footer{display:block}header{display:block}[hidden]{display:none}html{font-size:100%;overflow-y:scroll;-webkit-text-size-adjust:100%;-ms-text-size-adjust:100%}body{margin:0;font-size:13px;line-height:1.231}body{font-family:sans-serif;color:#222}button{font-family:sans-serif;color:#222}input{font-family:sans-serif;color:#222}textarea{font-family:sans-serif;color:#222}::-moz-selection{background:#2989C9;color:#fff;text-shadow:none}::selection{background:#2989C9;color:#fff;text-shadow:none}a{color:#00e}a:visited{color:#551a8b}a:hover{color:#06e}a:focus{outline:thin dotted}a:hover{outline:0}a:active{outline:0}b{font-weight:bold}pre{font-family:monospace,serif;_font-family:'courier new',monospace;font-size:1em}pre{white-space:pre;white-space:pre-wrap;word-wrap:break-word}ul{margin:1em 0;padding:0 0 0 40px}ol{margin:1em 0;padding:0 0 0 40px}nav ul{list-style:none;list-style-image:none;margin:0;padding:0}nav ol{list-style:none;list-style-image:none;margin:0;padding:0}form{margin:0}fieldset{border:0;margin:0;padding:0}legend{border:0;*margin-left:-7px;padding:0}button{font-size:100%;margin:0;vertical-align:baseline;*vertical-align:middle}input{font-size:100%;margin:0;vertical-align:baseline;*vertical-align:middle}textarea{font-size:100%;margin:0;vertical-align:baseline;*vertical-align:middle}button{line-height:normal;*overflow:visible}input{line-height:normal;*overflow:visible}table button{*overflow:auto}table input{*overflow:auto}button{cursor:pointer;-webkit-appearance:button}input[type="button"]{cursor:pointer;-webkit-appearance:button}input[type="reset"]{cursor:pointer;-webkit-appearance:button}input[type="submit"]{cursor:pointer;-webkit-appearance:button}input[type="checkbox"]{box-sizing:border-box}input[type="radio"]{box-sizing:border-box}input[type="search"]{-webkit-appearance:textfield;-moz-box-sizing:content-box;-webkit-box-sizing:content-box;box-sizing:content-box}input[type="search"]::-webkit-search-decoration{-webkit-appearance:none}button::-moz-focus-inner{border:0;padding:0}input::-moz-focus-inner{border:0;padding:0}textarea{overflow:auto;vertical-align:top;resize:vertical}input:invalid{background-color:#f0dddd}textarea:invalid{background-color:#f0dddd}table{border-collapse:collapse;border-spacing:0}@media print{a{text-decoration:underline}a:visited{text-decoration:underline}a[href]:after{content:" (" attr(href) ")"}a[href^="javascript:"]:after{content:""}a[href^="#"]:after{content:""}pre{border:1px solid #999;page-break-inside:avoid}p{orphans:3;widows:3}h2{orphans:3;widows:3}h2{page-break-after:avoid}}`
	if css != expected {
		t.Errorf("Unexpected output: %s", css)
	}
}
