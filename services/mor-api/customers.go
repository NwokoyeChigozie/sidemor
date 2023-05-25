package mor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/internal/models"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/services"
)

func GetCustomersService(c *gin.Context, extReq request.ExternalRequest, db postgresql.Databases) ([]models.Customer, postgresql.PaginationResponse, int, error) {
	var (
		paginator = postgresql.GetPagination(c)
		customer  = models.Customer{}
		search    = c.Query("search")
	)

	accessToken, err := services.GetAccessTokenByKeyFromRequest(extReq, c)
	if err != nil {
		return []models.Customer{}, postgresql.PaginationResponse{}, http.StatusInternalServerError, err
	}

	customer.AccountID = int64(accessToken.AccountID)

	customers, pagination, err := customer.GetCustomers(db.MOR, paginator, search)
	if err != nil {
		return customers, pagination, http.StatusInternalServerError, err
	}

	return customers, pagination, http.StatusOK, nil
}
