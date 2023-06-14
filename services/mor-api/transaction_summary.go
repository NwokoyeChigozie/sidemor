package mor

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/internal/models"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/services"
)

func GetMerchantTransactionsSummaryService(extReq request.ExternalRequest, db postgresql.Databases, accountID int) ([]models.TransactionSummary, int, error) {
	var (
		transaction = models.Transaction{
			MerchantID: int64(accountID),
		}
		isPaidOut = false
	)

	summaries, err := transaction.GetTransactionsSummary(db.MOR, &isPaidOut)
	if err != nil {
		return []models.TransactionSummary{}, http.StatusInternalServerError, err
	}

	summaries, err = GetTransactionsSummariesDetails(extReq, db, summaries)
	if err != nil {
		return []models.TransactionSummary{}, http.StatusInternalServerError, err
	}

	return summaries, http.StatusOK, nil
}

func GetTransactionsSummaryService(extReq request.ExternalRequest, db postgresql.Databases) ([]models.TransactionSummary, int, error) {
	var (
		transaction = models.Transaction{}
		isPaidOut   = false
	)

	summaries, err := transaction.GetTransactionsSummary(db.MOR, &isPaidOut)
	if err != nil {
		return []models.TransactionSummary{}, http.StatusInternalServerError, err
	}

	summaries, err = GetTransactionsSummariesDetails(extReq, db, summaries)
	if err != nil {
		return []models.TransactionSummary{}, http.StatusInternalServerError, err
	}

	return summaries, http.StatusOK, nil
}

func GetTransactionsSummariesDetails(extReq request.ExternalRequest, db postgresql.Databases, summaries []models.TransactionSummary) ([]models.TransactionSummary, error) {
	var (
		actualSummaries = []models.TransactionSummary{}
	)

	type summaryAndError struct {
		Summary models.TransactionSummary
		Err     error
	}

	var errs []string
	var wg sync.WaitGroup
	results := make(chan summaryAndError, len(summaries))

	for _, summary := range summaries {
		wg.Add(1)
		go func(extReq request.ExternalRequest, db postgresql.Databases, summary models.TransactionSummary, wg *sync.WaitGroup, results chan summaryAndError) {
			defer wg.Done()

			country, err := services.GetCountryByID(extReq, extReq.Logger, int(summary.CountryID))
			if err != nil {
				results <- summaryAndError{
					Summary: summary,
					Err:     err,
				}
				return
			}

			summary.Currency = country.CurrencyCode
			results <- summaryAndError{
				Summary: summary,
				Err:     err,
			}

		}(extReq, db, summary, &wg, results)
	}

	wg.Wait()
	close(results)

	// Collect the results from the channel and append them to the processedData slice.
	for result := range results {
		if result.Err != nil {
			errs = append(errs, result.Err.Error())
		} else {
			actualSummaries = append(actualSummaries, result.Summary)
		}
	}

	if len(errs) > 0 {
		extReq.Logger.Error(fmt.Sprintf("error getting mor summary details: %v", strings.Join(errs, ", ")))
	}

	return actualSummaries, nil
}
