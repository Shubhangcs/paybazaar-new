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

type BBPSInterface interface {
	CreatePostpaidMobileRecharge(echo.Context) error
	GetPostpaidMobileRechargeBalance(echo.Context) (*models.GetPostpaidMobileRechargeBillFetchAPIResponseModel, error)
	GetAllPostpaidMobileRecharge(echo.Context) ([]models.GetPostpaidMobileRechargeHistoryResponseModel, error)
	GetPostpaidMobileRechargeByRetailerID(echo.Context) ([]models.GetPostpaidMobileRechargeHistoryResponseModel, error)
	CreateElectricityBillPayment(echo.Context) error
	GetAllElectricityOperators(echo.Context) ([]models.GetElectricityOperatorResponseModel, error)
}

type bbpsRepository struct {
	db *database.Database
}

func NewBBPSRepository(db *database.Database) *bbpsRepository {
	return &bbpsRepository{
		db,
	}
}

func (bp *bbpsRepository) CreatePostpaidMobileRecharge(c echo.Context) error {
	var req models.CreatePostpaidMobileRechargeAPIRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	req.PartnerRequestID = uuid.NewString()

	apiUrl := `https://v2a.rechargkit.biz/recharge/postpaid`
	reqBody, err := json.Marshal(map[string]any{
		"mobile_no":          req.MobileNumber,
		"partner_request_id": req.PartnerRequestID,
		"operator_code":      req.OperatorCode,
		"circle":             req.OperatorCircle,
		"amount":             req.Amount,
		"recharge_type":      1,
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

	var res models.GetPostpaidMobileRechargeAPIResponseModel
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return err
	}

	fmt.Println(string(respBytes))

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	var status string
	if res.Status == 1 {
		status = "SUCCESS"
		if err := bp.db.CreatePostpaidMobileRechargeSuccessOrPendingQuery(ctx, req, res, status); err != nil {
			return err
		}
		return nil
	}

	if res.Status == 2 {
		status = "PENDING"
		if err := bp.db.CreatePostpaidMobileRechargeSuccessOrPendingQuery(ctx, req, res, status); err != nil {
			return err
		}
		return nil
	}

	if res.Status == 3 {
		status = "FAILED"
		if err := bp.db.CreatePostpaidMobileRechargeFailureQuery(ctx, req, res); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("invalid status from recharge kit")
}

func (bp *bbpsRepository) GetPostpaidMobileRechargeBalance(c echo.Context) (*models.GetPostpaidMobileRechargeBillFetchAPIResponseModel, error) {
	var req models.GetPostpaidMobileRechargeBillFetchAPIRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}
	apiUrl := fmt.Sprintf("https://v2a.rechargkit.biz/recharge/postPaidBillFetch?mobile_no=%s&operator_code=%d", req.MobileNumber, req.OperatorCode)

	apiRequest, err := http.NewRequest(
		http.MethodGet,
		apiUrl,
		nil,
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
	fmt.Println(string(respBytes))

	var res models.GetPostpaidMobileRechargeBillFetchAPIResponseModel
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (bp *bbpsRepository) GetAllPostpaidMobileRecharge(c echo.Context) ([]models.GetPostpaidMobileRechargeHistoryResponseModel, error) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	limit, offset := parsePagination(c)
	return bp.db.GetAllPostpaidMobileRechargeQuery(ctx, limit, offset)
}

func (bp *bbpsRepository) GetPostpaidMobileRechargeByRetailerID(c echo.Context) ([]models.GetPostpaidMobileRechargeHistoryResponseModel, error) {
	var retailerId = c.Param("retailer_id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	limit, offset := parsePagination(c)
	return bp.db.GetPostpaidMobileRechargeByRetailerIDQuery(ctx, retailerId, limit, offset)
}

func (bp *bbpsRepository) CreateElectricityBillPayment(c echo.Context) error {
	var req models.CreateElectricityBillPaymentRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	req.PartnerRequestID = uuid.NewString()

	apiUrl := `https://v2a.rechargkit.biz/recharge/billpayment`
	reqBody, err := json.Marshal(map[string]any{
		"p1":                 req.CustomerID,
		"partner_request_id": req.PartnerRequestID,
		"operator_code":      req.OperatorCode,
		"customer_email":     req.CustomerEmail,
		"amount":             req.Amount,
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

	var res models.GetElectricityBillPaymentAPIResponseModel
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return err
	}

	fmt.Println(string(respBytes))

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	var status string
	if res.Status == 1 {
		status = "SUCCESS"
		if err := bp.db.CreateElectricityBillPaymentSuccessOrPendingQuery(ctx, req, res, status); err != nil {
			return err
		}
		return nil
	}

	if res.Status == 2 {
		status = "PENDING"
		if err := bp.db.CreateElectricityBillPaymentSuccessOrPendingQuery(ctx, req, res, status); err != nil {
			return err
		}
		return nil
	}

	if res.Status == 3 {
		status = "FAILED"
		if err := bp.db.CreateElectricityBillPaymentFailureQuery(ctx, req, res); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("invalid status from recharge kit")
}

func (bp *bbpsRepository) GetAllElectricityOperators(c echo.Context) ([]models.GetElectricityOperatorResponseModel, error) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()
	return bp.db.GetElectricityOperatorsQuery(ctx)
}
