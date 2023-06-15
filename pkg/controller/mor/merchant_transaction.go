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

func (base *Controller) GetMerchantTransactions(c *gin.Context) {
	var (
		paginator = postgresql.GetPagination(c)
		req       = models.GetTransactionsRequest{
			Search:         c.Query("search"),
			CurrencyFilter: c.Query("currency"),
			Status:         c.Query("status"),
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

	user := models.MyIdentity
	if user == nil {
		msg := "error retrieving authenticated user"
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	transactions, pagination, code, err := mor.GetMerchantTransactionsService(base.ExtReq, base.Db, paginator, req, int(user.AccountID))
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", transactions, pagination)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetMerchantPayouts(c *gin.Context) {
	var (
		paginator = postgresql.GetPagination(c)
		req       = models.GetPayoutRequest{
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
	user := models.MyIdentity
	if user == nil {
		msg := "error retrieving authenticated user"
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	payouts, pagination, code, err := mor.GetMerchantPayoutsService(base.ExtReq, base.Db, paginator, req, int(user.AccountID))
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", payouts, pagination)
	c.JSON(http.StatusOK, rd)
}
