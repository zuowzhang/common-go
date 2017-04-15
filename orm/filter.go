package orm

type filter struct {
	sql  string
	args []interface{}
}

func (filter *filter)and(sql string, args ...interface{}) {
	filter.sql += " AND " + sql
	filter.args = append(filter.args, args...)
}

func (filter *filter)or(sql string, args ...interface{}) {
	filter.sql += " OR " + sql
	filter.args = append(filter.args, args...)
}
