package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type DTHRechargeInterface interface {
	GetAllDTHOperators(echo.Context) ([]models.GetDTHOperatorsResponseModel, error)
	CreateDTHRecharge(echo.Context) error
	GetAllDTHRecharges(echo.Context) ([]models.GetDTHRechargeHistoryResponseModel, error)
	GetDTHRechargesByRetailerID(echo.Context) ([]models.GetDTHRechargeHistoryResponseModel, error)
}

type dthRechargeRepository struct {
	db *database.Database
}

func NewDTHRechargeRepository(db *database.Database) *dthRechargeRepository {
	return &dthRechargeRepository{
		db,
	}
}

func (drr *dthRechargeRepository) GetAllDTHOperators(c echo.Context) ([]models.GetDTHOperatorsResponseModel, error) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	return drr.db.GetAllDTHOperatorsQuery(ctx)
}

func (drr *dthRechargeRepository) CreateDTHRecharge(c echo.Context) error {
	var req models.CreateDTHRechargeRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	req.PartnerRequestID = uuid.NewString()

	fmt.Println(req)

	apiUrl := `https://v2a.rechargkit.biz/recharge/dth`
	reqBody, err := json.Marshal(map[string]any{
		"customer_id":        req.CustomerID,
		"operator_code":      req.OperatorCode,
		"amount":             req.Amount,
		"partner_request_id": req.PartnerRequestID,
	})
	if err != nil {
		return err
	}

	apiRequest, err := http.NewRequest(
		http.MethodPost,
		apiUrl,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return err
	}

	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set("Authorization", "Bearer "+os.Getenv("RKIT_API_TOKEN"))

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(apiRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var apiResponse struct {
		Status  int    `json:"status"`
		Message string `json:"msg"`
	}
	if err := json.Unmarshal(respBytes, &apiResponse); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()

	fmt.Println(apiResponse)
	if apiResponse.Status == 1 {
		req.Status = "SUCCESS"
		if err := drr.db.CreateDTHRechargeSuccessOrPendingQuery(ctx, req); err != nil {
			return err
		}
		return nil
	}

	if apiResponse.Status == 2 {
		req.Status = "PENDING"
		if err := drr.db.CreateDTHRechargeSuccessOrPendingQuery(ctx, req); err != nil {
			return err
		}
		return nil
	}

	if apiResponse.Status == 3 {
		req.Status = "FAILED"
		if err := drr.db.CreateDTHRechargeFailedQuery(ctx, req); err != nil {
			return err
		}
		return fmt.Errorf("failed to recharge: %s", apiResponse.Message)
	}
	return fmt.Errorf("invalid status from recharge kit")

}

func (drr *dthRechargeRepository) GetAllDTHRecharges(c echo.Context) ([]models.GetDTHRechargeHistoryResponseModel, error) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()
	limit, offset := parsePagination(c)
	return drr.db.GetAllDTHRechargesQuery(ctx, limit, offset)
}

func (drr *dthRechargeRepository) GetDTHRechargesByRetailerID(c echo.Context) ([]models.GetDTHRechargeHistoryResponseModel, error) {
	var retailerID = c.Param("retailer_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()
	limit, offset := parsePagination(c)
	return drr.db.GetDTHRechargesByRetailerIDQuery(ctx, retailerID, limit, offset)
}