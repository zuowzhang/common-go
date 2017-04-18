package orm

var DeleteTag string = "Delete"

type Delete struct {
	table *Table
	filter
}

func (delete *Delete) Where(sql string, args ...interface{}) *Delete {
	delete.filter = filter{sql: " WHERE " + sql, args: args}
	return delete
}

func (delete *Delete) And(sql string, args ...interface{}) *Delete {
	delete.and(sql, args...)
	return delete
}

func (delete *Delete) Or(sql string, args ...interface{}) *Delete {
	delete.or(sql, args...)
	return delete
}

func (delete *Delete) Exec() (int64, error) {
	sql := "DELETE FROM " + delete.table.name + delete.sql
	if delete.table.showSql {
		ormLogger.D(DeleteTag, "Exec# %s", sql)
	}
	res, err := delete.table.db.Exec(sql, delete.args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
