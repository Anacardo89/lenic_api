package orm

import "database/sql"

var (
	Da DataAccess
)

type DataAccess struct {
	Db *sql.DB
}
