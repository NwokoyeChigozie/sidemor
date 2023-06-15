package mor

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/mor-api/internal/models"
	"github.com/vesicash/mor-api/services/mor-api"
	"github.com/vesicash/mor-api/utility"
)

func (base *Controller) SaveSettings(c *gin.Context) {
	var (
		req models.SaveSettingsRequest
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	user := models.MyIdentity
	if user == nil {
		msg := "error retrieving authenticated user"
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	settings, code, err := mor.SaveSettingsService(base.ExtReq, base.Db, *user, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successfully saved", settings)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetSettings(c *gin.Context) {

	user := models.MyIdentity
	if user == nil {
		msg := "error retrieving authenticated user"
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	settings, code, err := mor.GetSettingsService(base.ExtReq, base.Db, *user)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", settings)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) EnableOrDisablePaymentMethods(c *gin.Context) {
	var (
		req    models.EnableOrDisablePaymentMethodsRequest
		action = c.Param("action")
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	user := models.MyIdentity
	if user == nil {
		msg := "error retrieving authenticated user"
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	settings, code, err := mor.EnableOrDisablePaymentMethodsService(base.ExtReq, base.Db, *user, action, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successfully saved", settings)
	c.JSON(http.StatusOK, rd)

}
func (base *Controller) AddRemoveOrGetWallets(c *gin.Context) {
	var (
		req    models.AddRemoveOrGetWalletsRequest
		action = c.Param("action")
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	user := models.MyIdentity
	if user == nil {
		msg := "error retrieving authenticated user"
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	settings, code, err := mor.AddRemoveOrGetWalletsService(base.ExtReq, base.Db, *user, action, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", settings)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetVerificationSettings(c *gin.Context) {
	var (
		paginator = postgresql.GetPagination(c)
		req       = models.GetSettingsRequest{
			Search: c.Query("search"),
			Status: c.Query("status"),
		}
	)

	// search=transaction_id, filters:=wallet_type(currency), status, date range

	if c.Query("from") != "" {
		from, err := strconv.Atoi(c.Query("from"))
		if err != nil {
			msg := fmt.Sprintf("invalid from: %v, must be timestamp interger", err.Error())
			rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", msg, fmt.Errorf(msg), nil)
			c.JSON(http.StatusBadRequest, rd)
			return
		}
		req.FromTime = from
	}

	if c.Query("to") != "" {
		to, err := strconv.Atoi(c.Query("to"))
		if err != nil {
			msg := fmt.Sprintf("invalid to: %v, must be timestamp interger", err.Error())
			rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", msg, fmt.Errorf(msg), nil)
			c.JSON(http.StatusBadRequest, rd)
			return
		}
		req.ToTime = to
	}

	settings, pagination, code, err := mor.GetVerificationSettingsService(base.ExtReq, base.Db, paginator, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", settings, pagination)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) UpdateDocumentStatus(c *gin.Context) {
	var (
		req        models.UpdateDocumentStatusRequest
		settingsID = c.Param("id")
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	id, err := strconv.Atoi(settingsID)
	if err != nil {
		return
	}

	_, code, err := mor.UpdateDocumentStatusService(base.Db, id, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	if code != http.StatusOK {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "Document status updated", nil)
	c.JSON(http.StatusOK, rd)

}
