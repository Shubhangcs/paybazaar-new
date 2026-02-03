package repositories

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/levion-studio/paybazaar/internal/database"
	"github.com/levion-studio/paybazaar/internal/models"
)

type Beneficiary interface {
	GetBeneficiaries(echo.Context) (*[]models.BeneficiaryModel, error)
	AddNewBeneficiary(echo.Context) error
	VerifyBeneficiary(echo.Context) error
	DeleteBeneficiary(echo.Context) error
}

type beneficiaryRepo struct {
	query *database.Database
}

func NewBeneficiaryRepo(query *database.Database) *beneficiaryRepo {
	return &beneficiaryRepo{query: query}
}

func (r *beneficiaryRepo) GetBeneficiaries(e echo.Context) (*[]models.BeneficiaryModel, error) {
	var phone = e.Param("phone")
	if phone == "" {
		return nil, fmt.Errorf("phone number not found")
	}
	res, err := r.query.GetBeneficiaries(phone)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch benificiries")
	}
	return res, nil
}

func (r *beneficiaryRepo) AddNewBeneficiary(e echo.Context) error {
	req := &models.BeneficiaryModel{}
	if err := e.Bind(req); err != nil {
		return err
	}
	err := r.query.AddNewBeneficiary(req)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to add benificary")
	}
	return nil
}

func (r *beneficiaryRepo) VerifyBeneficiary(e echo.Context) error {
	beneficiaryId := e.Param("ben_id")
	if beneficiaryId == "" {
		return fmt.Errorf("beneficiary id not found")
	}
	if err := r.query.VerifyBenificary(beneficiaryId); err != nil {
		return fmt.Errorf("failed to verify beneficiary")
	}
	return nil
}

func (r *beneficiaryRepo) DeleteBeneficiary(e echo.Context) error {
	beneficiaryId := e.Param("ben_id")
	if beneficiaryId == "" {
		return fmt.Errorf("beneficiary id not found")
	}
	if err := r.query.DeleteBeneficiary(beneficiaryId); err != nil {
		return fmt.Errorf("failed to delete beneficiary")
	}
	return nil
}
