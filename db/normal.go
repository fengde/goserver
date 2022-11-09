package db

import (
	"fmt"
	"goserver/global"
	"strings"

	"github.com/fengde/gocommon/storex/mysqlx"
	"github.com/fengde/gocommon/taskx"
)

type Table interface {
	TableName() string
}

// 通用插入
func Insert[T Table](t T, kv map[string]any, session ...mysqlx.Session) (int64, error) {
	if len(session) > 0 {
		return session[0].Insert(t.TableName(), kv)
	}
	return global.DB.Insert(t.TableName(), kv)
}

// 通用更新
func Update[T Table](t T, kv map[string]any, wherekv map[string]any, session ...mysqlx.Session) (int64, error) {
	if len(session) > 0 {
		return session[0].Update(t.TableName(), kv, wherekv)
	}
	return global.DB.Update(t.TableName(), kv, wherekv)
}

func UpdateById[T Table](t T, kv map[string]any, id int64, session ...mysqlx.Session) (int64, error) {
	return Update(t, kv, map[string]any{"id": id}, session...)
}

// 通用删除
func Delete[T Table](t T, wherekv map[string]any, session ...mysqlx.Session) (int64, error) {
	if len(session) > 0 {
		return session[0].Delete(t.TableName(), wherekv)
	}
	return global.DB.Delete(t.TableName(), wherekv)
}

func DeleteById[T Table](t T, id int64, session ...mysqlx.Session) (int64, error) {
	return Delete(t, map[string]any{"id": id}, session...)
}

// 查询单条数据，根据id
func Get[T Table](t T, id int64, session ...mysqlx.Session) (*T, error) {
	var exist bool
	var err error

	if len(session) > 0 {
		exist, err = session[0].QueryOne(fmt.Sprintf(`SELECT * FROM %s WHERE id=? LIMIT 1`, t.TableName()), []any{id}, &t)
	} else {
		exist, err = global.DB.QueryOne(fmt.Sprintf(`SELECT * FROM %s WHERE id=? LIMIT 1`, t.TableName()), []any{id}, &t)
	}

	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, mysqlx.NotExistError
	}

	return &t, nil
}

type ColumnHelper struct {
	Column string
	Symbol string
	Value  any
}

type OrderHelper struct {
	Column string
	Symbol string
}

// 日常分页查询封装
func SearchHelper[T Table](t T, columns []ColumnHelper, orders []OrderHelper, pageIndex int64, pageSize int64, session ...mysqlx.Session) (total int64, rows []*T, err error) {
	var where = `1=1`
	var args = []any{}

	for _, c := range columns {
		where += fmt.Sprintf(` AND %s %s ?`, c.Column, c.Symbol)
		args = append(args, c.Value)
	}

	stg := taskx.NewSerialTaskGroup(
		func() error {
			// 查询总数
			sql := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, t.TableName(), where)
			if len(session) > 0 {
				_, err = session[0].QueryOne(sql, args, &total)
			} else {
				_, err = global.DB.QueryOne(sql, args, &total)
			}
			return err
		}, func() error {
			// 查询记录
			var sql string
			if len(orders) > 0 {
				var orderbys []string
				for _, order := range orders {
					orderbys = append(orderbys, order.Column+" "+order.Symbol)
				}
				sql = fmt.Sprintf(`SELECT * FROM %s WHERE %s ORDER BY %s LIMIT %d, %d`, t.TableName(), where, strings.Join(orderbys, ", "), (pageIndex-1)*pageSize, pageSize)
			} else {
				sql = fmt.Sprintf(`SELECT * FROM %s WHERE %s LIMIT %d, %d`, t.TableName(), where, (pageIndex-1)*pageSize, pageSize)
			}

			if len(session) > 0 {
				err = session[0].Query(sql, args, &rows)
			} else {
				err = global.DB.Query(sql, args, &rows)
			}
			return err
		})

	stg.Run()

	return
}
