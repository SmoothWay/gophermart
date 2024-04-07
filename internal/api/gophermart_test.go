package api_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/SmoothWay/gophermart/internal/api"
	"github.com/SmoothWay/gophermart/internal/logger"
	"github.com/SmoothWay/gophermart/internal/mocks"
	"github.com/SmoothWay/gophermart/internal/model"
	"github.com/SmoothWay/gophermart/internal/repository/postgres"
	"github.com/SmoothWay/gophermart/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	baseURL        = "http://localhost:8080"
	withdrawalReq  = baseURL + "/api/user/balance/withdraw"
	uploadOrder    = baseURL + "/api/user/orders"
	getOrders      = baseURL + "/api/user/orders"
	getWithdrawals = baseURL + "/api/user/withdrawals"
	getBalance     = baseURL + "/api/user/balance"
	loginUser      = baseURL + "/api/user/login"
	registerUser   = baseURL + "/api/user/register"
)

type apiTestSuite struct {
	suite.Suite
	r       *chi.Mux
	service *mocks.Service
}

func (suite *apiTestSuite) SetupTest() {
	logger.InitSlog("ERROR")
	srv := &mocks.Service{}
	swagger, err := api.GetSwagger()
	suite.NoError(err)
	swagger.Servers = nil
	r := chi.NewRouter()
	r.Use(middleware.OapiRequestValidator(swagger))
	g := api.NewGophermart(srv, []byte("secret"))
	sh := api.NewStrictHandler(g, nil)
	api.HandlerFromMux(sh, r)
	suite.r = r
	suite.service = srv
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(apiTestSuite))
}

func (suite *apiTestSuite) TestRegisterUser() {
	u := api.RegisterUserJSONRequestBody{
		Login:    "login",
		Password: "password",
	}
	b, err := json.Marshal(u)
	suite.NoError(err)

	suite.Run("Register user", func() {
		req, err := http.NewRequest(http.MethodPost, registerUser, bytes.NewReader(b))
		suite.NoError(err)
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.service.On("RegisterUser", mock.Anything, u.Login, u.Password).Once().Return(nil)
		suite.service.On("Authenticate", mock.Anything, u.Login, u.Password).Once().Return("token", nil)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()

		suite.Equal("Bearer token", res.Header.Get("Authorization"))
	})

	suite.Run("User already exists", func() {
		req, err := http.NewRequest(http.MethodPost, registerUser, bytes.NewReader(b))
		suite.NoError(err)
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		pqErr := &pq.Error{
			Code: pq.ErrorCode("23505"),
		}

		suite.service.On("RegisterUser", mock.Anything, u.Login, u.Password).Once().Return(pqErr)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.RegisterUser409JSONResponse{
			Message: "user already exist",
		}

		var got api.RegisterUser409JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})

	suite.Run("internal server error (Register user)", func() {
		req, err := http.NewRequest(http.MethodPost, registerUser, bytes.NewReader(b))
		suite.NoError(err)
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.service.On("RegisterUser", mock.Anything, u.Login, u.Password).Once().Return(errors.New("RegisterUser err"))
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.RegisterUser500JSONResponse{
			Message: api.MessageInternalError,
		}

		var got api.RegisterUser500JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})

	suite.Run("internal server error (Authenticate user)", func() {
		req, err := http.NewRequest(http.MethodPost, registerUser, bytes.NewReader(b))
		suite.NoError(err)
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.service.On("RegisterUser", mock.Anything, u.Login, u.Password).Once().Return(nil)
		suite.service.On("Authenticate", mock.Anything, u.Login, u.Password).Once().Return("", errors.New("Authenticate err"))
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.RegisterUser500JSONResponse{
			Message: api.MessageInternalError,
		}

		var got api.RegisterUser500JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})
}

func (suite *apiTestSuite) TestLoginUser() {
	u := api.LoginUserJSONRequestBody{
		Login:    "login",
		Password: "password",
	}
	b, err := json.Marshal(u)
	suite.NoError(err)

	suite.Run("Login user", func() {
		req, err := http.NewRequest(http.MethodPost, loginUser, bytes.NewReader(b))
		suite.NoError(err)
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.service.On("Authenticate", mock.Anything, u.Login, u.Password).Once().Return("token", nil)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()

		suite.Equal("Bearer token", res.Header.Get("Authorization"))
	})

	suite.Run("User unauthorized", func() {
		req, err := http.NewRequest(http.MethodPost, loginUser, bytes.NewReader(b))
		suite.NoError(err)
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.service.On("Authenticate", mock.Anything, u.Login, u.Password).Once().Return("", sql.ErrNoRows)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.LoginUser401JSONResponse{
			Message: "unauthorized",
		}

		var got api.LoginUser401JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})

	suite.Run("internal server error", func() {
		req, err := http.NewRequest(http.MethodPost, loginUser, bytes.NewReader(b))
		suite.NoError(err)
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.service.On("Authenticate", mock.Anything, u.Login, u.Password).Once().Return("", errors.New("Authenticate err"))
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.LoginUser500JSONResponse{
			Message: api.MessageInternalError,
		}

		var got api.LoginUser500JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})
}
func (suite *apiTestSuite) TestGetBalance() {
	userID, err := uuid.NewRandom()
	suite.NoError(err)
	ctx := api.NewContext(context.Background(), userID)

	suite.Run("Get balance", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, getBalance, nil)
		suite.NoError(err)

		rr := httptest.NewRecorder()

		suite.service.On("GetBalance", mock.Anything, userID).Once().Return(12.32, 5.04, nil)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.GetBalance200JSONResponse{
			Current:   12.32,
			Withdrawn: 5.04,
		}

		var got api.GetBalance200JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})

	suite.Run("No balance for current user", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, getBalance, nil)
		suite.NoError(err)

		rr := httptest.NewRecorder()

		suite.service.On("GetBalance", mock.Anything, userID).Once().Return(0.0, 0.0, sql.ErrNoRows)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.GetBalance200JSONResponse{
			Current:   0.0,
			Withdrawn: 0.0,
		}

		var got api.GetBalance200JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})

	suite.Run("Internal server error", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, getBalance, nil)
		suite.NoError(err)

		rr := httptest.NewRecorder()

		suite.service.On("GetBalance", mock.Anything, userID).Once().Return(0.0, 0.0, errors.New("GetBalance err"))
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.GetBalance500JSONResponse{
			Message: api.MessageInternalError,
		}

		var got api.GetBalance500JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})
}

