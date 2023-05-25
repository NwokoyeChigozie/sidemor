package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/internal/models"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
)

func HandleFlutterwaveMerchantWebhook(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases, requestBody []byte) error {
	var (
		req            models.FlutterwaveWebhookRequest
		data           models.FlutterwaveWebhookRequestData
		customer       models.Customer
		paymentHistory models.PaymentHistory
		accountIDStr   = c.Param("account_id")
	)

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		return fmt.Errorf("incorrect account_id: %v", err.Error())
	}

	err = json.Unmarshal(requestBody, &req)
	if err != nil {
		return err
	}

	if req.Data != nil {
		data = *req.Data
	}

	if data.Customer != nil {
		customer.AccountID = int64(accountID)
		if data.Customer.Email != nil {
			customer.Email = *data.Customer.Email
		}
		if data.Customer.PhoneNumber != nil {
			customer.PhoneNumber = *data.Customer.PhoneNumber
		}

		if data.Customer.Name != nil {
			namesSlice := strings.Split(*data.Customer.Name, "")
			customer.Lastname = namesSlice[0]
			if len(namesSlice) > 1 {
				customer.Firstname = namesSlice[1]
			}
		}

	}

	code, err := customer.GetCustomerByAccountIDAndEmail(db.MOR)
	if err != nil {
		if code == http.StatusInternalServerError {
			return err
		}

		err := customer.CreateCustomer(db.MOR)
		if err != nil {
			return err
		}
	}

	switch req.Event {
	case "charge.completed":
		paymentHistory, err = getFlutterwavePaymentHistoryForChargeCompleted(req, &customer)
	default:
		return fmt.Errorf("event type %v, not implemented", req.Event)
	}

	if err != nil {
		return err
	}

	err = paymentHistory.CreatePaymentHistory(db.MOR)
	if err != nil {
		return err
	}

	err = customer.UpdateAllFields(db.MOR)
	if err != nil {
		return err
	}

	return nil
}

func getFlutterwavePaymentHistoryForChargeCompleted(req models.FlutterwaveWebhookRequest, customer *models.Customer) (models.PaymentHistory, error) {
	var (
		paymentHistory = models.PaymentHistory{
			CustomerID: int64(customer.ID),
		}
		data = *req.Data
	)

	if data.TxRef != nil {
		paymentHistory.Reference = *data.TxRef
	}
	if data.Narration != nil {
		paymentHistory.Description = *data.Narration
	}

	if data.Amount != nil {
		paymentHistory.Amount = *data.Amount
	}

	if data.PaymentType != nil {
		paymentHistory.PaymentMethod = *data.PaymentType
	}

	if data.CreatedAt != nil {
		t, err := time.Parse("2006-01-02T15:04:05.000Z", *data.CreatedAt)
		if err != nil {
			return models.PaymentHistory{}, fmt.Errorf("Flutterwave webhhook log error, error parsing data.DateCreated, %v, %v", *data.CreatedAt, err.Error())
		}
		customer.LastPaymentMadeAt = t
	}

	if data.Status != nil {
		if *data.Status == "successful" {
			paymentHistory.Status = paymentHistorySuccessful
		}
	}

	customer.NumberOfPayments += 1

	return paymentHistory, nil
}
