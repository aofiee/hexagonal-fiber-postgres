package schema

import (
	"hexagonal/architecture/resolver"

	"github.com/graphql-go/graphql"
)

var (
	Customer = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Customer",
			Fields: graphql.Fields{
				"customer_id": &graphql.Field{
					Type: graphql.Int,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"date_of_birth": &graphql.Field{
					Type: graphql.String,
				},
				"city": &graphql.Field{
					Type: graphql.String,
				},
				"zip_code": &graphql.Field{
					Type: graphql.String,
				},
				"status": &graphql.Field{
					Type: graphql.Int,
				},
			},
		},
	)
)

type (
	customerSchema struct {
		customerResolver resolver.CustomerResolver
	}
)

func NewCustomerSchema(customerResolver resolver.CustomerResolver) customerSchema {
	return customerSchema{
		customerResolver: customerResolver,
	}
}

func (c customerSchema) Query() *graphql.Object {
	objectConfig := graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"GetCustomers": &graphql.Field{
				Type:        graphql.NewList(Customer),
				Description: "Get all Customer",
				Resolve:     c.customerResolver.GetCustomers,
			},
			"GetCustomer": &graphql.Field{
				Type:        Customer,
				Description: "Get Customer By ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: c.customerResolver.GetCustomer,
			},
		},
	}

	return graphql.NewObject(objectConfig)
}

func (c customerSchema) Mutation() *graphql.Object {
	objectConfig := graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"CreateCustomer": &graphql.Field{
				Type:        graphql.String,
				Description: "Store a new customer",
				Args: graphql.FieldConfigArgument{
					"CustomerID": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"Name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"DateOfBirth": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"City": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"ZipCode": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"Status": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: c.customerResolver.CreateCustomer,
			},
		},
	}
	return graphql.NewObject(objectConfig)
}