func (suite *apiTestSuite) TestGetWithdrawals() {
	userID, err := uuid.NewRandom()
	suite.NoError(err)
	ctx := api.NewContext(context.Background(), userID)

	suite.Run("Get withdrawals", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, getWithdrawals, nil)
		suite.NoError(err)

		rr := httptest.NewRecorder()

		withdrawals := []model.Withdrawal{
			{Order: "125", Sum: 255.54, ProcessedAt: time.Date(2024, time.January, 1, 1, 1, 0, 0, time.UTC)},
			{Order: "2006", Sum: 1024.0, ProcessedAt: time.Date(2024, time.January, 2, 2, 2, 0, 0, time.UTC)},
		}

		suite.service.On("GetWithdrawals", mock.Anything, userID).Once().Return(withdrawals, nil)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := make(api.GetWithdrawals200JSONResponse, len(withdrawals))
		for i, withdrawal := range withdrawals {
			processedAt := withdrawal.ProcessedAt
			expected[i] = api.Withdrawal{
				Order:       withdrawal.Order,
				Sum:         withdrawal.Sum,
				ProcessedAt: &processedAt,
			}
		}

		var got api.GetWithdrawals200JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})

	suite.Run("No balance for current user", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, getWithdrawals, nil)
		suite.NoError(err)

		rr := httptest.NewRecorder()

		suite.service.On("GetWithdrawals", mock.Anything, userID).Once().Return(nil, sql.ErrNoRows)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()

		suite.Equal(http.StatusNoContent, res.StatusCode)
	})

	suite.Run("Internal server error", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, getWithdrawals, nil)
		suite.NoError(err)

		rr := httptest.NewRecorder()

		suite.service.On("GetWithdrawals", mock.Anything, userID).Once().Return(nil, errors.New("GetWithdrawals err"))
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.GetWithdrawals500JSONResponse{
			Message: api.MessageInternalError,
		}

		var got api.GetWithdrawals500JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})
}

func (suite *apiTestSuite) TestGetOrders() {
	userID, err := uuid.NewRandom()
	suite.NoError(err)
	ctx := api.NewContext(context.Background(), userID)

	suite.Run("Get orders", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, getOrders, nil)
		suite.NoError(err)

		rr := httptest.NewRecorder()

		orders := []model.Order{
			{Number: "125", Accrual: 255.54, Status: "PROCESSED", UploadedAt: time.Date(2024, time.January, 1, 2, 2, 2, 0, time.UTC)},
			{Number: "2006", Accrual: 1024.0, Status: "PROCESSED", UploadedAt: time.Date(2024, time.January, 2, 1, 1, 1, 0, time.UTC)},
		}

		suite.service.On("GetOrders", mock.Anything, userID).Once().Return(orders, nil)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := make(api.GetOrders200JSONResponse, len(orders))
		for i, order := range orders {
			accrual := order.Accrual
			expected[i] = api.Order{
				Number:     order.Number,
				Status:     order.Status,
				Accrual:    &accrual,
				UploadedAt: order.UploadedAt,
			}
		}

		var got api.GetOrders200JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})

	suite.Run("No orders for current user", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, getOrders, nil)
		suite.NoError(err)

		rr := httptest.NewRecorder()

		suite.service.On("GetOrders", mock.Anything, userID).Once().Return(nil, sql.ErrNoRows)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()

		suite.Equal(http.StatusNoContent, res.StatusCode)
	})

	suite.Run("Internal server error", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, getOrders, nil)
		suite.NoError(err)

		rr := httptest.NewRecorder()

		suite.service.On("GetOrders", mock.Anything, userID).Once().Return(nil, errors.New("GetOrders err"))
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.GetOrders500JSONResponse{
			Message: api.MessageInternalError,
		}

		var got api.GetOrders500JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})
}

