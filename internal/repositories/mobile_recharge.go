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

type MobileRechargeInterface interface {
	CreateMobileRecharge(echo.Context) error
	GetAllMobileRechargeCircles(echo.Context) ([]models.GetMobileRechargeCircleResponseModel, error)
	GetAllMobileRechargeOperators(echo.Context) ([]models.GetMobileRechargeOperatorsResponseModel, error)
	GetAllPlansBasedOnCircleAndOperator(echo.Context) (any, error)
	GetAllMobileRecharges(echo.Context) ([]models.GetMobileRechargeHistoryResponseModel, error)
	GetMobileRechargesByRetailerID(echo.Context) ([]models.GetMobileRechargeHistoryResponseModel, error)
}

type mobileRechargeRepository struct {
	db *database.Database
}

func NewMobileRechargeRepository(db *database.Database) *mobileRechargeRepository {
	return &mobileRechargeRepository{
		db,
	}
}

func (mrr *mobileRechargeRepository) CreateMobileRecharge(c echo.Context) error {
	var req models.CreateMobileRechargeRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	req.PartnerRequestID = uuid.NewString()
	apiUrl := `https://v2bapi.rechargkit.biz/recharge/prepaid`
	reqBody, err := json.Marshal(map[string]any{
		"mobile_no":      req.MobileNumber,
		"operator_code":  req.OperatorCode,
		"amount":         req.Amount,
		"partner_req_id": req.PartnerRequestID,
		"circle":         req.CircleCode,
		"recharge_type":  1,
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
	if apiResponse.Status == 1 {
		req.Status = "SUCCESS"
		mrr.db.CreateMobileRechargeSuccessOrPendingQuery(ctx, req)
		return nil
	}

	if apiResponse.Status == 2 {
		req.Status = "PENDING"
		mrr.db.CreateMobileRechargeSuccessOrPendingQuery(ctx, req)
		return nil
	}

	if apiResponse.Status == 3 {
		req.Status = "FAILED"
		mrr.db.CreateMobileRechargeFailedQuery(ctx, req)
		return fmt.Errorf("failed to recharge: %s", apiResponse.Message)
	}
	return fmt.Errorf("invalid status from recharge kit")
}

func (mrr *mobileRechargeRepository) GetAllMobileRechargeCircles(c echo.Context) ([]models.GetMobileRechargeCircleResponseModel, error) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()
	return mrr.db.GetAllMobileRechargeCirclesQuery(ctx)
}

func (mrr *mobileRechargeRepository) GetAllMobileRechargeOperators(c echo.Context) ([]models.GetMobileRechargeOperatorsResponseModel, error) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()
	return mrr.db.GetAllMobileRechargeOperatorsQuery(ctx)
}

func (mrr *mobileRechargeRepository) GetAllPlansBasedOnCircleAndOperator(c echo.Context) (any, error) {
	var req models.GetMobileRechargePlansRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}
	apiUrl := `https://v2bapi.rechargkit.biz/recharge/prepaidPlanFetch`
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	apiRequest, err := http.NewRequest(
		http.MethodPost,
		apiUrl,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, err
	}
	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set("Authorization", "Bearer "+os.Getenv("RKIT_API_TOKEN"))

	client := &http.Client{Timeout: 20 * time.Second}

	resp, err := client.Do(apiRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res models.GetMobileRechargePlansResponseModel
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, err
	}

	if res.Error == 1 {
		return nil, fmt.Errorf("failed to fetch plan: %s", res.Message)
	}

	return res, nil
}

func (mrr *mobileRechargeRepository) GetAllMobileRecharges(c echo.Context) ([]models.GetMobileRechargeHistoryResponseModel, error) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()
	limit, offset := parsePagination(c)
	return mrr.db.GetAllMobileRechargesQuery(ctx, limit, offset)
}

func (mrr *mobileRechargeRepository) GetMobileRechargesByRetailerID(c echo.Context) ([]models.GetMobileRechargeHistoryResponseModel, error) {
	var retailerID = c.Param("retailer_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()
	limit, offset := parsePagination(c)
	return mrr.db.GetMobileRechargesByRetailerIDQuery(ctx, retailerID, limit, offset)
}
