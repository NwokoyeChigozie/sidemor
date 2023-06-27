package mor

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/mor-api/internal/models"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/services/mor-api"
	"github.com/vesicash/mor-api/utility"
)

func (base *Controller) RequestWithdrawal(c *gin.Context) {
	var (
		req models.RequestWithdrawalRequest
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

	code, err := mor.RequestWithdrawalService(base.ExtReq, base.Db, *user, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", nil)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetWithdrawals(c *gin.Context) {
	var (
		paginator = postgresql.GetPagination(c)
		req       = models.GetWithdrawalRequest{
			Search:         c.Query("search"),
			CurrencyFilter: c.Query("currency"),
			Status:         c.Query("status"),
		}
	)

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

	withdrawals, pagination, code, err := mor.GetWithdrawalsService(base.ExtReq, base.Db, paginator, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", withdrawals, pagination)
	c.JSON(http.StatusOK, rd)

}
func (base *Controller) CompleteWithdrawal(c *gin.Context) {
	var (
		withdrawalIDStr = c.Param("withdrawal_id")
	)

	withdrawalID, err := strconv.Atoi(withdrawalIDStr)
	if err != nil {
		err = fmt.Errorf("invalid id: %v", err)
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	code, err := mor.CompleteWithdrawalService(base.ExtReq, base.Db, withdrawalID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", nil)
	c.JSON(http.StatusOK, rd)

}

// CompleteWithdrawal
