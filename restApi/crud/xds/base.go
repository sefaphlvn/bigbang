package xds

import (
	"github.com/sefaphlvn/bigbang/restapi/crud"
	"github.com/sefaphlvn/bigbang/restapi/db"
)

type DBHandler crud.DbHandler

func NewXDSHandler(db *db.MongoDB) *DBHandler {
	return &DBHandler{
		DB: db,
	}
}
