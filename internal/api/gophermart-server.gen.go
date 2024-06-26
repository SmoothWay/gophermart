// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get user's balance
	// (GET /api/user/balance)
	GetBalance(w http.ResponseWriter, r *http.Request)
	// Withdrawal request
	// (POST /api/user/balance/withdraw)
	WithdrawalRequest(w http.ResponseWriter, r *http.Request)
	// Logs user into the system
	// (POST /api/user/login)
	LoginUser(w http.ResponseWriter, r *http.Request)
	// Get list of orders
	// (GET /api/user/orders)
	GetOrders(w http.ResponseWriter, r *http.Request)
	// Upload order number
	// (POST /api/user/orders)
	UploadOrder(w http.ResponseWriter, r *http.Request)
	// Register user
	// (POST /api/user/register)
	RegisterUser(w http.ResponseWriter, r *http.Request)
	// Get list of successful withdrawals
	// (GET /api/user/withdrawals)
	GetWithdrawals(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Get user's balance
// (GET /api/user/balance)
func (_ Unimplemented) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Withdrawal request
// (POST /api/user/balance/withdraw)
func (_ Unimplemented) WithdrawalRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Logs user into the system
// (POST /api/user/login)
func (_ Unimplemented) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get list of orders
// (GET /api/user/orders)
func (_ Unimplemented) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Upload order number
// (POST /api/user/orders)
func (_ Unimplemented) UploadOrder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Register user
// (POST /api/user/register)
func (_ Unimplemented) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get list of successful withdrawals
// (GET /api/user/withdrawals)
func (_ Unimplemented) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetBalance operation middleware
func (siw *ServerInterfaceWrapper) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetBalance(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// WithdrawalRequest operation middleware
func (siw *ServerInterfaceWrapper) WithdrawalRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.WithdrawalRequest(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// LoginUser operation middleware
func (siw *ServerInterfaceWrapper) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.LoginUser(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetOrders operation middleware
func (siw *ServerInterfaceWrapper) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetOrders(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// UploadOrder operation middleware
func (siw *ServerInterfaceWrapper) UploadOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UploadOrder(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// RegisterUser operation middleware
func (siw *ServerInterfaceWrapper) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RegisterUser(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetWithdrawals operation middleware
func (siw *ServerInterfaceWrapper) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetWithdrawals(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/user/balance", wrapper.GetBalance)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/api/user/balance/withdraw", wrapper.WithdrawalRequest)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/api/user/login", wrapper.LoginUser)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/user/orders", wrapper.GetOrders)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/api/user/orders", wrapper.UploadOrder)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/api/user/register", wrapper.RegisterUser)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/user/withdrawals", wrapper.GetWithdrawals)
	})

	return r
}

type GetBalanceRequestObject struct {
}

type GetBalanceResponseObject interface {
	VisitGetBalanceResponse(w http.ResponseWriter) error
}

type GetBalance200JSONResponse Balance

func (response GetBalance200JSONResponse) VisitGetBalanceResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetBalance401JSONResponse Error

func (response GetBalance401JSONResponse) VisitGetBalanceResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type GetBalance500JSONResponse Error

func (response GetBalance500JSONResponse) VisitGetBalanceResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type WithdrawalRequestRequestObject struct {
	Body *WithdrawalRequestJSONRequestBody
}

type WithdrawalRequestResponseObject interface {
	VisitWithdrawalRequestResponse(w http.ResponseWriter) error
}

type WithdrawalRequest200Response struct {
}

func (response WithdrawalRequest200Response) VisitWithdrawalRequestResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type WithdrawalRequest401JSONResponse Error

func (response WithdrawalRequest401JSONResponse) VisitWithdrawalRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type WithdrawalRequest402JSONResponse Error

func (response WithdrawalRequest402JSONResponse) VisitWithdrawalRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(402)

	return json.NewEncoder(w).Encode(response)
}

type WithdrawalRequest422JSONResponse Error

func (response WithdrawalRequest422JSONResponse) VisitWithdrawalRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)

	return json.NewEncoder(w).Encode(response)
}

type WithdrawalRequest500JSONResponse Error

func (response WithdrawalRequest500JSONResponse) VisitWithdrawalRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type LoginUserRequestObject struct {
	Body *LoginUserJSONRequestBody
}

type LoginUserResponseObject interface {
	VisitLoginUserResponse(w http.ResponseWriter) error
}

type LoginUser200ResponseHeaders struct {
	Authorization string
}

type LoginUser200Response struct {
	Headers LoginUser200ResponseHeaders
}

func (response LoginUser200Response) VisitLoginUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Authorization", fmt.Sprint(response.Headers.Authorization))
	w.WriteHeader(200)
	return nil
}

type LoginUser400JSONResponse Error

func (response LoginUser400JSONResponse) VisitLoginUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type LoginUser401JSONResponse Error

func (response LoginUser401JSONResponse) VisitLoginUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type LoginUser500JSONResponse Error

