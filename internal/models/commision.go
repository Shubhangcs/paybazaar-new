package models

import "time"

type CreateCommisionModel struct {
	UserID                     string  `json:"user_id" validate:"required"`
	TotalCommision             float64 `json:"total_commision" validate:"required,gt=0"`
	AdminCommision             float64 `json:"admin_commision" validate:"gte=0"`
	MasterDistributorCommision float64 `json:"master_distributor_commision" validate:"gte=0"`
	DistributorCommision       float64 `json:"distributor_commision" validate:"gte=0"`
	RetailerCommision          float64 `json:"retailer_commision" validate:"gte=0"`
}

type UpdateCommisionModel struct {
	TotalCommision             *float64 `json:"total_commision,omitempty"`
	AdminCommision             *float64 `json:"admin_commision,omitempty"`
	MasterDistributorCommision *float64 `json:"master_distributor_commision,omitempty"`
	DistributorCommision       *float64 `json:"distributor_commision,omitempty"`
	RetailerCommision          *float64 `json:"retailer_commision,omitempty"`
}

type GetCommisionModel struct {
	CommisionID                int64     `json:"commision_id"`
	UserID                     string    `json:"user_id"`
	TotalCommision             float64   `json:"total_commision"`
	AdminCommision             float64   `json:"admin_commision"`
	MasterDistributorCommision float64   `json:"master_distributor_commision"`
	DistributorCommision       float64   `json:"distributor_commision"`
	RetailerCommision          float64   `json:"retailer_commision"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}
