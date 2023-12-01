package iodb

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
	"ims-server/pkg/util"
	"strings"
)

type Op string

const (
	OpGt    Op = ">"
	OpLt    Op = "<"
	OpEq    Op = "="
	OpNotEq Op = "!=" // 不等于
	OpGtEq  Op = ">="
	OpLtEq  Op = "<="
	OpLike  Op = "like"
)

func (o Op) String() string {
	return string(o)
}

var OpMap = map[string]Op{
	"gt":   OpGt,
	">":    OpGt,
	"lt":   OpLt,
	"<":    OpLt,
	"eq":   OpEq,
	"=":    OpEq,
	"nte":  OpNotEq,
	"!=":   OpNotEq,
	"gte":  OpGtEq,
	">=":   OpGtEq,
	"lte":  OpLtEq,
	"<=":   OpGtEq,
	"like": OpLike,
}

// 操作符的优先级
var operatorPriority = map[string]int{
	OpGt.String():    2,
	OpLt.String():    2,
	OpEq.String():    2,
	OpNotEq.String(): 2,
	OpGtEq.String():  2,
	OpLtEq.String():  2,
	OpLike.String():  2,

	"AND": 0,
	"NOT": 0,
	"OR":  1,
}

type Order string

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"
)

func (o Order) Valid() bool {
	switch o {
	case OrderAsc, OrderDesc:
		return true
	}
	return false
}

type Express struct {
	expr string
	vars []interface{}
	src  util.Stack
}

