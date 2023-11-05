package xds

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/restServer/crud"
)

type DBHandler crud.DbHandler

func NewXDSHandler(db *db.MongoDB) *DBHandler {
	return &DBHandler{
		DB: db,
	}
}
