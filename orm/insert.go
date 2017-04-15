package orm

import "reflect"

var InsertTag string = "Insert"

type Insert struct {
	table *Table
	beans []interface{}
}

func (insert *Insert)Values(bean ...interface{}) *Insert {
	insert.beans = append(insert.beans, bean...)
	return insert
}

func (insert *Insert)Exec() (affected int64, err error) {
	sql := "INSERT INTO " + insert.table.name + " VALUES"
	var args []interface{}
	for idx, bean := range insert.beans {
		if idx != 0 {
			sql += ", "
		}
		var first bool = true
		for _, v := range insert.table.nameMapper {
			if first {
				first = false
				sql += "("
			} else {
				sql += ", "
			}
			sql += "?"
			value := reflect.ValueOf(bean)
			args = append(args, value.FieldByName(v.K).Interface())
		}
		sql += ")"
	}
	if insert.table.showSql {
		ormLogger.D(InsertTag, "Exec# %s; args%v", sql, args)
	}
	res, err := insert.table.db.Exec(sql, args...)
	if err == nil {
		affected, err = res.RowsAffected()
	}
	return
}
