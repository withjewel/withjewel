package model

import (
	"database/sql"
	"withjewel/jewel/jedb"

	"github.com/jinzhu/gorm"
)

type Artical struct {
	gorm.Model
	User    User
	UserID  int
	Source  sql.NullString
	Content sql.NullString
}

func CreateArtical(userID int, source string) {
	db := jedb.Accquire()
	defer jedb.Revert(db)

	newArtical := &Artical{UserID: userID, Source: sql.NullString{String: source}}
	db.Create(newArtical)
}
