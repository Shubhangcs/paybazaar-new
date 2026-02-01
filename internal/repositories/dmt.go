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

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type DMTInterface interface {
	CheckDMTWalletExists(echo.Context) (*models.CheckDMTWalletExistsResponseModel, error)
	CreateDMTWallet(echo.Context) (*models.CreateDMTWalletResponseModel, error)
	VerifyDMTWalletCreation(echo.Context) (*models.CreateDMTWalletVerifyResponseModel, error)
	CreateDMTBeneficiary(echo.Context) (*models.CreateDMTBeneficiaryResponseModel, error)
	GetDMTBeneficieries(echo.Context) (*models.GetDMTBeneficiariesResponseModel, error)
	DeleteDMTBeneficiary(echo.Context) (*models.DeleteDMTBeneficiaryResponseModel, error)
	VerifyDMTBeneficiaryDelete(echo.Context) (*models.DeleteDMTBeneficiaryVerificationResponseModel, error)
}

type dmtRepository struct {
	db *database.Database
}

func NewDMTRepository(db *database.Database) *dmtRepository {
	return &dmtRepository{
		db,
	}
}

func (dr *dmtRepository) CheckDMTWalletExists(c echo.Context) (*models.CheckDMTWalletExistsResponseModel, error) {
	var req models.CheckDMTWalletExistsRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}
	apiUrl := `https://v2bapi.rechargkit.biz/moneytransfer/checkWalletExist`
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

	var res models.CheckDMTWalletExistsResponseModel
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, err
	}
	fmt.Println(res)
	if res.Error != 0 {
		return nil, fmt.Errorf("error from rechargekit: %s", res.Description)
	}
	return &res, nil
}

func (dr *dmtRepository) CreateDMTWallet(c echo.Context) (*models.CreateDMTWalletResponseModel, error) {
	var req models.CreateDMTWalletRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}
	apiUrl := `https://v2bapi.rechargkit.biz/moneytransfer/createWalletRequest`
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

	var res models.CreateDMTWalletResponseModel
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, err
	}
	fmt.Println(res)
	if res.Error != 0 {
		return nil, fmt.Errorf("error from rechargekit: %s", res.Description)
	}
	return &res, nil
}

func (dr *dmtRepository) VerifyDMTWalletCreation(c echo.Context) (*models.CreateDMTWalletVerifyResponseModel, error) {
	var req models.CreateDMTWalletVerifyRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*20)
	defer cancel()
	if err := dr.db.GetUserDetailsForDMTWalletCreation(ctx, &req); err != nil {
		return nil, err
	}
	apiUrl := `https://v2bapi.rechargkit.biz/moneytransfer/verifyOtp`
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

	var res models.CreateDMTWalletVerifyResponseModel
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, err
	}
	if res.Error != 0 {
		return nil, fmt.Errorf("error from rechargekit: %s", res.Description)
	}
	return &res, nil
}

func (dr *dmtRepository) CreateDMTBeneficiary(c echo.Context) (*models.CreateDMTBeneficiaryResponseModel, error) {
	var req models.CreateDMTBeneficiaryRequestModel
	if err := bindAndValidate(c, &req); err != nil {
		return nil, err
	}
	apiUrl := `https://v2bapi.rechargkit.biz/addBeneficiaryRequest`
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

	var res models.CreateDMTBeneficiaryResponseModel
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, err
	}
	if res.Error != 0 {
		return nil, fmt.Errorf("error from rechargekit: %s", res.Description)
	}
	return &res, nil
}

func (dr *dmtRepository) GetDMTBeneficieries(c echo.Context) (*models.GetDMTBeneficiariesResponseModel, error) {
	var req models.GetDMTBeneficiariesRequestModel
	if err := bindAndValidate(c,&req); err != nil {
		return nil, err
	}
	apiUrl := `https://v2bapi.rechargkit.biz/getUserDetails`
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

	var res models.GetDMTBeneficiariesResponseModel
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, err
	}
	if res.Error != 0 {
		return nil, fmt.Errorf("error from rechargekit: %s", res.Description)
	}
	return &res, nil
}

func (dr *dmtRepository) DeleteDMTBeneficiary(c echo.Context) (*models.DeleteDMTBeneficiaryResponseModel, error) {
	var req models.DeleteDMTBeneficiaryRequestModel
	if err := bindAndValidate(c,&req); err != nil {
		return nil, err
	}
	apiUrl := `https://v2bapi.rechargkit.biz/deleteBeneficiaryRequest`
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

	var res models.DeleteDMTBeneficiaryResponseModel
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, err
	}
	if res.Error != 0 {
		return nil, fmt.Errorf("error from rechargekit: %s", res.Description)
	}
	return &res, nil
}

func (dr *dmtRepository) VerifyDMTBeneficiaryDelete(c echo.Context) (*models.DeleteDMTBeneficiaryVerificationResponseModel, error) {
	var req models.DeleteDMTBeneficiaryVerificationRequestModel
	if err := bindAndValidate(c,&req); err != nil {
		return nil, err
	}
	apiUrl := `https://v2bapi.rechargkit.biz/moneytransfer/confirmDeleteBeneficiary`
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

	var res models.DeleteDMTBeneficiaryVerificationResponseModel
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, err
	}
	if res.Error != 0 {
		return nil, fmt.Errorf("error from rechargekit: %s", res.Description)
	}
	return &res, nil
}
