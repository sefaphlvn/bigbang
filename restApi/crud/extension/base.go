package extension

import (
	"github.com/sefaphlvn/bigbang/restApi/crud"
	"github.com/sefaphlvn/bigbang/restApi/db"
)

type DBHandler crud.DbHandler

func NewExtensionHandler(db *db.MongoDB) *DBHandler {
	return &DBHandler{
		DB: db,
	}
}
