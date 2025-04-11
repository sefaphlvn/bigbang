package custom

import (
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sefaphlvn/bigbang/rest/crud"
)

type AppHandler crud.Application

func NewCustomHandler(context *db.AppContext) *AppHandler {
	return &AppHandler{
		Context: context,
	}
}
