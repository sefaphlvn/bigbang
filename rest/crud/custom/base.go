package custom

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/rest/crud"
)

type DBHandler crud.DbHandler

func NewCustomHandler(db *db.WTF) *DBHandler {
	return &DBHandler{
		DB: db,
	}
}
