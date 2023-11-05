package extension

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/restServer/crud"
)

type DBHandler crud.DbHandler

func NewExtensionHandler(db *db.MongoDB) *DBHandler {
	return &DBHandler{
		DB: db,
	}
}
