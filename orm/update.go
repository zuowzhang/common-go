package orm

import (
	"errors"
	"reflect"
)

var UpdateTag string = "Update"

type Update struct {
	table *Table
	filter
	columns []string
	bean    interface{}
}

func (update *Update) Where(sql string, args ...interface{}) *Update {
	update.filter = filter{sql: " WHERE " + sql, args: args}
	return update
}

func (update *Update) And(sql string, args ...interface{}) *Update {
	update.and(sql, args...)
	return update
}

func (update *Update) Or(sql string, args ...interface{}) *Update {
	update.or(sql, args...)
	return update
}

func (update *Update) Columns(column ...string) *Update {
	update.columns = append(update.columns, column...)
	return update
}

func (update *Update) Values(bean interface{}) *Update {
	update.bean = bean
	return update
}

func (update *Update) Exec() (int64, error) {
	sql := "UPDATE " + update.table.name + " SET "
	if len(update.columns) == 0 {
		for column := range update.table.nameMapper {
			update.columns = append(update.columns, column)
		}
	}
	var args []interface{}
	for idx, column := range update.columns {
		if idx > 0 {
			sql += ", "
		}
		sql += column + " = ?"
		if pair, ok := update.table.nameMapper[column]; ok {
			args = append(args, reflect.ValueOf(update.bean).FieldByName(pair.K).Interface())
		} else {
			return 0, errors.New("can not find column " + column)
		}
	}
	sql += update.sql
	args = append(args, update.args...)
	if update.table.showSql {
		ormLogger.D(UpdateTag, "Exec# %s; %v", sql, args)
	}
	res, err := update.table.db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
