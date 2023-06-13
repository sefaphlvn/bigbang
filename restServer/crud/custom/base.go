package custom

import (
	"github.com/sefaphlvn/bigbang/restServer/crud"
	"github.com/sefaphlvn/bigbang/restServer/db"
)

type DBHandler crud.DbHandler

func NewCustomHandler(db *db.MongoDB) *DBHandler {
	return &DBHandler{
		DB: db,
	}
}
