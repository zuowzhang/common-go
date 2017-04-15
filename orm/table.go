package orm

import (
	"database/sql"
	"reflect"
	"strings"
	"bytes"
	"common/log"
	_ "github.com/go-sql-driver/mysql"
)

var TableTag string = "Table"

var ormLogger *log.Logger = log.NewLogger()

type pair struct {
	K, V string
}

type Table struct {
	db *sql.DB
	name string
	beanType reflect.Type
	nameMapper map[string]pair
	showSql bool
}

func getTagNameMapper(beanType reflect.Type) map[string]pair {
	mapper := make(map[string]pair)
	for idx := 0; idx < beanType.NumField(); {
		f := beanType.Field(idx)
		idx++
		ignore := f.Tag.Get("ignore")
		if strings.ToLower(ignore) == "true" {
			continue
		}
		name := f.Tag.Get("name")
		if name == "" {
			buffer := bytes.NewBuffer([]byte{})
			for idx, r := range f.Name {
				if r >= 'A' && r <= 'Z' {
					if idx != 0 {
						buffer.WriteByte('_')
					}
					r += 'a' - 'A'
				}
				buffer.WriteRune(r)
			}
			name = string(buffer.Bytes())
		}
		dbType := f.Tag.Get("type")
		if dbType == "" {
			switch f.Type.Kind() {
			case reflect.Bool:
				dbType = "TINYINT(1)"
			case reflect.Int8, reflect.Uint8:
				dbType = "TINYINT"
			case reflect.Float32:
				dbType = "FLOAT"
			case reflect.Float64:
				dbType = "DOUBLE"
			case reflect.Int16, reflect.Uint16:
				dbType = "SMALLINT"
			case reflect.Int, reflect.Uint32:
				dbType = "INTEGER"
			case reflect.Int64, reflect.Uint64:
				dbType = "BIGINT"
			default:
				dbType = "VARCHAR(255)"
			}
		}
		mapper[name] = pair{K:f.Name, V:dbType}
	}
	return mapper
}

func NewTable(beanType reflect.Type, driver, dataSourceName string) (*Table, error) {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Table{db:db,
		name:strings.ToLower(beanType.Name()),
		beanType:beanType,
		nameMapper:getTagNameMapper(beanType)}, nil
}

func (table *Table)Close() error {
	return table.db.Close()
}

func (table *Table)ShowSql(show bool) *Table {
	table.showSql = show
	return table
}

func (table *Table)CreateDbTable() error {
	buffer := bytes.NewBufferString("CREATE TABLE IF NOT EXISTS ")
	buffer.WriteString(table.name)
	buffer.WriteString("(")
	idx := 0
	for key, value := range table.nameMapper {
		if idx != 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(key + " " + value.V)
		idx++
	}
	buffer.WriteString(")")
	sql := string(buffer.Bytes())
	if table.showSql {
		ormLogger.D(TableTag, "create table# %s", sql)
	}
	_, err := table.db.Exec(sql)
	return err
}

func (table *Table)NewQuery() *Query {
	return &Query{table:table}
}

func (table *Table)NewInsert() *Insert {
	return &Insert{table:table}
}

func (table *Table)NewDelete() *Delete {
	return &Delete{table:table}
}

func (table *Table)NewUpdate() *Update {
	return &Update{table:table}
}