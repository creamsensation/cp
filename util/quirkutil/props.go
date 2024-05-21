package quirkutil

import (
	"fmt"
	"strings"
	
	"github.com/creamsensation/cp"
	"github.com/creamsensation/quirk"
	"github.com/creamsensation/util"
)

type Props struct {
	Fulltext string
	Offset   int
	Limit    int
	Order    []string
	Columns  Columns
}

type Columns struct {
	Fulltext []string
	Order    map[string]string
}

func CreateProps(c cp.Ctx) Props {
	var props Props
	c.Parse().MustQuery(Fulltext, &props.Fulltext)
	c.Parse().MustQuery(Offset, &props.Offset)
	c.Parse().MustQuery(Limit, &props.Limit)
	var order string
	c.Parse().MustQuery(Order, &order)
	if len(order) > 0 {
		props.Order = append(props.Order, order)
	}
	return props
}

func UseProps(q *quirk.Quirk, props Props) {
	useFulltextProps(q, props.Fulltext, props.Columns.Fulltext)
	useOrderProps(q, props.Order, props.Columns.Order)
	useOffsetProps(q, props.Offset)
	useLimitProps(q, props.Limit)
}

func useFulltextProps(q *quirk.Quirk, fulltext string, columns []string) {
	if len(fulltext) == 0 || len(columns) == 0 {
		return
	}
	startWord := "WHERE"
	if q.WhereExists() {
		startWord = "AND"
	}
	conditions := make([]string, len(columns))
	args := make(quirk.Map)
	for i := range conditions {
		name := fmt.Sprintf("fulltext%d", i+1)
		args[name] = quirk.CreateTsQuery(fulltext)
		conditions[i] = columns[i] + " @@ to_tsquery(@" + name + ")"
	}
	q.Q(startWord+` (`+strings.Join(conditions, " OR ")+`)`, args)
}

func useOrderProps(q *quirk.Quirk, order []string, columns map[string]string) {
	if len(order) == 0 || len(columns) == 0 {
		return
	}
	r := make([]string, 0)
	for _, o := range order {
		if !strings.Contains(o, ":") {
			continue
		}
		parts := strings.Split(o, ":")
		if len(parts) < 2 || parts[1] == "" {
			continue
		}
		column, ok := columns[parts[0]]
		if !ok {
			continue
		}
		r = append(r, util.EscapeString(column)+" "+util.EscapeString(strings.ToUpper(parts[1])))
	}
	if len(r) == 0 {
		return
	}
	q.Q(`ORDER BY ` + strings.Join(r, ","))
}

func useOffsetProps(q *quirk.Quirk, offset int) {
	q.Q(`OFFSET @offset`, quirk.Map{Offset: offset})
}

func useLimitProps(q *quirk.Quirk, limit int) {
	if limit == -1 {
		return
	}
	if limit == 0 {
		limit = DefaultLimit
	}
	q.Q(`LIMIT @limit`, quirk.Map{Limit: limit})
}
