package orm

import "database/sql"

type User struct {
	ID                 int
	Username, Password string
	Email              sql.NullString
}
