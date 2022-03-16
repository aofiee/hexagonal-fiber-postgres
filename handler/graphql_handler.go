package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"

	"github.com/gofiber/fiber/v2"

	"context"

	"github.com/graphql-go/graphql/gqlerrors"
)

const (
	/*ContentTypeJSON as request type*/
	ContentTypeJSON = "application/json"
	/*ContentTypeGraphQL as request type*/
	ContentTypeGraphQL = "application/graphql"
	/*ContentTypeFormURLEncoded as request type*/
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

/*ResultCallbackFn is function type which required on callback*/
type ResultCallbackFn func(
	ctx context.Context,
	params *graphql.Params,
	result *graphql.Result,
	responseBody []byte,
)

/*Handler is used to handle net/http handler*/
type Handler struct {
	Schema           *graphql.Schema
	pretty           bool
	graphiql         bool
	playground       bool
	rootObjectFn     RootObjectFn
	resultCallbackFn ResultCallbackFn
	formatErrorFn    func(err error) gqlerrors.FormattedError
}

/*RequestOptions are the available options for graphql request*/
type RequestOptions struct {
	Query         string                 `json:"query" url:"query" schema:"query"`
	Variables     map[string]interface{} `json:"variables" url:"variables" schema:"variables"`
	OperationName string                 `json:"operationName" url:"operationName" schema:"operationName"`
}

// a workaround for getting`variables` as a JSON string
type requestOptionsCompatibility struct {
	Query         string `json:"query" url:"query" schema:"query"`
	Variables     string `json:"variables" url:"variables" schema:"variables"`
	OperationName string `json:"operationName" url:"operationName" schema:"operationName"`
}

func getFromForm(c *fiber.Ctx) *RequestOptions {
	query := c.Query("query")
	if query != "" {
		// TODO: For fiber => get variables map. Remove static 100 with query len
		variables := make(map[string]interface{}, 100)
		// variablesStr := c.Query("variables")
		// json.Unmarshal([]byte(variablesStr), &variables)

		return &RequestOptions{
			Query:         query,
			Variables:     variables,
			OperationName: c.Query("operationName"),
		}
	}

	return nil
}

/*NewRequestOptions Parses a http.Request into GraphQL request options struct*/
func NewRequestOptions(c *fiber.Ctx) *RequestOptions {
	if reqOpt := getFromForm(c); reqOpt != nil {
		return reqOpt
	}

	if c.Method() != http.MethodPost {
		return &RequestOptions{}
	}

	if c.Body() == nil {
		return &RequestOptions{}
	}

	// TODO: improve Content-Type handling
	contentTypeStr := string(c.Request().Header.ContentType())
	contentTypeTokens := strings.Split(contentTypeStr, ";")
	contentType := contentTypeTokens[0]

	switch contentType {
	case ContentTypeGraphQL:
		body := c.Body()
		if body == nil {
			return &RequestOptions{}
		}
		return &RequestOptions{
			Query: string(body),
		}
		// TODO: Do for fiber handler
	// case ContentTypeFormURLEncoded:
	// 	if err := r.ParseForm(); err != nil {
	// 		return &RequestOptions{}
	// 	}

	// 	if reqOpt := getFromForm(r.PostForm); reqOpt != nil {
	// 		return reqOpt
	// 	}

	// 	return &RequestOptions{}

	case ContentTypeJSON:
		fallthrough
	default:
		opts := new(RequestOptions)
		err := c.BodyParser(opts)
		if err != nil {
			// Probably `variables` was sent as a string instead of an object.
			// So, we try to be polite and try to parse that as a JSON string
			var optsCompatible requestOptionsCompatibility
			json.Unmarshal(c.Body(), &optsCompatible)
			json.Unmarshal([]byte(optsCompatible.Variables), &opts.Variables)
		}
		return (opts)
	}
}

// ContextHandler provides an entrypoint into executing graphQL queries with a
// user-provided context.
func (h *Handler) ContextHandler(
	ctx context.Context,
	c *fiber.Ctx,
) {
	// get query
	opts := NewRequestOptions(c)

	// execute graphql query
	params := graphql.Params{
		Schema:         *h.Schema,
		RequestString:  opts.Query,
		VariableValues: opts.Variables,
		OperationName:  opts.OperationName,
		Context:        ctx,
	}
	if h.rootObjectFn != nil {
		params.RootObject = h.rootObjectFn(ctx, c)
	}
	result := graphql.Do(params)

	if formatErrorFn := h.formatErrorFn; formatErrorFn != nil && len(result.Errors) > 0 {
		formatted := make([]gqlerrors.FormattedError, len(result.Errors))
		for i, formattedError := range result.Errors {
			formatted[i] = formatErrorFn(formattedError.OriginalError())
		}
		result.Errors = formatted
	}

	// if h.graphiql {
	// 	acceptHeader := c.Response().Header.//r.Header.Get("Accept")
	// 	_, raw := r.URL.Query()["raw"]
	// 	if !raw && !strings.Contains(acceptHeader, "application/json") && strings.Contains(acceptHeader, "text/html") {
	// 		renderGraphiQL(w, params)
	// 		return
	// 	}
	// }

	// if h.playground {
	// 	acceptHeader := r.Header.Get("Accept")
	// 	_, raw := r.URL.Query()["raw"]
	// 	if !raw && !strings.Contains(acceptHeader, "application/json") && strings.Contains(acceptHeader, "text/html") {
	// 		renderPlayground(w, r)
	// 		return
	// 	}
	// }

	// use proper JSON Header
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	c.Status(fiber.StatusOK)

	var buff []byte
	if h.pretty {
		buff, _ = json.MarshalIndent(result, "", "\t")
	} else {
		buff, _ = json.Marshal(result)
	}
	c.Send(buff)

	if h.resultCallbackFn != nil {
		h.resultCallbackFn(ctx, &params, result, buff)
	}
}

// ServeHTTP provides an entrypoint into executing graphQL queries.
func (h *Handler) ServeHTTP(c *fiber.Ctx) error {
	h.ContextHandler(c.Context(), c)
	return nil
}

// RootObjectFn allows a user to generate a RootObject per request
type RootObjectFn func(
	ctx context.Context,
	c *fiber.Ctx,
) map[string]interface{}

/*Config are the required configuration for net/http handler*/
type Config struct {
	Schema           *graphql.Schema
	Pretty           bool
	GraphiQL         bool
	Playground       bool
	RootObjectFn     RootObjectFn
	ResultCallbackFn ResultCallbackFn
	FormatErrorFn    func(err error) gqlerrors.FormattedError
}

/*NewConfig returns config object with default values*/
func NewConfig() *Config {
	return &Config{
		Schema:     nil,
		Pretty:     true,
		GraphiQL:   true,
		Playground: false,
	}
}

/*New creates a handler based on config param*/
func New(p *Config) *Handler {
	if p == nil {
		p = NewConfig()
	}

	if p.Schema == nil {
		panic("undefined GraphQL schema")
	}

	return &Handler{
		Schema:           p.Schema,
		pretty:           p.Pretty,
		graphiql:         p.GraphiQL,
		playground:       p.Playground,
		rootObjectFn:     p.RootObjectFn,
		resultCallbackFn: p.ResultCallbackFn,
		formatErrorFn:    p.FormatErrorFn,
	}
}
