package sql_tool

import (
	"database/sql"
	"github.com/feiin/ploto"
	"landing-page/pkg/constant"
	"landing-page/pkg/utils"
	"reflect"
)

type Table struct {
	Name string
	AIColumns []string
}

type SqlTool struct {
	Table   Table
	Mapping interface{}
	Fields  []string
	db      *sql.DB
	Type    string
}

func NewSqlTool(db *sql.DB, table Table, object interface{}, fields ...string) SqlTool {
	return SqlTool{
		db:      db,
		Table:   table,
		Mapping: object,
		Fields:  fields,
		Type:    constant.GET,
	}
}

func (t *SqlTool) SetType(queryType string) {
	t.Type = queryType
}

func (tool *SqlTool) GetQueryColumnsList() []string {
	if len(tool.Fields) > 0 {
		return tool.Fields
	}
	result := make([]string, 0)
	val := reflect.ValueOf(tool.Mapping)
	typeIgnoreAI := []string{ constant.CREATE, constant.UPDATE }
	for i := 0; i < val.Type().NumField(); i++ {
		fieldName := val.Type().Field(i).Tag.Get("json")
		if !utils.IsStringSliceContain(typeIgnoreAI, tool.Type) ||
			(utils.IsStringSliceContain(typeIgnoreAI, tool.Type) && !utils.IsStringSliceContain(tool.Table.AIColumns, fieldName)) {
			result = append(result, fieldName)
		}
	}
	return result
}

func (tool *SqlTool) Query(object interface{}, query string, args ...interface{}) error {
	rows, err := tool.db.Query(query, args...)
	if err != nil {
		return err
	}
	err = ploto.ScanResult(rows, &object)
	return err
}

func (tool *SqlTool) GetFillValue(object interface{}) []interface{} {
	val := reflect.ValueOf(object)
	result := make([]interface{}, 0)
	for i := 0; i < val.NumField(); i++ {
		result = append(result, val.Field(i))
	}
	return result
}