func (response LoginUser500JSONResponse) VisitLoginUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetOrdersRequestObject struct {
}

type GetOrdersResponseObject interface {
	VisitGetOrdersResponse(w http.ResponseWriter) error
}

type GetOrders200JSONResponse []Order

func (response GetOrders200JSONResponse) VisitGetOrdersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetOrders204Response struct {
}

func (response GetOrders204Response) VisitGetOrdersResponse(w http.ResponseWriter) error {
	w.WriteHeader(204)
	return nil
}

type GetOrders401JSONResponse Error

func (response GetOrders401JSONResponse) VisitGetOrdersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type GetOrders500JSONResponse Error

func (response GetOrders500JSONResponse) VisitGetOrdersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type UploadOrderRequestObject struct {
	Body *UploadOrderTextRequestBody
}

type UploadOrderResponseObject interface {
	VisitUploadOrderResponse(w http.ResponseWriter) error
}

type UploadOrder200Response struct {
}

func (response UploadOrder200Response) VisitUploadOrderResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type UploadOrder202Response struct {
}

func (response UploadOrder202Response) VisitUploadOrderResponse(w http.ResponseWriter) error {
	w.WriteHeader(202)
	return nil
}

type UploadOrder400JSONResponse Error

func (response UploadOrder400JSONResponse) VisitUploadOrderResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type UploadOrder401JSONResponse Error

func (response UploadOrder401JSONResponse) VisitUploadOrderResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type UploadOrder409JSONResponse Error

func (response UploadOrder409JSONResponse) VisitUploadOrderResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

type UploadOrder422JSONResponse Error

func (response UploadOrder422JSONResponse) VisitUploadOrderResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)

	return json.NewEncoder(w).Encode(response)
}

type UploadOrder500JSONResponse Error

func (response UploadOrder500JSONResponse) VisitUploadOrderResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type RegisterUserRequestObject struct {
	Body *RegisterUserJSONRequestBody
}

type RegisterUserResponseObject interface {
	VisitRegisterUserResponse(w http.ResponseWriter) error
}

type RegisterUser200ResponseHeaders struct {
	Authorization string
}

type RegisterUser200Response struct {
	Headers RegisterUser200ResponseHeaders
}

func (response RegisterUser200Response) VisitRegisterUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Authorization", fmt.Sprint(response.Headers.Authorization))
	w.WriteHeader(200)
	return nil
}

type RegisterUser400JSONResponse Error

func (response RegisterUser400JSONResponse) VisitRegisterUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type RegisterUser409JSONResponse Error

func (response RegisterUser409JSONResponse) VisitRegisterUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

type RegisterUser500JSONResponse Error

func (response RegisterUser500JSONResponse) VisitRegisterUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetWithdrawalsRequestObject struct {
}

type GetWithdrawalsResponseObject interface {
	VisitGetWithdrawalsResponse(w http.ResponseWriter) error
}

type GetWithdrawals200JSONResponse []Withdrawal

func (response GetWithdrawals200JSONResponse) VisitGetWithdrawalsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetWithdrawals204Response struct {
}

func (response GetWithdrawals204Response) VisitGetWithdrawalsResponse(w http.ResponseWriter) error {
	w.WriteHeader(204)
	return nil
}

type GetWithdrawals401JSONResponse Error

func (response GetWithdrawals401JSONResponse) VisitGetWithdrawalsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type GetWithdrawals500JSONResponse Error

