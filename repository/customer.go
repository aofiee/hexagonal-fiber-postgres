package repository

import (
	"time"

	"gorm.io/gorm"
)

// Ports
type (
	BarrothModel struct {
		ID        uint           `gorm:"primaryKey" json:"id"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
		DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`
	}
	Customer struct {
		BarrothModel
		CustomerID  int    `gorm:"type:INT4" json:"customer_id"`
		Name        string `gorm:"type:VARCHAR(100)" json:"name"`
		DateOfBirth string `gorm:"type:DATE" json:"date_of_birth"`
		City        string `gorm:"type:VARCHAR(100)" json:"city"`
		ZipCode     string `gorm:"type:VARCHAR(10)" json:"zip_code"`
		Status      int    `gorm:"type:INT2" json:"status"`
	}
	CustomerRepository interface {
		GetAll() ([]Customer, error)
		GetByID(id int) (*Customer, error)
		CreateCustomer(customer *Customer) (Customer, error)
	}
)
