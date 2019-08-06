module github.com/seaiiok/collector

go 1.12

replace github.com/seaiiok/snet/snet.v4/clients => /snet/snet.v4/clients

require (
	github.com/boltdb/bolt v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/julienschmidt/httprouter v1.2.0
	github.com/mattn/go-sqlite3 v1.11.0
)
