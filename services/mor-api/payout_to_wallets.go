package mor

import (
	"fmt"
	"net/http"

	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/internal/models"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/services"
	"github.com/vesicash/mor-api/utility"
)

func PayOutToWalletsService(extReq request.ExternalRequest, db postgresql.Databases, req models.PayoutToWalletRequest) (string, int, error) {
	if len(req.Merchants) > 0 && len(req.Merchants) < 10 {
		code, err := PayoutToUsers(extReq, db, req.Merchants)
		if err != nil {
			return "", code, err
		}

		return "payout successful", http.StatusOK, nil
	} else if len(req.Merchants) > 10 {
		go PayoutToUsers(extReq, db, req.Merchants)
		return "payout started", http.StatusOK, nil
	}
	return "", http.StatusOK, nil
}

func PayoutToUsers(extReq request.ExternalRequest, db postgresql.Databases, accountIDs []int) (int, error) {
	for _, accountID := range accountIDs {
		code, err := PayoutToUser(extReq, db, accountID)
		if err != nil {
			return code, err
		}
	}
	return http.StatusOK, nil
}

func PayoutToUser(extReq request.ExternalRequest, db postgresql.Databases, accountID int) (int, error) {
	var (
		IsPaidOut                 = false
		currenciesTransactionsMap map[int][]models.Transaction
	)

	_, err := services.GetUserWithAccountID(extReq, accountID)
	if err != nil {
		return http.StatusBadRequest, err
	}

	transaction := models.Transaction{MerchantID: int64(accountID), IsPaidOut: IsPaidOut, Status: models.TransactionSuccessful}
	transactions, err := transaction.GetTransactionsAll(db.MOR, &IsPaidOut)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, tx := range transactions {
		txs := currenciesTransactionsMap[int(tx.CountryID)]
		txs = append(txs, tx)
		currenciesTransactionsMap[int(tx.CountryID)] = txs
	}

	for countryID, trxs := range currenciesTransactionsMap {
		country, err := services.GetCountryByID(extReq, extReq.Logger, countryID)
		if err != nil {
			extReq.Logger.Error(fmt.Sprintf("error getting country with id %v: %v", countryID, err.Error()))
			continue
		}

		var totalAmount float64

		for _, trx := range trxs {
			totalAmount += trx.Amount
		}

		payout := models.Payout{
			MerchantID: int64(accountID),
			Reference:  utility.RandomString(25),
			Amount:     totalAmount,
			CountryID:  int64(countryID),
			Status:     models.TransactionSuccessful,
		}

		_, err = services.CreditWallet(extReq, db, totalAmount, country.CurrencyCode, accountID, false, "no", "yes", "")
		if err != nil {
			payout.Status = models.TransactionFailed
			extReq.Logger.Error(fmt.Sprintf("error crediting mor wallet %v, amount %v: %v", country.CurrencyCode, totalAmount, err.Error()))
		}

		err = payout.CreatePayout(db.MOR)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		if payout.Status != models.TransactionSuccessful {
			continue
		}

		for _, trx := range trxs {
			trx.IsPaidOut = true
			trx.PayoutID = int64(payout.ID)
			err := trx.UpdateAllFields(db.MOR)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}

	}

	return http.StatusOK, nil
}
