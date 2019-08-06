package sqlite3

import (
	"database/sql"
	"reflect"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var sqlite3 *sql.DB
var once sync.Once

//Init mssql driver init
func Init(drivername, datasourcename string) (err error) {
	once.Do(func() {
		sqlite3, err = sql.Open(drivername, datasourcename)
	})
	return sqlite3.Ping()
}

func Exec(constr string) error {
	stmt, err := sqlite3.Prepare(constr)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

// Querys  query more
func Querys(constr string) (res [][]string, err error) {
	res = make([][]string, 0)
	stmt, err := sqlite3.Prepare(constr)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil || rows == nil {
		return
	}
	cols, err := rows.Columns()
	if err != nil {
		return
	}
	for rows.Next() {
		arr := make([]interface{}, len(cols))
		re := make([]string, len(cols))
		for i := 0; i < len(cols); i++ {
			arr[i] = new(sql.NullString)
		}
		err = rows.Scan(arr...)
		if err != nil {
			return
		}

		for i := 0; i < len(arr); i++ {
			arrtemp := reflect.ValueOf(arr[i])
			arrtem := arrtemp.Interface().(*sql.NullString)
			re[i] = arrtem.String
		}
		res = append(res, re)
	}
	defer rows.Close()
	return
}

//Inserts insert more rows
func Inserts(constr string, res [][]string) (err error) {
	conn, err := sqlite3.Begin()

	for c := 0; c < len(res); c++ {
		re := make([]interface{}, len(res[c]))
		for i := 0; i < len(res[c]); i++ {
			re[i] = res[c][i]
		}
		stmt, err := sqlite3.Prepare(constr)
		if err != nil {
			conn.Rollback()
		}
		defer stmt.Close()
		_, err = stmt.Exec(re...)
		if err != nil {
			conn.Rollback()
		}
		defer stmt.Close()
	}
	conn.Commit()
	return
}
