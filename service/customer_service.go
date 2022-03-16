package service

import "hexagonal/architecture/repository"

type (
	customerService struct {
		customerRepository repository.CustomerRepository
	}
)

func NewCustomerService(cus repository.CustomerRepository) CustomerService {
	return customerService{customerRepository: cus}
}

func (c customerService) GetCustomer(id int) (CustomerRes, error) {
	customer, err := c.customerRepository.GetByID(id)
	if err != nil {
		return CustomerRes{}, err
	}
	return CustomerRes{
		CustomerID:  customer.CustomerID,
		Name:        customer.Name,
		DateOfBirth: customer.DateOfBirth,
		City:        customer.City,
		ZipCode:     customer.ZipCode,
		Status:      customer.Status,
	}, nil
}

func (c customerService) GetCustomers() ([]CustomerRes, error) {
	customers, err := c.customerRepository.GetAll()
	if err != nil {
		return nil, err
	}
	var customerRes []CustomerRes
	for _, v := range customers {
		customerRes = append(customerRes, CustomerRes{
			CustomerID:  v.CustomerID,
			Name:        v.Name,
			DateOfBirth: v.DateOfBirth,
			City:        v.City,
			ZipCode:     v.ZipCode,
			Status:      v.Status,
		})
	}
	return customerRes, nil
}