func (response GetWithdrawals500JSONResponse) VisitGetWithdrawalsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Get user's balance
	// (GET /api/user/balance)
	GetBalance(ctx context.Context, request GetBalanceRequestObject) (GetBalanceResponseObject, error)
	// Withdrawal request
	// (POST /api/user/balance/withdraw)
	WithdrawalRequest(ctx context.Context, request WithdrawalRequestRequestObject) (WithdrawalRequestResponseObject, error)
	// Logs user into the system
	// (POST /api/user/login)
	LoginUser(ctx context.Context, request LoginUserRequestObject) (LoginUserResponseObject, error)
	// Get list of orders
	// (GET /api/user/orders)
	GetOrders(ctx context.Context, request GetOrdersRequestObject) (GetOrdersResponseObject, error)
	// Upload order number
	// (POST /api/user/orders)
	UploadOrder(ctx context.Context, request UploadOrderRequestObject) (UploadOrderResponseObject, error)
	// Register user
	// (POST /api/user/register)
	RegisterUser(ctx context.Context, request RegisterUserRequestObject) (RegisterUserResponseObject, error)
	// Get list of successful withdrawals
	// (GET /api/user/withdrawals)
	GetWithdrawals(ctx context.Context, request GetWithdrawalsRequestObject) (GetWithdrawalsResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// GetBalance operation middleware
func (sh *strictHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	var request GetBalanceRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetBalance(ctx, request.(GetBalanceRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetBalance")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetBalanceResponseObject); ok {
		if err := validResponse.VisitGetBalanceResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// WithdrawalRequest operation middleware
func (sh *strictHandler) WithdrawalRequest(w http.ResponseWriter, r *http.Request) {
	var request WithdrawalRequestRequestObject

	var body WithdrawalRequestJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.WithdrawalRequest(ctx, request.(WithdrawalRequestRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "WithdrawalRequest")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(WithdrawalRequestResponseObject); ok {
		if err := validResponse.VisitWithdrawalRequestResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// LoginUser operation middleware
func (sh *strictHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var request LoginUserRequestObject

	var body LoginUserJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.LoginUser(ctx, request.(LoginUserRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "LoginUser")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(LoginUserResponseObject); ok {
		if err := validResponse.VisitLoginUserResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetOrders operation middleware
func (sh *strictHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	var request GetOrdersRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetOrders(ctx, request.(GetOrdersRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetOrders")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetOrdersResponseObject); ok {
		if err := validResponse.VisitGetOrdersResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// UploadOrder operation middleware
func (sh *strictHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	var request UploadOrderRequestObject

	data, err := io.ReadAll(r.Body)
	if err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't read body: %w", err))
		return
	}
	body := UploadOrderTextRequestBody(data)
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.UploadOrder(ctx, request.(UploadOrderRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "UploadOrder")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(UploadOrderResponseObject); ok {
		if err := validResponse.VisitUploadOrderResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// RegisterUser operation middleware
func (sh *strictHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request RegisterUserRequestObject

	var body RegisterUserJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.RegisterUser(ctx, request.(RegisterUserRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "RegisterUser")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(RegisterUserResponseObject); ok {
		if err := validResponse.VisitRegisterUserResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetWithdrawals operation middleware
func (sh *strictHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	var request GetWithdrawalsRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetWithdrawals(ctx, request.(GetWithdrawalsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetWithdrawals")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetWithdrawalsResponseObject); ok {
		if err := validResponse.VisitGetWithdrawalsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xX32vjRhD+V5ZtoS9u5POlD+e3Bko5OBq4Eu6hhDKWxvYe0q5uZhSfG/S/l92VZDlS",
	"4lzPMaH4bdH++ObH982M7nXqitJZtMJ6fq85XWMBYXkFOdgU/bIkVyKJwbCRVkRoxS+XjgoQPdeZqxY5",
	"6omWbYl6rm1VLJB0PdEbI+uMYGOfdb6eaMIvlSHM9PyvDqr/zG13yS0+Yyoe5DciR0NLC2SGVXAhQ07J",
	"lGKc1fN4XrXb3YMsZOxqYEV7bgz5mjIcQYY0pQryZ8aoWc3vHxoy0SwgFY9uVWXuIMPsb3iQCxD8WUxx",
	"2LEGt0PZf3PM3Rse8zZ3K2NHbSyBeeMo2zOw+3jIvvhu75Uxkz41zIjR3jfMtdkZGkYuReZvit5Ec1X8",
	"FxpHM+L1oQv+tLFLF+w0kvu93125RiqAPPfvkDjydnrx5mLqDXElWiiNnuu3F9OLtyFGsg4+J1CapGKk",
	"ZLFT8AqDlz424DXwPvMgKK3Ivb1cOssxbrPpNCjdWWmUDmWZmzTcTT6zs7ti4Vc/Ei71XP+Q7KpJ0pSS",
	"pIUIfu7L8M8q9UlYVrnqLPPeXU7fHA0+loYR8BsLlawdmX8w86C/HNHnR0HfW0GykCtGukNS2BwM3CiA",
	"tjEtyufvJ1aLLj0CK/ZU8hv61l8Y5Dlpa2TQgeORjO+08hG/VMiiI1GR5cpl26P539NkXUc1DNn1irlw",
	"OZ2dggtcLZcmNWily7QHn50E/A5yk6lQmtSuFb0aFewopKjj6lMq6JrQOPU/+O3Qvl6G8uHp7yL7RK8R",
	"MqRw7deGj3Frz4qHPbMOhD1B2q4g63JxMmUyksocsrJOFH41EfvV0PSDW3Go1spYcUrWqHjLgsUBtgbd",
	"8VOt+Tqe+M7ObAQLPhSBOMTW3WQCRLD9loY9m14OSf6HU42X55aOonLDotyyjcmAHpNHCtdNGMmvmxHy",
	"8dIl+FWSMgdzqFz0Z1OhCp9VsQK+WjpSzT9ZJD3khJBtozA5UmH22G3DKiUEaXvs/7VkDYaJdy8PGiM8",
	"TMd5mNh2Gtq37+n6TLgyLM3P7qguPzYnzjPFUQX67kQzxVAqr4atLbNCiT3A0003Jj85THzqHTvFRNH/",
	"ATz+WNF3+jxb9GYL3gVys5fxhxSqu0/32kLhsxO26tv63wAAAP//iPwOCBMWAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
