package orm

import (
	"reflect"
	"testing"
	"os"
)

type User struct {
	Name        string `name:"name"`
	Address     string `ignore:"true"`
	Age         int
	PhoneNumber string
}

func TestGetTagNameMapper(t *testing.T) {
	t.Logf("%#v\n", getTagNameMapper(reflect.TypeOf(User{})))
}

func getTable() (*Table, error) {
	table, err := NewTable(reflect.TypeOf(User{}), "mysql", "root:123456@/uc?charset=utf8")
	if err != nil {
		return nil, err
	}
	table.ShowSql(true)
	return table, nil
}

func TestNewTable(t *testing.T) {
	table, err := getTable()
	if err != nil {
		t.Fatal(err)
	}
	err = table.CreateDbTable()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("success")
}

func TestQuery_Get(t *testing.T) {
	table, err := getTable()
	if err != nil {
		t.Fatal(err)
	}
	beans, err := table.NewQuery().
		Where("age = ?", 18).
		Or("name = ?", "zhangsan").
		Get()

	if err != nil {
		t.Fatal(err)
	}
	t.Log(beans...)
}

func TestQuery_Count(t *testing.T) {
	table, err := getTable()
	if err != nil {
		t.Fatal(err)
	}
	count, err := table.NewQuery().Count("name")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("count = ", count)
}

func TestInsert_Exec(t *testing.T) {
	table, err := getTable()
	if err != nil {
		t.Fatal(err)
	}
	affected, err := table.NewInsert().Values(User{Age: 18}, User{Name: "zhangsan"}).Exec()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("affected = ", affected)
}

func TestDelete_Exec(t *testing.T) {
	table, err := getTable()
	if err != nil {
		t.Fatal(err)
	}
	affected, err := table.NewDelete().Where("age = 18").Exec()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(affected)
}

func TestUpdate_Exec(t *testing.T) {
	table, err := getTable()
	if err != nil {
		t.Fatal(err)
	}
	affected, err := table.NewUpdate().Where("name = ?", "zhangsan").
		Values(User{Age: 19, PhoneNumber: "13866666666", Name: "zhangsan"}).
		Exec()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(affected)
}
