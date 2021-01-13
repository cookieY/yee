package yee

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

var test = `
type GenCodeVal struct {
	${ATTRIBUTE}
}
`

func jk(t string) {
	test = strings.Replace(test, "${ATTRIBUTE}", t, -1)
	f, err := os.OpenFile("koala.go", os.O_WRONLY&os.O_CREATE, 0666)
	if err != nil {
		log.Println(err.Error())
	}
	_, err = f.Write([]byte(test))
	if err != nil {
		log.Println(err.Error())
	}
	f.Close()
}

func TestNew(t *testing.T) {
	c := []map[string]string{{"a": "int"}, {"k": "string"}}
	//l := ""
	for _, i := range c {
		fmt.Println(i)
		//fmt.Println(c[j])
	}
}

var exprCase = []expr{{Name: "Username", Expr: "username =?", TP: "string"}, {Name: "Age", Expr: "age > ?", TP: "int"}}

func TestGenQueryExpr(t *testing.T) {

	GenQueryExpr(exprCase)
}

func TestGenerateRestfulAPI(t *testing.T) {
	k := GenCodeVal{
		Package:   "manage",
		QueryExpr: exprCase,
		Page:      "20",
		Modal:     "modal.core_account",
	}
	fmt.Println(GenerateRestfulAPI(k))
}
