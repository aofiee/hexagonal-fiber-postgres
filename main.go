package main

import (
	"hexagonal/architecture/handler"
	"hexagonal/architecture/repository"
	"hexagonal/architecture/resolver"
	myschema "hexagonal/architecture/schema"
	"hexagonal/architecture/service"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/graphql-go/graphql"
	gqlHandler "github.com/graphql-go/handler"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dns := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	dial := postgres.Open(dns)
	// dns := "root:password@tcp(localhost:3306)/dd_hakka?charset=utf8mb4&parseTime=True&loc=Local"
	// dial := mysql.Open(dns)
	db, err := createDatabaseConnection(dial, dns)
	if err != nil {
		log.Println(err)
	}

	customerRepo := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerService)

	cResolver := resolver.NewCustomerResolver(customerService)
	cSchema := myschema.NewCustomerSchema(cResolver)
	graphqlSchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: cSchema.Query(),
	})
	if err != nil {
		log.Println(err)
	}
	gh := gqlHandler.New(&gqlHandler.Config{
		Schema:   &graphqlSchema,
		GraphiQL: true,
		Pretty:   true,
	})

	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024,
	})
	app.Static("/static", "../public")
	app.Use(requestid.New())
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format:     "[${time}] ${method} ${path}",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Asia/Bangkok",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept,Authorization",
	}))

	app.Get("/customer/:id", customerHandler.GetCustomer)
	app.Get("/customers", customerHandler.GetCustomers)

	app.Get("/graph", func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			gh.ServeHTTP(writer, request)
		})(c.Context())
		return nil
	})
	app.Post("/graph", func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			gh.ServeHTTP(writer, request)
		})(c.Context())
		return nil
	})

	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}

func createDatabaseConnection(dial gorm.Dialector, dns string) (db *gorm.DB, err error) {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, err
	}
	time.Local = loc
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Error,
			Colorful:      true,
		},
	)

	mydb, err := gorm.Open(dial, &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := mydb.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetMaxOpenConns(100)

	mydb.AutoMigrate(&repository.Customer{})
	return mydb, nil
}

func PlaygroundHandler(c *fiber.Ctx) error {
	h := playground.Handler("GraphQL", "/query")
	fasthttpadaptor.NewFastHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		h.ServeHTTP(writer, request)
	})(c.Context())
	return nil
}
