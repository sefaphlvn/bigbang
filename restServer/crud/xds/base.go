package xds

import (
	"github.com/sefaphlvn/bigbang/restServer/crud"
	"github.com/sefaphlvn/bigbang/restServer/db"
)

type DBHandler crud.DbHandler

func NewXDSHandler(db *db.MongoDB) *DBHandler {
	return &DBHandler{
		DB: db,
	}
}
