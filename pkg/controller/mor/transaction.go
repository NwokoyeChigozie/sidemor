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

func (base *Controller) RecordTransaction(c *gin.Context) {
	var (
		req models.RecordTransactionRequest
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

	vr := postgresql.ValidateRequestM{Logger: base.Logger, Test: base.ExtReq.Test}
	err = vr.ValidateRequest(req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	transaction, code, err := mor.RecordTransactionService(base.ExtReq, base.Db, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "successfully created", transaction)
	c.JSON(http.StatusCreated, rd)

}

func (base *Controller) GetTransaction(c *gin.Context) {
	var (
		id = c.Param("id")
	)

	transactionID, err := strconv.Atoi(id)
	if err != nil {
		msg := fmt.Sprintf("invalid id: %v", err.Error())
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", msg, fmt.Errorf(msg), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	transaction, code, err := mor.GetTransactionService(base.ExtReq, base.Db, transactionID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successfully created", transaction)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) GetTransactions(c *gin.Context) {
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

	transactions, pagination, code, err := mor.GetTransactionsService(base.ExtReq, base.Db, paginator, req)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "successful", transactions, pagination)
	c.JSON(http.StatusOK, rd)

}
