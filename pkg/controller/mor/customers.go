package mor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/mor-api/services/mor-api"
	"github.com/vesicash/mor-api/utility"
)

func (base *Controller) GetCustomers(c *gin.Context) {

	customers, pagination, code, err := mor.GetCustomersService(c, base.ExtReq, base.Db)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", customers, pagination)
	c.JSON(http.StatusOK, rd)

}
