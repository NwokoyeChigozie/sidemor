package mor

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/mor-api/internal/models"
	"github.com/vesicash/mor-api/services/mor-api"
	"github.com/vesicash/mor-api/utility"
)

func (base *Controller) GetCustomers(c *gin.Context) {
	user := models.MyIdentity
	if user == nil {
		msg := "error retrieving authenticated user"
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}
	customers, pagination, code, err := mor.GetCustomersService(c, base.ExtReq, base.Db, *user)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", customers, pagination)
	c.JSON(http.StatusOK, rd)

}
