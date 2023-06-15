package mor

import (
	"net/http"
	"strings"

	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/internal/models"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/services"
)

func GetMerchantTransactionsService(extReq request.ExternalRequest, db postgresql.Databases, paginator postgresql.Pagination, req models.GetTransactionsRequest, merchantID int) ([]models.Transaction, postgresql.PaginationResponse, int, error) {
	var (
		transaction = models.Transaction{MerchantID: int64(merchantID)}
		isPaidOut   = false
	)

	if req.CurrencyFilter != "" {
		country, _ := services.GetCountryByCurrency(extReq, extReq.Logger, strings.ToUpper(req.CurrencyFilter))
		transaction.CountryID = int64(country.ID)
	}

	if req.Status != "" {
		transaction.Status = models.TransactionStatus(req.Status)
	}

	transactions, pagination, err := transaction.GetTransactions(db.MOR, paginator, req.Search, req.FromTime, req.ToTime, &isPaidOut)
	if err != nil {
		return transactions, pagination, http.StatusInternalServerError, err
	}

	transactions, err = GetMorTransactionsDetails(extReq, db, transactions)
	if err != nil {
		return transactions, pagination, http.StatusInternalServerError, err
	}

	return transactions, pagination, http.StatusOK, nil
}

func GetMerchantPayoutsService(extReq request.ExternalRequest, db postgresql.Databases, paginator postgresql.Pagination, req models.GetPayoutRequest, merchantID int) ([]models.Payout, postgresql.PaginationResponse, int, error) {
	var (
		payout = models.Payout{MerchantID: int64(merchantID)}
	)

	if req.CurrencyFilter != "" {
		country, _ := services.GetCountryByCurrency(extReq, extReq.Logger, strings.ToUpper(req.CurrencyFilter))
		payout.CountryID = int64(country.ID)
	}

	if req.Status != "" {
		payout.Status = models.TransactionStatus(req.Status)
	}

	payouts, pagination, err := payout.GetPayouts(db.MOR, paginator, req.Search, req.FromTime, req.ToTime)
	if err != nil {
		return payouts, pagination, http.StatusInternalServerError, err
	}

	payouts, err = GetMorPayoutsDetails(extReq, db, payouts)
	if err != nil {
		return payouts, pagination, http.StatusInternalServerError, err
	}

	return payouts, pagination, http.StatusOK, nil
}
