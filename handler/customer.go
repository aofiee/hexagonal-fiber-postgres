package handler

import (
	"hexagonal/architecture/service"

	"github.com/gofiber/fiber/v2"
)

type (
	customerHandler struct {
		customerService service.CustomerService
	}
)

func NewCustomerHandler(cus service.CustomerService) customerHandler {
	return customerHandler{customerService: cus}
}

func (c customerHandler) GetCustomer(f *fiber.Ctx) error {
	id, err := f.ParamsInt("id", 1)
	if err != nil {
		return err
	}
	customer, err := c.customerService.GetCustomer(id)
	if err != nil {
		return err
	}
	return f.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":  customer,
		"error": nil,
	})
}

func (c customerHandler) GetCustomers(f *fiber.Ctx) error {
	customers, err := c.customerService.GetCustomers()
	if err != nil {
		return err
	}
	return f.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":  customers,
		"error": nil,
	})
}
