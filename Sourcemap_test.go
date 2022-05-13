package starlet

import (
	"strings"
	"testing"

	"github.com/aerogo/codetree"
)

func TestSourcemap(t *testing.T) {
	list := []string{"testdata/test.strlt", "testdata/test2.strlt"}
	outputs := []string{
		`{"version":3,"sources":["testdata/test.strlt","testdata/test2.strlt"],"names":["a","color black",":hover","color blue","div","color red","body","display flex","flex-direction row","color text-color","background-color rgb(32, 32, 32)","p","color text-hover-color","width 100%","img","height 100%",":active","#content","color bg-color","\u003e div","color orange","border none","[aria-class=\"button\"]","color green","color white","span, address","display none","h1","h2","0%","opacity 0","100%","opacity 1","animation-name appear"],"mappings":";;;;;A;CAMAA;;;AACAC;CANAC;;;AACAC;CACAC;;;AACAC;CAHAH;;;AACAC;CACAC;;;AACAC;CCkBAC;CAZAC;CACAC;CAaAC;;;AACAC;CAGAC;CAlBAJ;CACAC;CAmBAL;;;AACAO;CAGAV;;;AACAK;CArBAH;;;AACAU;CACAR;;;AACAS;CACAC;;;AACAC;CAoBAC;;;AACAb;CAGAc;;;AACAC;CACAhB;;;AACAG;CACAc;;;AACAC;CACAN;;;AACAO;CACAC;;;AACAC;CAIArB;;;AACAsB;CACAC;;;AACAC;CADAD;;;AACAC;CAHAxB;;;AACAsB;CACAC;;;AACAC;CADAD;;;AACAC;CAGAC;;;AAEAD;CADAE;;;;AACAF;CD3CAG;;;AACAC;CACAC;;;;;AACAC;CATArB;CACAJ;;;;;AAEA0B;CCuDA3B"}
`,
		`{"version":3,"sources":["testdata/test.strlt","testdata/test2.strlt"],"names":["a","color black",":hover","color blue","div","color red","body","display flex","flex-direction row","color text-color","background-color rgb(32, 32, 32)","p","color text-hover-color","width 100%","img","height 100%",":active","#content","color bg-color","\u003e div","color orange","border none","[aria-class=\"button\"]","color green","color white","span, address","display none","h1","h2","0%","opacity 0","100%","opacity 1","animation-name appear"],"mappings":"gJ,EAMAA,YACAC,SANAC,WACAC,aACAC,UACAC,SAHAH,WACAC,aACAC,UACAC,KCkBAC,aAZAC,mBACAC,wBAaAC,yBACAC,EAGAC,aAlBAJ,mBACAC,WAmBAL,yBACAO,EAGAV,UACAK,SArBAH,8BACAU,aACAR,WACAS,iBACAC,YACAC,UAoBAC,WACAb,SAGAc,UACAC,gBACAhB,UACAG,cACAc,aACAC,aACAN,YACAO,+BACAC,YACAC,WAIArB,YACAsB,SACAC,aACAC,YADAD,aACAC,SAHAxB,YACAsB,OACAC,aACAC,UADAD,aACAC,GAGAC,aAEAD,GADAE,+BACAF,GD3CAG,UACAC,KACAC,4CACAC,EATArB,aACAJ,0DAEA0B,KCuDA3B"}
`,
	}

	tree, err := codetree.FromFilelist(list)
	if err != nil {
		t.Error(err)
	}
	compiler := FromCodeTree(tree)
	for i, output := range outputs {
		cssBuilder := &MappingBuilder{}
		compiler.Render(cssBuilder, i != 1)
		mapBuilder := &strings.Builder{}
		err = cssBuilder.WriteTo(mapBuilder)
		if err != nil {
			t.Error(err)
		}
		sMap := mapBuilder.String()
		t.Log(sMap)
		if sMap != output {
			t.Error("unexpected output")
		}
	}
}
