package repository

import (
	"gorm.io/gorm"
)

// Adapter
type (
	customerRepository struct {
		db *gorm.DB
	}
)

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return customerRepository{db: db}
}

func (c customerRepository) GetAll() ([]Customer, error) {
	var customers []Customer
	err := c.db.Find(&customers).Error
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (c customerRepository) GetByID(id int) (*Customer, error) {
	var customer Customer
	err := c.db.Where("id = ?", id).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}
