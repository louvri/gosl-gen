
package test_table 

import (
   "fmt"
   "strings"
   _sql "database/sql"
   "errors"
	"time"
    "unicode"
    helper "github.com/louvri/gosl/builder"
	db "github.com/louvri/gosl-gen/test/generated/model/test_table"
)

func buildSelection(indexFilter map[string]bool) string {
	var builder strings.Builder
	isFilterEmpty := len(indexFilter) == 0
            if isFilterEmpty || indexFilter["`id`"] {
                if builder.Len() > 0 {
                    builder.WriteString(",")
                }
                
                    builder.WriteString("COALESCE(")
                    builder.WriteString("`id`")
                    builder.WriteString(",")
                    builder.WriteString("0") // int empty
                    builder.WriteString(") AS ")
                    builder.WriteString("`id`")
                  
            }
            if isFilterEmpty || indexFilter["`field_value`"] {
                if builder.Len() > 0 {
                    builder.WriteString(",")
                }
                
                    builder.WriteString("COALESCE(")
                    builder.WriteString("`field_value`")
                    builder.WriteString(",")
                    builder.WriteString("''") // string empty
                    builder.WriteString(") AS ")
                    builder.WriteString("`field_value`")
                   
            }
	return builder.String()
}

func buildIndex(params helper.QueryParams) map[string][]map[string]interface{} {
    statements := make(map[string][]map[string]interface{})
    for key,values := range params.BetweenTime {
        if statements[key] == nil {
            statements[key] = make([]map[string]interface{},0)
        }
        statements[key] = append(statements[key],map[string]interface{}{
            "betweenTime": map[string][]time.Time {
                key: values,
            },
        })
    }
    for key,values := range params.In { 
        if statements[key] == nil {
            statements[key] = make([]map[string]interface{},0)
        }
        statements[key] = append(statements[key],map[string]interface{}{
            "in": map[string]interface{}{
                key: values,
            },
        })
    }
    for key,values := range params.Notin {
        if statements[key] == nil {
            statements[key] = make([]map[string]interface{},0)
        }
        statements[key] = append(statements[key],map[string]interface{}{
            "notIn": map[string]interface{}{
                key: values,
            },
        })
    }
    if obj,ok := params.Object.(db.Model); ok {
        for key, value := range obj.ToMap(nil) { 
            if statements[key] == nil {
                statements[key] = make([]map[string]interface{},0)
            }
            statements[key] = append(statements[key],map[string]interface{}{
                "object": map[string]interface{}{
                    key: value,
                },
            })
        }
    } else if obj,ok := params.Object.(*db.Model); ok {
         for key, value := range obj.ToMap(nil) { 
            if statements[key] == nil {
                statements[key] = make([]map[string]interface{},0)
            }
            statements[key] = append(statements[key],map[string]interface{}{
                "object": map[string]interface{}{
                    key: value,
                },
            })
        }
    } else if obj,ok := params.Object.(map[string]interface{}); ok {
         for key, value := range obj { 
            if statements[key] == nil {
                statements[key] = make([]map[string]interface{},0)
            }
            statements[key] = append(statements[key],map[string]interface{}{
                "object": map[string]interface{}{
                    key: value,
                },
            })
        }
    } 
    
    for _, condition := range params.Conditions {
        if statements[condition.Key] == nil {
            statements[condition.Key] = make([]map[string]interface{},0)
        }
        statements[condition.Key] = append(statements[condition.Key],map[string]interface{}{
            "custom": condition,
        })
    }
    return statements
}

func buildStatement(statement helper.Builder, params helper.QueryParams, priorities []string) helper.Builder {
    priorityIndex := make(map[string]bool)
    stmts := buildIndex(params) 
    build := func(data map[string]interface{}) {
        for action, parameters := range data {
            switch action {
                case "betweenTime": {
                    for key, value := range parameters.(map[string][]time.Time) { 
                        column := key 
                        timeframe := value
                        if len(timeframe) == 2 {
                              _,nWhere,_ := statement.Status()
                            if nWhere > 0 {
                                statement = statement.And()
                            }
                            statement = statement.BetweenTime(column, timeframe[0], timeframe[1])
                        }
                    }
                     
                }
                case "in": {
                    for key,value := range parameters.(map[string]interface{}) {
                        if value != nil {
                            _,nWhere,_ := statement.Status()
                            if nWhere > 0 {
                                statement = statement.And()
                            }
                            in := make(map[string]interface{})
                            in[key] = value
                            statement = statement.In(in)                
                        }
                    }
                }
                case "notIn": {
                    for key,value := range parameters.(map[string]interface{}) {
                        if value != nil {
                             _, nWhere, _ := statement.Status()
                            if nWhere > 0 {
                                statement = statement.And()
                            }
                            notin := make(map[string]interface{})
                            notin[key] = value
                            statement = statement.Not(helper.New().In(notin))
                        }
                    }
                }
                case "object": {
                    for key,value := range parameters.(map[string]interface{}) {
                        if value != nil {
                            _,nWhere,_ := statement.Status()
                            if nWhere > 0 {
                                statement = statement.And()
                            }
                            var builder strings.Builder
                            var column strings.Builder
                            for _,character := range key {
                                if character >= 'A' && character <= 'Z' {
                                    column.WriteString("_")
                                }
                                column.WriteRune(unicode.ToLower(character))
                            }
                            builder.WriteString("`") 
                            builder.WriteString(column.String())
                            builder.WriteString("`")
                            builder.WriteString(" = ?")
                            statement = statement.Statement(builder.String(),[]interface{}{value})
                        }
                    }
                }
                case "custom": {
                    cond := parameters.(helper.Condition)
                    if cond.Key != "" && cond.Value != nil {
                        _, nWhere, _ := statement.Status()
                        if nWhere > 0 {
                            statement = statement.And()
                        }
                        statement = statement.Compare([]helper.Condition{cond})
                    }
                }
            }

        }         
    }
    if priorities == nil {
       priorities = make([]string,0)
    }
    if len(priorities) == 0 {
            duplicate := make(map[string]bool)
                    if !duplicate["id"] {
                        priorities = append (priorities, "id")
                        duplicate["id"] = true
                    }
    }
    //build statement
    for _,priority := range priorities {
        priorityIndex[priority] = true
        for _,stmt := range stmts[priority] {
            build(stmt)
        }
    }
    for key,ops := range stmts {
        if priorityIndex[key] {
            continue
        }
        for _,stmt := range ops {
            build(stmt)
        }
    }
    return statement
}

