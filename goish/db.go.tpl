package main

import "database/sql"
import "time"
import "errors"

//import "log"
import _ "github.com/lib/pq"

type DBCon struct {
	db *sql.DB
}

func (r DBCon) Type() string {
	return "database-connection"
}

func (r DBCon) String() string {
	return r.Type()
}

func (r DBCon) Lib() Searchable {
	return protoDBCon
}

type DBTx struct {
	tx *sql.Tx
}

func (r DBTx) Type() string {
	return "database-transaction"
}

func (r DBTx) String() string {
	return r.Type()
}

func (r DBTx) Lib() Searchable {
	return protoDBTx
}

type DBRes struct {
	rows *sql.Rows
}

func (r DBRes) Type() string {
	return "database-result"
}

func (r DBRes) String() string {
	return r.Type()
}

func (r DBRes) Lib() Searchable {
	return protoDBRes
}

var libSQL *Map
var protoDBCon *Map
var protoDBTx *Map
var protoDBRes *Map
var Nsql Any

func init() {
	onInit = append(onInit, func() {

		sqlArgs := func(aa *Args) []interface{} {
			args := []interface{}{}
			for i := 2; i < aa.len(); i++ {
				arg := aa.get(i)
				switch arg.(type) {
				case Int:
					args = append(args, int64(arg.(Int)))
				case Dec:
					args = append(args, float64(arg.(Dec)))
				default:
					if arg == nil {
						args = append(args, nil)
					} else {
						args = append(args, totext(arg))
					}
				}
			}
			return args
		}

		protoDBCon = NewMap(MapData{

			Text("close"): Func(func(vm *VM, aa *Args) *Args {
				con := aa.get(0).(DBCon)
				vm.da(aa)
				return join(vm, NewStatus(con.db.Close()))
			}),

			Text("begin"): Func(func(vm *VM, aa *Args) *Args {
				con := aa.get(0).(DBCon)
				vm.da(aa)
				tx, err := con.db.Begin()
				return join(vm, NewStatus(err), DBTx{tx})
			}),

			Text("execute"): Func(func(vm *VM, aa *Args) *Args {
				con := aa.get(0).(DBCon)
				sql := totext(aa.get(1))
				args := sqlArgs(aa)
				vm.da(aa)
				result, err := con.db.Exec(sql, args...)
				if err != nil {
					return join(vm, NewStatus(err))
				}
				id, _ := result.LastInsertId()
				affected, _ := result.RowsAffected()
				return join(vm, NewStatus(err), Int(id), Int(affected))
			}),

			Text("query"): Func(func(vm *VM, aa *Args) *Args {
				con := aa.get(0).(DBCon)
				sql := totext(aa.get(1))
				args := sqlArgs(aa)
				vm.da(aa)
				rows, err := con.db.Query(sql, args...)
				if err != nil {
					return join(vm, NewStatus(err))
				}
				return join(vm, NewStatus(err), DBRes{rows})
			}),
		})
		protoDBCon.meta = protoDef

		protoDBTx = NewMap(MapData{

			Text("commit"): Func(func(vm *VM, aa *Args) *Args {
				tx := aa.get(0).(DBTx)
				vm.da(aa)
				return join(vm, NewStatus(tx.tx.Commit()))
			}),

			Text("rollback"): Func(func(vm *VM, aa *Args) *Args {
				tx := aa.get(0).(DBTx)
				vm.da(aa)
				return join(vm, NewStatus(tx.tx.Rollback()))
			}),

			Text("execute"): Func(func(vm *VM, aa *Args) *Args {
				tx := aa.get(0).(DBTx)
				sql := totext(aa.get(1))
				args := sqlArgs(aa)
				vm.da(aa)
				result, err := tx.tx.Exec(sql, args...)
				if err != nil {
					return join(vm, NewStatus(err))
				}
				id, _ := result.LastInsertId()
				affected, _ := result.RowsAffected()
				return join(vm, NewStatus(err), Int(id), Int(affected))
			}),

			Text("query"): Func(func(vm *VM, aa *Args) *Args {
				tx := aa.get(0).(DBTx)
				sql := totext(aa.get(1))
				args := sqlArgs(aa)
				vm.da(aa)
				rows, err := tx.tx.Query(sql, args...)
				if err != nil {
					return join(vm, NewStatus(err))
				}
				return join(vm, NewStatus(err), DBRes{rows})
			}),
		})
		protoDBTx.meta = protoDef

		protoDBRes = NewMap(MapData{

			Text("close"): Func(func(vm *VM, aa *Args) *Args {
				res := aa.get(0).(DBRes)
				vm.da(aa)
				return join(vm, NewStatus(res.rows.Close()))
			}),

			Text("iterate"): Func(func(vm *VM, aa *Args) *Args {

				res := aa.get(0).(DBRes)
				cols, _ := res.rows.Columns()
				vm.da(aa)

				return join(vm, Func(func(vm *VM, aa *Args) *Args {

					if res.rows.Next() {

						columns := make([]interface{}, len(cols))
						columnPointers := make([]interface{}, len(cols))
						for i, _ := range columns {
							columnPointers[i] = &columns[i]
						}

						if err := res.rows.Scan(columnPointers...); err != nil {
							return join(vm, NewStatus(err))
						}

						m := NewMap(MapData{})
						for i, colName := range cols {
							val := columnPointers[i].(*interface{})
							switch (*val).(type) {
							case int64:
								m.Set(Text(colName), Int((*val).(int64)))
							case float64:
								m.Set(Text(colName), Dec((*val).(float64)))
							case bool:
								m.Set(Text(colName), Bool((*val).(bool)))
							case []byte:
								m.Set(Text(colName), Blob((*val).([]byte)))
							case time.Time:
								m.Set(Text(colName), Instant((*val).(time.Time)))
							case string:
								m.Set(Text(colName), Text((*val).(string)))
							default:
								if *val == nil {
									m.Set(Text(colName), nil)
								} else {
									m.Set(Text(colName), NewStatus(errors.New("unknown type")))
								}
							}
						}
						return join(vm, m)
					}

					return join(vm, nil)
				}))
			}),
		})
		protoDBRes.meta = protoDef

		libSQL = NewMap(MapData{

			Text("db"): protoDBCon,
			Text("rs"): protoDBRes,
			Text("tx"): protoDBTx,

			Text("open"): Func(func(vm *VM, aa *Args) *Args {
				driver := totext(aa.get(0))
				connStr := totext(aa.get(1))
				vm.da(aa)
				db, err := sql.Open(driver, connStr)
				if err != nil {
					return join(vm, NewStatus(err))
				}
				return join(vm, NewStatus(err), DBCon{db})
			}),
		})
		Nsql = libSQL
	})
}
