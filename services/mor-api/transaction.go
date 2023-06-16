package mor

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/internal/models"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/services"
)

func RecordTransactionService(extReq request.ExternalRequest, db postgresql.Databases, req models.RecordTransactionRequest) (models.Transaction, int, error) {
	var (
		transaction = models.Transaction{}
	)

	checkTime := int(time.Now().Add(336 * time.Hour).Unix())
	if req.TransactionCreatedAt > checkTime {
		return models.Transaction{}, http.StatusBadRequest, fmt.Errorf("invalid timestamp, time must not be more that 2 weeks after today")
	}

	transaction.MerchantID = req.AccountID
	transaction.Reference = req.Reference
	transaction.Description = req.Description
	transaction.CountryID = int64(req.Country)
	transaction.Amount = req.Amount
	transaction.TaxFee = req.TaxFee
	transaction.ProcessingFee = req.ProcessingFee
	transaction.TransactionDate = time.Unix(int64(req.TransactionCreatedAt), 0)
	transaction.Status = models.TransactionSuccessful
	err := transaction.CreateTransaction(db.MOR)
	if err != nil {
		return transaction, http.StatusInternalServerError, err
	}

	transaction, err = GetMorTransactionDetails(extReq, db, transaction)
	if err != nil {
		return transaction, http.StatusInternalServerError, err
	}

	return transaction, http.StatusOK, nil
}

func GetTransactionService(extReq request.ExternalRequest, db postgresql.Databases, transactionID int) (models.Transaction, int, error) {
	var (
		transaction = models.Transaction{ID: uint(transactionID)}
	)

	code, err := transaction.GetTransactionByID(db.MOR)
	if err != nil {
		return transaction, code, err
	}

	transaction, err = GetMorTransactionDetails(extReq, db, transaction)
	if err != nil {
		return transaction, http.StatusInternalServerError, err
	}

	return transaction, http.StatusOK, nil
}

func GetTransactionsService(extReq request.ExternalRequest, db postgresql.Databases, paginator postgresql.Pagination, req models.GetTransactionsRequest) ([]models.Transaction, postgresql.PaginationResponse, int, error) {
	var (
		transaction = models.Transaction{}
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

func GetMorTransactionDetails(extReq request.ExternalRequest, db postgresql.Databases, transaction models.Transaction) (models.Transaction, error) {
	user, err := services.GetUserWithAccountID(extReq, int(transaction.MerchantID))
	if err != nil {
		return transaction, err
	}

	country, err := services.GetCountryByID(extReq, extReq.Logger, int(transaction.CountryID))
	if err != nil {
		return transaction, err
	}

	customer := models.Customer{ID: uint(transaction.CustomerID)}
	code, err := customer.GetCustomerByID(db.MOR)
	if code == http.StatusInternalServerError {
		return transaction, err
	}

	transaction.CustomerName = customer.Lastname + " " + customer.Firstname
	transaction.MerchantName = user.Lastname + " " + user.Firstname
	transaction.MerchantEmail = user.EmailAddress
	transaction.Country = country.Name
	transaction.Currency = country.CurrencyCode

	return transaction, nil
}

func GetMorTransactionsDetails(extReq request.ExternalRequest, db postgresql.Databases, transactions []models.Transaction) ([]models.Transaction, error) {

	type transactionAndError struct {
		Transaction models.Transaction
		Err         error
	}

	var newTransactions []models.Transaction
	var errs []string
	var wg sync.WaitGroup
	results := make(chan transactionAndError, len(transactions))

	// Loop through the data slice and spawn a goroutine for each item.
	for _, transaction := range transactions {
		wg.Add(1)
		go func(extReq request.ExternalRequest, db postgresql.Databases, transaction models.Transaction, wg *sync.WaitGroup, results chan transactionAndError) {
			defer wg.Done()
			transaction, err := GetMorTransactionDetails(extReq, db, transaction)
			results <- transactionAndError{
				Transaction: transaction,
				Err:         err,
			}

		}(extReq, db, transaction, &wg, results)
	}

	wg.Wait()
	close(results)

	// Collect the results from the channel and append them to the processedData slice.
	for result := range results {
		if result.Err != nil {
			errs = append(errs, result.Err.Error())
		} else {
			newTransactions = append(newTransactions, result.Transaction)
		}
	}

	if len(errs) > 0 {
		extReq.Logger.Error(fmt.Sprintf("error getting mor transaction details: %v", strings.Join(errs, ", ")))
	}

	return newTransactions, nil
}
