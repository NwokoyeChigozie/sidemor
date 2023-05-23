package mor

import (
	"github.com/go-playground/validator/v10"
	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/utility"
)

type Controller struct {
	Db        postgresql.Databases
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}
