package sqlite3

import (
	"testing"
)

func TestSqlite3(t *testing.T) {
	err := Init("sqlite3", "./me.db")
	t.Log(err)
	create := `CREATE TABLE "CacheFiles" (
		"ID" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"FILEPATH" TEXT,
		"FILESTATUS" TEXT,
		"FILEREADTIME" TEXT,
		"FILEMODTIME" TEXT,
		"FILENAME" TEXT,
		"FILESIZA" TEXT
	  );`
	err = Exec(create)
	t.Log(err)

	inserts := `insert into T1(id,name) values (?,?);`

	value1 := make([][]string, 0)
	v1 := make([]string, 0)
	v1 = []string{"1", "Tom1"}
	value1 = append(value1, v1)

	v1 = []string{"2", "Tom2"}
	value1 = append(value1, v1)

	v1 = []string{"3", "Tom3"}
	value1 = append(value1, v1)

	err = Inserts(inserts, value1)
	t.Log(err)

	query := `select * from T1;`
	res, err := Querys(query)
	t.Log(err)
	t.Log(res)

}
