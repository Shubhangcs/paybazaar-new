package models

import "time"

type CreateCommisionRequestModel struct {
	UserID                     string  `json:"user_id" validate:"required"`
	Service                    string  `json:"service" validate:"required"`
	TotalCommision             float64 `json:"total_commision" validate:"required"`
	AdminCommision             float64 `json:"admin_commision" validate:"required"`
	MasterDistributorCommision float64 `json:"master_distributor_commision" validate:"required"`
	DistributorCommision       float64 `json:"distributor_commision" validate:"required"`
	RetailerCommision          float64 `json:"retailer_commision" validate:"required"`
}

type UpdateCommisionRequestModel struct {
	CommisionID                int64    `json:"commision_id" validate:"required"`
	TotalCommision             *float64 `json:"total_commision" validate:"omitempty"`
	AdminCommision             *float64 `json:"admin_commision" validate:"omitempty"`
	MasterDistributorCommision *float64 `json:"master_distributor_commision" validate:"omitempty"`
	DistributorCommision       *float64 `json:"distributor_commision" validate:"omitempty"`
	RetailerCommision          *float64 `json:"retailer_commision" validate:"omitempty"`
}

type GetCommisionResponseModel struct {
	CommisionID                int64     `json:"commision_id"`
	UserID                     string    `json:"user_id"`
	Service                    string    `json:"service"`
	TotalCommision             float64   `json:"total_commision"`
	AdminCommision             float64   `json:"admin_commision"`
	MasterDistributorCommision float64   `json:"master_distributor_commision"`
	DistributorCommision       float64   `json:"distributor_commision"`
	RetailerCommision          float64   `json:"retailer_commision"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

type GetTDSCommisionResponseModel struct {
	TDSCommisionID int64     `json:"tds_commision_id"`
	TransactionID  string    `json:"transaction_id"`
	UserID         string    `json:"user_id"`
	UserName       string    `json:"user_name"`
	Commision      float64   `json:"commision"`
	TDS            float64   `json:"tds"`
	PaidCommision  float64   `json:"paid_commision"`
	PANNumber      string    `json:"pan_number"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