func buildQuery(params helper.QueryParams, key string, operation helper.MergeOperation, tracks []interface{},priorities []string,shouldSort ...bool) helper.Builder{
    q := helper.New()
    q = q.From("test_table")
    q = buildStatement(q,params,priorities)
    if len(shouldSort) == 0 || len(shouldSort) > 0 && shouldSort[0] {
        if params.Next != 0 {
            _, nWhere, _ := q.Status()
            if nWhere > 0 {
                q = q.And()
            }
            q = q.Compare([]helper.Condition{
                {
                    Key:      key,
                    Operator: ">",
                    Value:    params.Next,
                },
            })
            q = q.Order(key, "asc")
        } else {
            if params.Page != 0 {
                q = q.Page(params.Page)
            }
            if len(params.Orderby) > 0 {
                q = q.Orders(params.Orderby)
            } else {
                q = q.Order(key, "asc")
            }
        }
    }
    if len(params.Groupby) > 0 {
         q = q.Groups(params.Groupby)
    }
    if params.Size != 0 {
        q = q.Size(params.Size)
    }
    if len(tracks) > 0 {
       switch operation {
            case helper.Identifier:
                {
                    in := make(map[string]interface{})
                    in[key] = tracks
                    if _, w, _ := q.Status(); w > 0 {
                    q = q.And()
                    }
                    q = q.Not(helper.New().In(in))
                }
            case helper.Statement:
                {
                    for i, item := range tracks {
                        //id not in (statement)
                        if _, w, _ := q.Status(); w > 0 {
                            q = q.And()
                        }
                        if statement, ok := item.(helper.Builder); ok {
                            statement.Reset("select")
                            statement.Select(key)
                            statement.Reset("from")
                            statement.From("test_table", fmt.Sprintf("%c", 'a'+i))
                            q = q.Not(helper.New().
                                Exists(
                                    statement,
                                    helper.Condition{
                                        Key: fmt.Sprintf("%s.%s",
                                            statement.Alias("test_table"),
                                            key,
                                        ),
                                        Operator: "=",
                                        Value: fmt.Sprintf("%s.`%s`",
                                            q.Alias("test_table"),
                                            key,
                                        ),
                                    },
                                ),
                            )
                        }
                    }
                }
       }
    }
    return q
}

func normalize(params *helper.QueryParams) {
	params.In = helper.ResolveColumnNameMap(params.In)
	params.Notin = helper.ResolveColumnNameMap(params.Notin)
	params.BetweenTime = helper.ResolveColumnNameMapInTime(params.BetweenTime)
	params.ColumnFilter = helper.ResolveColumnNameCollections(params.ColumnFilter)
	params.Orderby = helper.ResolveColumnNameMap(params.Orderby)
    params.Groupby = helper.ResolveColumnNameCollections(params.Groupby)
    if params.Merge != nil {
        params.Merge.Track = helper.ResolveColumnName(params.Merge.Track)
    }
	tmp := make([]helper.Condition, 0)
	for _, item := range params.Conditions {
        if !strings.Contains(item.Key,"`"){
		    item.Key = helper.ResolveColumnName(item.Key)
        }
        tmp = append(tmp, item)
	}
	params.Conditions = tmp
}

func scan(rows interface{}, filter map[string]bool) (*db.Model, error) {
    var err error
    var result db.Model
    var _row *_sql.Row
    var _rows *_sql.Rows
    var ok bool
    _row, ok = rows.(*_sql.Row)
    if !ok {
        _rows, ok = rows.(*_sql.Rows)
        if !ok {
            return nil,errors.New("wrong type of sql.row")
        }
    }
    tobeScanned := make([]interface{}, 0)
    isFilterEmpty := len(filter) == 0
            if isFilterEmpty || filter["`id`"] {
                tobeScanned = append(tobeScanned, &result.Id)
            }
            if isFilterEmpty || filter["`field_value`"] {
                tobeScanned = append(tobeScanned, &result.FieldValue)
            }
    if _rows != nil {
        err = _rows.Scan(tobeScanned...)
    } else if _row != nil {
        err = _row.Scan(tobeScanned...)
    }
    if err != nil && err != _sql.ErrNoRows {
        return nil,err
    } else if err != nil && err == _sql.ErrNoRows {
        return nil,nil
    }
    return &result,nil
}

func track(list []interface{}, data db.Model, track ...string)[]interface{} {
    if len(track) > 0 && track[0] != "" {
            if strings.ToLower(track[0]) == "id" {
                list =  append(list, data.Id)    
            }
            if strings.ToLower(track[0]) == "fieldvalue" {
                list =  append(list, data.FieldValue)    
            }
    } else {
        list = append(list, data.Id)
    }
    return list
}
