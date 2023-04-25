package db

import (
	"sync"
	"time"

	// "corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/confighelper"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gocraft/dbr"
)

// Hold a single global connection (pooling provided by sql driver)
var sqlConnection *dbr.Connection
var connectionError error
var sqlOnce sync.Once

//GetSQLConnection to db
func GetSQLConnection() (*dbr.Connection, error) {
	sqlOnce.Do(func() {
		// create a connection db(e.g. "postgres", "mysql", or "sqlite3")
		// connection, err := dbr.Open("mysql", confighelper.GetConfig("MysqlDSN"), nil)
		connection, err := dbr.Open("mysql", "root:root@tcp(localhost:3306)/project?loc=Asia%2FKolkata&parseTime=true", nil)
		// connection, err := dbr.Open("mysql", "root:root@tcp(localhost:3306)/MAHAVAN_GIS?loc=Asia%2FKolkata&parseTime=true", nil)
		if err != nil {
			connectionError = err
		}
		connection.SetMaxIdleConns(100)
		connection.SetMaxOpenConns(5000)
		connection.SetConnMaxLifetime(3 * time.Second)
		sqlConnection = connection
	})
	return sqlConnection, connectionError
}
