package resolver

import (
	"fmt"
	"hexagonal/architecture/service"
	"log"

	"github.com/graphql-go/graphql"
)

type (
	customerResolver struct {
		customerService service.CustomerService
	}
	CustomerResolver interface {
		GetCustomer(params graphql.ResolveParams) (interface{}, error)
		GetCustomers(params graphql.ResolveParams) (interface{}, error)
		CreateCustomer(params graphql.ResolveParams) (interface{}, error)
	}
)

func NewCustomerResolver(cus service.CustomerService) CustomerResolver {
	return customerResolver{
		customerService: cus,
	}
}

func (c customerResolver) GetCustomer(params graphql.ResolveParams) (interface{}, error) {
	var (
		id int
		ok bool
	)
	if id, ok = params.Args["id"].(int); !ok || id == 0 {
		return nil, fmt.Errorf("id is not integer or zero")
	}
	customer, err := c.customerService.GetCustomer(id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (c customerResolver) GetCustomers(params graphql.ResolveParams) (interface{}, error) {
	customers, err := c.customerService.GetCustomers()
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (c customerResolver) CreateCustomer(params graphql.ResolveParams) (interface{}, error) {
	var (
		cusID       int
		cusName     string
		dateOfBirth string
		city        string
		zipCode     string
		status      int
		ok          bool
	)
	if cusID, ok = params.Args["CustomerID"].(int); !ok || cusID == 0 {
		return nil, fmt.Errorf("id is not int or 0")
	}
	if cusName, ok = params.Args["Name"].(string); !ok || cusName == "" {
		return nil, fmt.Errorf("id is not string or nil")
	}
	if dateOfBirth, ok = params.Args["DateOfBirth"].(string); !ok || dateOfBirth == "" {
		return nil, fmt.Errorf("id is not string or nil")
	}
	if city, ok = params.Args["City"].(string); !ok || city == "" {
		return nil, fmt.Errorf("id is not string or nil")
	}
	if zipCode, ok = params.Args["ZipCode"].(string); !ok || zipCode == "" {
		return nil, fmt.Errorf("id is not string or nil")
	}
	if status, ok = params.Args["Status"].(int); !ok || status == 0 {
		return nil, fmt.Errorf("id is not int or 0")
	}
	log.Println(`input`, cusID, cusName, dateOfBirth, city, zipCode, status)
	newCus := service.CustomerRes{
		CustomerID:  cusID,
		Name:        cusName,
		DateOfBirth: dateOfBirth,
		City:        city,
		ZipCode:     zipCode,
		Status:      status,
	}
	res, err := c.customerService.CreateCustomer(&newCus)
	if err != nil {
		return nil, err
	}
	return res, nil
}
