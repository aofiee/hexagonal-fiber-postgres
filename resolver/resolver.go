package resolver

import (
	"fmt"
	"hexagonal/architecture/service"
	"log"

	"github.com/graphql-go/graphql"
)

type (
	customerResolver struct {
		customerResolver service.CustomerService
	}
	CustomerResolver interface {
		GetCustomer(params graphql.ResolveParams) (interface{}, error)
		GetCustomers(params graphql.ResolveParams) (interface{}, error)
	}
)

func NewCustomerResolver(cus service.CustomerService) CustomerResolver {
	return customerResolver{
		customerResolver: cus,
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
	customer, err := c.customerResolver.GetCustomer(id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (c customerResolver) GetCustomers(params graphql.ResolveParams) (interface{}, error) {
	customers, err := c.customerResolver.GetCustomers()
	if err != nil {
		return nil, err
	}
	log.Println(`customers`, customers)
	return customers, nil
}