func (suite *apiTestSuite) TestUploadOrder() {
	userID, err := uuid.NewRandom()
	suite.NoError(err)
	ctx := api.NewContext(context.Background(), userID)

	suite.Run("Upload order", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadOrder, strings.NewReader("125"))
		suite.NoError(err)
		req.Header.Add("Content-Type", "text/plain")
		rr := httptest.NewRecorder()

		suite.service.On("UploadOrder", mock.Anything, userID, mock.Anything).Once().Return(nil)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()

		suite.Equal(http.StatusAccepted, res.StatusCode)
	})

	suite.Run("Invalid order number", func() {
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadOrder, strings.NewReader("126"))
		suite.NoError(err)
		request.Header.Add("Content-Type", "text/plain")

		rr := httptest.NewRecorder()

		suite.r.ServeHTTP(rr, request)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.UploadOrder422JSONResponse{
			Message: "invalid order number",
		}

		var got api.UploadOrder422JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})

	suite.Run("Order already exists", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadOrder, strings.NewReader("125"))
		suite.NoError(err)
		req.Header.Add("Content-Type", "text/plain")

		rr := httptest.NewRecorder()

		pqErr := &pq.Error{
			Code: "23505",
		}

		suite.service.On("UploadOrder", mock.Anything, userID, mock.Anything).Once().Return(pqErr)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.UploadOrder409JSONResponse{
			Message: "order already exists",
		}

		var got api.UploadOrder409JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})

	suite.Run("Order for current user already exists", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadOrder, strings.NewReader("125"))
		suite.NoError(err)
		req.Header.Add("Content-Type", "text/plain")

		rr := httptest.NewRecorder()

		suite.service.On("UploadOrder", mock.Anything, userID, mock.Anything).Once().Return(service.ErrOrderAlreadyExistThisUser)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()

		suite.Equal(http.StatusOK, res.StatusCode)
	})

	suite.Run("Internal server error", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadOrder, strings.NewReader("125"))
		suite.NoError(err)
		req.Header.Add("Content-Type", "text/plain")

		rr := httptest.NewRecorder()

		suite.service.On("UploadOrder", mock.Anything, userID, mock.Anything).Once().Return(errors.New("UploadOrder err"))
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.UploadOrder500JSONResponse{
			Message: api.MessageInternalError,
		}

		var got api.UploadOrder500JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})
}

func (suite *apiTestSuite) TestWithdrawalRequest() {
	userID, err := uuid.NewRandom()
	suite.NoError(err)
	ctx := api.NewContext(context.Background(), userID)

	processedAt := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	withdrawal := api.WithdrawalRequestJSONRequestBody{
		Order:       "125",
		Sum:         125.78,
		ProcessedAt: &processedAt,
	}
	b, err := json.Marshal(withdrawal)
	suite.NoError(err)

	suite.Run("Withdrawal request", func() {
		//ctx := context.WithValue(context.Background(), api.KeyUserID, userID)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, withdrawalReq, bytes.NewReader(b))
		suite.NoError(err)
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.service.On("WithdrawalRequest", mock.Anything, userID, withdrawal.Order, withdrawal.Sum).Once().Return(nil)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()

		suite.Equal(http.StatusOK, res.StatusCode)
	})

	suite.Run("Invalid order number", func() {
		w := api.WithdrawalRequestJSONRequestBody{
			Order:       "126",
			Sum:         125.78,
			ProcessedAt: &processedAt,
		}
		b, err := json.Marshal(w)
		suite.NoError(err)

		//ctx := context.WithValue(context.Background(), api.KeyUserID, userID)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, withdrawalReq, bytes.NewReader(b))
		suite.NoError(err)
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.WithdrawalRequest422JSONResponse{
			Message: "invalid order number",
		}

		var got api.WithdrawalRequest422JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})

	suite.Run("Insufficient balance", func() {
		//ctx := context.WithValue(context.Background(), api.KeyUserID, userID)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, withdrawalReq, bytes.NewReader(b))
		suite.NoError(err)
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.service.On("WithdrawalRequest", mock.Anything, userID, withdrawal.Order, withdrawal.Sum).Once().Return(postgres.ErrNotEnoughFunds)
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.WithdrawalRequest402JSONResponse{
			Message: postgres.ErrNotEnoughFunds.Error(),
		}

		var got api.WithdrawalRequest402JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})

	suite.Run("Internal server error", func() {
		//ctx := context.WithValue(context.Background(), api.KeyUserID, userID)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, withdrawalReq, bytes.NewReader(b))
		suite.NoError(err)
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.service.On("WithdrawalRequest", mock.Anything, userID, withdrawal.Order, withdrawal.Sum).Once().Return(errors.New("WithdrawalRequest err"))
		suite.r.ServeHTTP(rr, req)
		res := rr.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		suite.NoError(err)

		expected := api.WithdrawalRequest500JSONResponse{
			Message: api.MessageInternalError,
		}

		var got api.WithdrawalRequest500JSONResponse
		err = json.Unmarshal(resBody, &got)
		suite.NoError(err)

		suite.Equal(expected, got)
	})
}
