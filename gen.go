package yee

import (
	"fmt"
	"strings"
)

// gen.go used to generate Restful code

const PREFIX = `
package ${PACKAGE}

import (
    "net/http"
	"github.com/cookieY/yee"
	"github.com/jinzhu/gorm"
)

${TP}

func Paging(page interface{}, total int) (start int, end int) {
	start = i*total - total
	end = total
	return
}

func Fetch${PACKAGE}(c yee.Context) (err error) {
	u := new(FinderPrefix)
	if err = c.Bind(u); err != nil {
		c.Logger().Error(err.Error())
		return
	}

	var order []${MODAL}
	
	start, end := lib.Paging(u.Page, ${PAGE})

	if u.Find.Valve {
		model.DB().Model(&${MODAL}{}).
			Scopes(
				${QUERY_EXPR}
			).Count(&pg).Order("id desc").Offset(start).Limit(end).Find(&order)
	} else {
		model.DB().Model(&model.CoreSqlOrder{}).Count(&pg).Order("id desc").Offset(start).Limit(end).Find(&order)
	}

    return c.JSON(http.StatusOK, map[string]interface{}{"data": order, "page": pg})
}

${QUERYFUNC}

`

var QueryExprPrefix = `
func AccordingTo${EXPR_NAME}(val string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(${QUERY_EXPR},val)
	}
}
`

var FinderPrefix = `
type FinderPrefix struct {
    Valve   bool   // 自行添加tag
	${FINDER_EXPR}
}
`

type expr struct {
	Name string `json:"name"`
	Expr string `json:"expr"`
	TP   string `json:"tp"`
}

type GenCodeVal struct {
	Flag      string `json:"flag"`       // 根据哪个字段进行CURD
	Package   string `json:"package"`    // 项目名,根据项目名生成package name
	QueryExpr []expr `json:"query_expr"` // 查询条件
	Page      string `json:"page"`       //分页大小
	Modal     string `json:"modal"`
}

func GenerateRestfulAPI(GenCodeVal GenCodeVal) string {
	empty := strings.Replace(PREFIX, "${PACKAGE}", GenCodeVal.Package, -1)
	empty = strings.Replace(empty, "${MODAL}", GenCodeVal.Modal, -1)
	empty = strings.Replace(empty, "${PAGE}", GenCodeVal.Page, -1)
	f, s,l := GenQueryExpr(GenCodeVal.QueryExpr)
	empty = strings.Replace(empty, "${QUERYFUNC}", f, -1)
	empty = strings.Replace(empty, "${TP}", s, -1)
	empty = strings.Replace(empty, "${QUERY_EXPR}", l, -1)
	return empty
}

func GenQueryExpr(QueryExpr []expr) (string, string, string) {
	funcEmpty := ""
	structEmpty := ""
	exprList := ""
	for _, i := range QueryExpr {
		tmpText := ""
		tmpText = strings.Replace(QueryExprPrefix, "${EXPR_NAME}", i.Name, -1)
		tmpText = strings.Replace(tmpText, "${QUERY_EXPR}", i.Expr, -1)
		funcEmpty += tmpText + "\n"
		structEmpty += fmt.Sprintf("%s    %s \n    ", i.Name, i.TP)
		exprList += fmt.Sprintf("AccordingTo%s(u.%s),\n                ", i.Name, strings.ToLower(i.Name))
	}
	structEmpty = strings.Replace(FinderPrefix, "${FINDER_EXPR}", structEmpty, -1)
	return funcEmpty, structEmpty, exprList
}
