package orm

import "reflect"

var QueryTag string = "Query"

type Query struct {
	table *Table
	filter
}

func (query *Query) Where(sql string, args ...interface{}) *Query {
	query.filter = filter{sql: " WHERE " + sql, args: args}
	return query
}

func (query *Query) And(sql string, args ...interface{}) *Query {
	query.and(sql, args...)
	return query
}

func (query *Query) Or(sql string, args ...interface{}) *Query {
	query.or(sql, args...)
	return query
}

func (query *Query) Get() (beans []interface{}, err error) {
	sql := "SELECT * FROM " + query.table.name + query.sql
	if query.table.showSql {
		ormLogger.D(QueryTag, "Exec# %s", sql)
	}
	rows, err := query.table.db.Query(sql, query.args...)
	if err == nil {
		columnNames, err := rows.Columns()
		if err == nil {
			for rows.Next() {
				bean := reflect.Indirect(reflect.New(query.table.beanType))
				var filedValues []interface{}
				for _, name := range columnNames {
					filedValues = append(filedValues,
						bean.FieldByName(query.table.nameMapper[name].K).Addr().Interface())
				}
				err = rows.Scan(filedValues...)
				beans = append(beans, bean)
			}
		}
	}
	return
}

func (query *Query) Count(columns ...string) (count int64, err error) {
	sql := "SELECT COUNT("
	if len(columns) == 0 {
		sql += "*"
	} else {
		sql += columns[0]
	}
	sql += ") FROM " + query.table.name
	if query.table.showSql {
		ormLogger.D(QueryTag, "Count# %s", sql)
	}
	rows, err := query.table.db.Query(sql)
	if err == nil {
		for rows.Next() {
			rows.Scan(&count)
		}
	}
	return
}
