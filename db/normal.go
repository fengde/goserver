package db

import (
	"fmt"
	"server/global"

	"github.com/fengde/gocommon/storex/mysqlx"
)

type Table interface {
	TableName() string
}

// 通用插入
func Insert[T Table](t T, kv map[string]interface{}, session ...mysqlx.Session) (int64, error) {
	if len(session) > 0 {
		return session[0].Insert(t.TableName(), kv)
	}
	return global.DB.Insert(t.TableName(), kv)
}

// 通用更新
func Update[T Table](t T, kv map[string]interface{}, wherekv map[string]interface{}, session ...mysqlx.Session) (int64, error) {
	if len(session) > 0 {
		return session[0].Update(t.TableName(), kv, wherekv)
	}
	return global.DB.Update(t.TableName(), kv, wherekv)
}

func UpdateById[T Table](t T, kv map[string]interface{}, id int64, session ...mysqlx.Session) (int64, error) {
	return Update(t, kv, map[string]interface{}{"id": id}, session...)
}

// 通用删除
func Delete[T Table](t T, wherekv map[string]interface{}, session ...mysqlx.Session) (int64, error) {
	if len(session) > 0 {
		return session[0].Delete(t.TableName(), wherekv)
	}
	return global.DB.Delete(t.TableName(), wherekv)
}

func DeleteById[T Table](t T, id int64, session ...mysqlx.Session) (int64, error) {
	return Delete(t, map[string]interface{}{"id": id}, session...)
}

// 查询单条数据，根据id
func Get[T Table](t T, id int64, session ...mysqlx.Session) (*T, error) {
	var exist bool
	var err error

	if len(session) > 0 {
		exist, err = session[0].QueryOne(fmt.Sprintf(`SELECT * FROM %s WHERE id=? LIMIT 1`, t.TableName()), []interface{}{id}, &t)
	} else {
		exist, err = global.DB.QueryOne(fmt.Sprintf(`SELECT * FROM %s WHERE id=? LIMIT 1`, t.TableName()), []interface{}{id}, &t)
	}

	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, mysqlx.NotExistError
	}

	return &t, nil
}

// 通用查询
func Search[T Table](t T, sql string, args []interface{}, session ...mysqlx.Session) (rows []*T, err error) {
	if len(session) > 0 {
		err = session[0].Query(sql, args, &rows)
	} else {
		err = global.DB.Query(sql, args, &rows)
	}
	return
}
