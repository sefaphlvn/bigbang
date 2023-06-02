package xds

import (
	"github.com/sefaphlvn/bigbang/restApi/crud"
	"github.com/sefaphlvn/bigbang/restApi/db"
)

type DBHandler crud.DbHandler

func NewXDSHandler(db *db.MongoDB) *DBHandler {
	return &DBHandler{
		DB: db,
	}
}
