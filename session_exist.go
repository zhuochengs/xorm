// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"errors"
	"fmt"
	builder "github.com/zhuochengs/xorm-builder"
	"reflect"

	"github.com/zhuochengs/xorm-core"
)

// Exist returns true if the record exist otherwise return false
func (session *Session) Exist(bean ...interface{}) (bool, error) {
	if session.isAutoClose {
		defer session.Close()
	}

	var sqlStr string
	var args []interface{}
	var err error

	if session.statement.RawSQL == "" {
		if len(bean) == 0 {
			tableName := session.statement.TableName()
			if len(tableName) <= 0 {
				return false, ErrTableNotFound
			}

			if session.statement.cond.IsValid() {
				condSQL, condArgs, err := builder.ToSQL(session.statement.cond)
				if err != nil {
					return false, err
				}

				if session.engine.dialect.DBType() == core.MSSQL {
					sqlStr = fmt.Sprintf("SELECT top 1 * FROM %s WHERE %s", tableName, condSQL)
				} else {
					sqlStr = fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT 1", tableName, condSQL)
				}
				args = condArgs
			} else {
				if session.engine.dialect.DBType() == core.MSSQL {
					sqlStr = fmt.Sprintf("SELECT top 1 * FROM %s", tableName)
				} else {
					sqlStr = fmt.Sprintf("SELECT * FROM %s LIMIT 1", tableName)
				}
				args = []interface{}{}
			}
		} else {
			beanValue := reflect.ValueOf(bean[0])
			if beanValue.Kind() != reflect.Ptr {
				return false, errors.New("needs a pointer")
			}

			if beanValue.Elem().Kind() == reflect.Struct {
				if err := session.statement.setRefBean(bean[0]); err != nil {
					return false, err
				}
			}

			if len(session.statement.TableName()) <= 0 {
				return false, ErrTableNotFound
			}
			session.statement.Limit(1)
			sqlStr, args, err = session.statement.genGetSQL(bean[0])
			if err != nil {
				return false, err
			}
		}
	} else {
		sqlStr = session.statement.RawSQL
		args = session.statement.RawParams
	}

	rows, err := session.queryRows(sqlStr, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}