type PageRequest struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`

	Filter string `json:"filter" form:"filter"`
	Search string `json:"search" form:"search"` // 模糊搜索（like)
	Sort   string `json:"sort" form:"sort"`

	expr  []Express
	order []Express

	filterFields util.Set[string]
}

// SetFilterFields 设置允许的查询字段，默认全不允许
func (pr *PageRequest) SetFilterFields(allowedFieldsSet util.Set[string]) {
	for k, _ := range allowedFieldsSet {
		allowedFieldsSet.Remove(k)
		allowedFieldsSet.Add(strcase.ToSnake(k))
	}
	pr.filterFields = allowedFieldsSet
}

func (pr *PageRequest) splicing(o string, expr *Express) string {
	var resp = ""

	value := expr.src.Pop()
	key := expr.src.Pop()

	switch o {
	case OpGt.String(), OpLt.String(), OpEq.String(), OpNotEq.String(), OpGtEq.String(), OpLtEq.String():
		if _, ok := pr.filterFields[key.(string)]; !ok {
			break
		}
		expr.vars = append(expr.vars, value)
		resp = fmt.Sprintf("%s %s ?", key, o)
	case OpLike.String():
		if _, ok := pr.filterFields[key.(string)]; !ok {
			break
		}
		expr.vars = append(expr.vars, "%"+value.(string)+"%")
		resp = fmt.Sprintf("%s like ?", key)
	case "OR", "AND":
		resp = fmt.Sprintf("%s %s %s", key, o, value)
		if (key == nil || key == "") && (value != nil && value != "") {
			resp = value.(string)
		} else if (key != nil && key != "") && (value == nil || value == "") {
			resp = key.(string)
		}
	case "NOT":
		expr.src.Push(key)
		if value == nil {
			break
		}
		resp = fmt.Sprintf("NOT %s", value)
	}
	return resp
}

// TODO: 能否用 pigeon 优化（https://pkg.go.dev/github.com/mna/pigeon#section-readme）
func (pr *PageRequest) handleFilter() {
	var (
		symbol = util.NewStack()
	)
	pr.Filter = strings.TrimSpace(pr.Filter)
	if pr.Filter != "" {
		expr := Express{}

		pr.Filter = "( " + pr.Filter + " )"
		parts := strings.Split(pr.Filter, " ")
		for _, part := range parts {
			// 转义比较操作
			if _, ok := OpMap[part]; ok {
				part = OpMap[part].String()
			}

			switch part {
			case "(":
				symbol.Push(part)
			case ")":
				for sym := symbol.Pop(); sym != "("; sym = symbol.Pop() {
					s := pr.splicing(sym.(string), &expr)
					expr.src.Push(s)
				}
				expr.src.Push("(" + expr.src.Pop().(string) + ")")
			case "AND", "OR", "NOT", OpGt.String(), OpLt.String(), OpEq.String(), OpNotEq.String(), OpGtEq.String(), OpLtEq.String(), OpLike.String():
				// filter 不应该能使用 'like', 转成 '='
				if part == OpLike.String() {
					part = OpEq.String()
				}
				pre := symbol.Pop().(string)

				for pre != "(" && operatorPriority[part] <= operatorPriority[pre] {
					s := pr.splicing(pre, &expr)
					expr.src.Push(s)
					pre = symbol.Pop().(string)
				}
				symbol.Push(pre)
				symbol.Push(part)
			default:
				part = strings.TrimPrefix(part, "'")
				part = strings.TrimSuffix(part, "'")
				expr.src.Push(part)
			}
		}
		expr.expr = expr.src.Pop().(string)
		expr.expr = expr.expr[1 : len(expr.expr)-1]
		// 加入 expr
		pr.expr = append(pr.expr, expr)
	}
}

func (pr *PageRequest) handleSearch() {
	var (
		symbol = util.NewStack()
	)
	pr.Search = strings.TrimSpace(pr.Search)
	if pr.Search != "" {
		expr := Express{}
		pr.Search = "( " + pr.Search + " )"
		parts := strings.Split(pr.Search, " ")
		for _, part := range parts {
			switch part {
			case "(":
				symbol.Push(part)
			case ")":
				for sym := symbol.Pop(); sym != "("; sym = symbol.Pop() {
					s := pr.splicing(sym.(string), &expr)
					expr.src.Push(s)
				}
				expr.src.Push("(" + expr.src.Pop().(string) + ")")
			case OpEq.String():
				// search 的 '=' 起的是 'like' 的作用
				part = OpLike.String()

				pre := symbol.Pop().(string)
				for pre != "(" && operatorPriority[part] <= operatorPriority[pre] {
					s := pr.splicing(pre, &expr)
					expr.src.Push(s)
					pre = symbol.Pop().(string)
				}
				symbol.Push(pre)
				symbol.Push(part)
			default:
				expr.src.Push(part)
			}
		}
		expr.expr = expr.src.Pop().(string)
		expr.expr = expr.expr[1 : len(expr.expr)-1]
		// 加入 expr
		pr.expr = append(pr.expr, expr)
	}
}

func (pr *PageRequest) handleSort() {
	pr.Sort = strings.TrimSpace(pr.Sort)
	if pr.Sort != "" {
		parts := strings.Split(pr.Search, " ")
		for _, part := range parts {
			kv := strings.Split(part, "|")
			if len(kv) != 2 {
				continue
			}
			key := kv[0]
			value := kv[1]
			switch value {
			case OrderAsc, OrderDesc:
				pr.order = append(pr.order, Express{
					expr: key + " " + value,
				})
			default:
				continue
			}
		}
	}
}

func (pr *PageRequest) Build() {
	// 处理 filter
	pr.handleFilter()
	// 处理 search
	pr.handleSearch()
	// 处理 Sort
	pr.handleSort()
}

func (pr *PageRequest) ToFilterDB(db *gorm.DB) *gorm.DB {
	filterDB := db
	for _, v := range pr.expr {
		filterDB = filterDB.Where(v.expr, v.vars...)
	}
	for _, v := range pr.order {
		filterDB = filterDB.Order(v.expr)
	}
	return filterDB
}

func (pr *PageRequest) ToPageDB(db *gorm.DB) *gorm.DB {
	resp := pr.ToFilterDB(db)
	resp = resp.Limit(pr.Size).Offset(pr.Page * pr.Size)
	return resp
}
