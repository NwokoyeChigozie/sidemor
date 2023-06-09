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

func (base *Controller) GetPayout(c *gin.Context) {
	var (
		id = c.Param("id")
	)

	payoutID, err := strconv.Atoi(id)
	if err != nil {
		msg := fmt.Sprintf("invalid id: %v", err.Error())
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	payout, code, err := mor.GetPayoutService(base.ExtReq, base.Db, payoutID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", payout)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetPayouts(c *gin.Context) {
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

	payouts, pagination, code, err := mor.GetPayoutsService(base.ExtReq, base.Db, paginator, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", payouts, pagination)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) PayOutToWallets(c *gin.Context) {
	var (
		req models.PayoutToWalletRequest
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

	msg, code, err := mor.PayOutToWalletsService(base.ExtReq, base.Db, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, msg, nil)
	c.JSON(http.StatusOK, rd)

}
